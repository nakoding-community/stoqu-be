package usecase

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gitlab.com/stoqu/stoqu-be/internal/config"

	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
	errConstant "gitlab.com/stoqu/stoqu-be/pkg/util/response/constant"
	"gitlab.com/stoqu/stoqu-be/pkg/util/str"
	"gitlab.com/stoqu/stoqu-be/pkg/util/trxmanager"
)

type Stock interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.StockViewResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.StockResponse, error)
	Transaction(ctx context.Context, payload dto.TransactionStockRequest) (dto.StockTransactionResponse, error)
	Convertion(ctx context.Context, payload dto.ConvertionStockRequest) (dto.StockConvertionResponse, error)
	Movement(ctx context.Context, payload dto.MovementStockRequest) (result dto.StockMovementResponse, err error)
	History(ctx context.Context, filterParam abstraction.Filter) ([]dto.StockHistoryResponse, abstraction.PaginationInfo, error)

	// helper
	TransactionProcess(ctx context.Context, payload dto.TransactionStockRequest) (result dto.StockTransactionResponse, err error)
}

type stock struct {
	Cfg  *config.Configuration
	Repo repository.Factory
}

func NewStock(cfg *config.Configuration, f repository.Factory) Stock {
	return &stock{cfg, f}
}

func (u *stock) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.StockViewResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := `
			lower(products.code) LIKE ? OR 
			lower(products.name) LIKE ? OR 
			lower(brands.name) LIKE ? OR 
			lower(users.name) LIKE ? OR 
			lower(variants.name) LIKE ? OR 
			lower(units.name) LIKE ? OR
			packets.value = ?
		`
		searchVal := "%" + strings.ToLower(filterParam.Search) + "%"
		searchValFloat, _ := strconv.ParseFloat(filterParam.Search, 64)
		search = &abstraction.Search{
			Query: searchQuery,
			Args: []interface{}{
				searchVal,
				searchVal,
				searchVal,
				searchVal,
				searchVal,
				searchVal,
				searchValFloat,
			},
		}
	}

	stocks, info, err := u.Repo.Stock.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, stock := range stocks {
		stockRacks, err := u.Repo.StockRack.FindByStockID(ctx, stock.ID)
		if err == nil {
			for j, stockRack := range stockRacks {
				rack, err := u.Repo.Rack.FindByID(ctx, stockRack.RackID)
				if err == nil {
					stockRacks[j].Rack = rack
				}
			}
			stock.StockRacks = stockRacks
		}

		result = append(result, dto.StockViewResponse{
			StockView: stock,
		})
	}

	return result, pagination, nil
}

func (u *stock) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.StockResponse, error) {
	var result dto.StockResponse

	stock, err := u.Repo.Stock.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.StockResponse{
		StockModel: *stock,
	}

	return result, nil
}

func (u *stock) Transaction(ctx context.Context, payload dto.TransactionStockRequest) (result dto.StockTransactionResponse, err error) {
	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		result, err = u.TransactionProcess(ctx, payload)
		return err
	}); err != nil {
		return result, err
	}

	result.Status = constant.STATUS_SUCCESS
	return result, nil
}

func (u *stock) Convertion(ctx context.Context, payload dto.ConvertionStockRequest) (result dto.StockConvertionResponse, err error) {
	wrapper := &convertionDataWrapper{}
	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		// validate data
		err := u.convertionValidateData(ctx, payload, wrapper)
		if err != nil {
			return err
		}

		// validate calculation
		err = u.convertionValidateCalculation(ctx, payload, wrapper)
		if err != nil {
			return err
		}

		// process
		err = u.convertionProcess(ctx, payload, wrapper)
		if err != nil {
			return err
		}

		// mutation
		err = u.convertionMutation(ctx, wrapper)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result.Products = append(result.Products, dto.StockConvertionProductResponse{
		ID:          wrapper.destination.product.ID,
		LookupCodes: []string{wrapper.destination.product.Code},
	})

	result.Status = constant.STATUS_SUCCESS
	return result, nil
}

func (u *stock) Movement(ctx context.Context, payload dto.MovementStockRequest) (result dto.StockMovementResponse, err error) {
	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		// validate
		if payload.Origin.RackID == payload.Destination.RackID {
			return res.ErrorBuilder(res.Constant.Error.Validation, errors.New("can't move stock in same rack id"))
		}
		stock, err := u.Repo.Stock.FindByProductID(ctx, payload.Origin.ProductID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find stock by product id error")
		}
		stockRackOrigin, err := u.Repo.StockRack.FindByStockAndRackID(ctx, stock.ID, payload.Origin.RackID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find stock rack origin by stock and rack id error")
		}
		stockRackDestination, err := u.upsertStockRack(ctx, stock.ID, payload.Destination.RackID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "upsert stock rack destination by stock and rack id error")
		}
		stockLookupOrigins, err := u.Repo.StockLookup.FindByIDs(ctx, payload.Origin.StockLookupIDs, "")
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find stock lookups origin by ids error")
		}

		// process
		var totalSeal, totalNotSeal int
		for i, stockLookup := range stockLookupOrigins {
			if stockLookup.IsSeal {
				totalSeal++
			} else {
				totalNotSeal++
			}

			stockLookupOrigins[i].StockRackID = stockRackDestination.ID
		}
		stockRackOrigin.Total -= int64(len(stockLookupOrigins))
		stockRackOrigin.TotalSeal -= int64(totalSeal)
		stockRackOrigin.TotalNotSeal -= int64(totalNotSeal)
		stockRackDestination.Total += int64(len(stockLookupOrigins))
		stockRackDestination.TotalSeal += int64(totalSeal)
		stockRackDestination.TotalNotSeal += int64(totalNotSeal)

		// mutation
		for _, stockLookup := range stockLookupOrigins {
			_, err = u.Repo.StockLookup.UpdateByID(ctx, stockLookup.ID, stockLookup)
			if err != nil {
				return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "update stock lookup origin by id error")
			}
		}
		_, err = u.Repo.StockRack.UpdateByID(ctx, stockRackOrigin.ID, *stockRackOrigin)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "update stock rack origin by id error")
		}
		_, err = u.Repo.StockRack.UpdateByID(ctx, stockRackDestination.ID, *stockRackDestination)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "update stock destination origin by id error")
		}

		return nil
	}); err != nil {
		return result, err
	}

	result.Status = constant.STATUS_SUCCESS
	return result, nil
}

func (u *stock) History(ctx context.Context, filterParam abstraction.Filter) (result []dto.StockHistoryResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := `
			lower(stock_trxs.code) LIKE ? OR 
			lower(stock_trxs.trx_type) LIKE ? OR 
		`
		searchVal := "%" + strings.ToLower(filterParam.Search) + "%"
		search = &abstraction.Search{
			Query: searchQuery,
			Args: []interface{}{
				searchVal,
				searchVal,
			},
		}
	}

	stockTrxs, info, err := u.Repo.StockTrx.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, err
	}
	pagination = *info

	for i, stockTrx := range stockTrxs {
		stockTrxItems, _, err := u.Repo.StockTrxItem.Find(ctx, abstraction.Filter{
			Query: []abstraction.FilterQuery{
				{
					Field: "stock_trx_id",
					Value: stockTrx.ID,
				},
			},
		}, search)
		if err != nil {
			return nil, pagination, err
		}

		for j, stockTrxItem := range stockTrxItems {
			page, limit := 1, 1
			products, _, err := u.Repo.Product.Find(ctx, abstraction.Filter{
				Query: []abstraction.FilterQuery{
					{
						Field: "products.id",
						Value: stockTrxItem.ProductID,
					},
				},
				Pagination: abstraction.Pagination{
					Page:  &page,
					Limit: &limit,
				},
			}, search)
			if err == nil && len(products) > 0 {
				stockTrxItems[j].Product = &products[0]
			}

			stockTrxItemLookups, _, err := u.Repo.StockTrxItemLookup.Find(ctx, abstraction.Filter{
				Query: []abstraction.FilterQuery{
					{
						Field: "stock_trx_item_id",
						Value: stockTrxItem.ID,
					},
				},
			}, search)
			if err == nil {
				stockTrxItems[j].StockTrxItemLookups = stockTrxItemLookups
			}
		}
		stockTrxs[i].StockTrxItems = stockTrxItems

		result = append(result, dto.StockHistoryResponse{
			StockTrxModel: stockTrxs[i],
		})
	}

	return result, pagination, nil
}

// helper
func (u *stock) upsertStockRack(ctx context.Context, stockID, rackID string) (result *model.StockRackModel, err error) {
	_, err = u.Repo.Rack.FindByID(ctx, rackID)
	if err != nil {
		if res.ErrorResponse(err).ErrorCode() == errConstant.E_DATA_NOTFOUND {
			return result, res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find rack by id not found")
		}
		return result, err
	}

	result, err = u.Repo.StockRack.FindByStockAndRackID(ctx, stockID, rackID)
	if err != nil {
		if res.ErrorResponse(err).ErrorCode() == errConstant.E_DATA_NOTFOUND {
			stockRackInsert, err := u.Repo.StockRack.Create(ctx, model.StockRackModel{
				Entity: model.Entity{
					ID: uuid.New().String(),
				},
				StockRackEntity: model.StockRackEntity{
					StockID: stockID,
					RackID:  rackID,
				},
			})
			if err != nil {
				return result, res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "create stock rack error")
			}
			result = &stockRackInsert
		} else {
			return result, res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find stock rack error")
		}
	}

	return result, nil
}

type transactionDataWrapper struct {
	req dto.TransactionStockProductRequest

	stock     *model.StockModel
	stockRack *model.StockRackModel
	packet    *model.PacketModel

	stockTrxItemID      string
	stockTrxItemLookups []model.StockTrxItemLookupModel
	qtySeal, qtyNotSeal int
}

func (u *stock) TransactionProcess(ctx context.Context, payload dto.TransactionStockRequest) (result dto.StockTransactionResponse, err error) {
	stockTrxID := uuid.New().String()
	stockTrx := model.StockTrxModel{
		Entity: model.Entity{
			ID: stockTrxID,
		},
		StockTrxEntity: model.StockTrxEntity{
			Code:       str.GenCode(constant.CODE_STOCK_TRX_PREFIX),
			TrxType:    payload.TrxType,
			OrderTrxID: payload.OrderTrxID,
		},
	}
	stockTrxItems := []model.StockTrxItemModel{}
	stockTrxItemLookups := []model.StockTrxItemLookupModel{}

	for _, v := range payload.Products {
		stockTrxItemID := uuid.New().String()

		product, err := u.Repo.Product.FindByID(ctx, v.ID)
		if err != nil {
			return result, res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find product by id error")
		}
		stock, err := u.Repo.Stock.FindByProductID(ctx, v.ID)
		if err != nil {
			return result, res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find stock by product id error")
		}
		stockRack, err := u.upsertStockRack(ctx, stock.ID, v.RackID)
		if err != nil {
			return result, err
		}
		packet, err := u.Repo.Packet.FindByID(ctx, stock.PacketID)
		if err != nil {
			return result, res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find packet by id error")
		}

		wrapper := &transactionDataWrapper{
			req:            v,
			stock:          stock,
			stockRack:      stockRack,
			packet:         packet,
			stockTrxItemID: stockTrxItemID,
		}
		if payload.TrxType == constant.TRX_TYPE_IN {
			err = u.transactionTypeIn(ctx, wrapper)
			if err != nil {
				return result, err
			}
		} else {
			err = u.transactionTypeOut(ctx, wrapper)
			if err != nil {
				return result, err
			}
		}

		_, err = u.Repo.StockRack.UpdateByID(ctx, stockRack.ID, *wrapper.stockRack)
		if err != nil {
			return result, res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "update stock rack error")
		}
		_, err = u.Repo.Stock.UpdateByID(ctx, stock.ID, *wrapper.stock)
		if err != nil {
			return result, res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "update stock error")
		}

		stockTrxItems = append(stockTrxItems, model.StockTrxItemModel{
			Entity: model.Entity{
				ID: stockTrxItemID,
			},
			StockTrxItemEntity: model.StockTrxItemEntity{
				TotalSeal:    wrapper.qtySeal,
				TotalNotSeal: wrapper.qtyNotSeal,
				StockTrxID:   stockTrx.Entity.ID,
				StockID:      stock.ID,
				ProductID:    stock.ProductID,
			},
		})
		stockTrxItemLookups = append(stockTrxItemLookups, wrapper.stockTrxItemLookups...)

		resultProduct := dto.StockTransactionProductResponse{
			ID:          v.ID,
			LookupCodes: []string{product.Code},
		}
		result.Products = append(result.Products, resultProduct)
	}

	_, err = u.Repo.StockTrx.Create(ctx, stockTrx)
	if err != nil {
		return result, res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "create stock trx error")
	}
	_, err = u.Repo.StockTrxItem.Creates(ctx, stockTrxItems)
	if err != nil {
		return result, res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "create stock trx items error")
	}
	_, err = u.Repo.StockTrxItemLookup.Creates(ctx, stockTrxItemLookups)
	if err != nil {
		return result, res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "create stock trx item lookups error")
	}

	return result, nil
}

func (u *stock) transactionTypeIn(ctx context.Context, wrapper *transactionDataWrapper) (err error) {
	stockLookups := []model.StockLookupModel{}

	if len(wrapper.req.StockTrxItemLookupIDs) == 0 {
		wrapper.qtySeal = wrapper.req.Quantity
		for i := 0; i < wrapper.req.Quantity; i++ {
			stockLookup := model.StockLookupModel{
				StockLookupEntity: model.StockLookupEntity{
					Code:           str.GenCode(constant.CODE_STOCK_LOOKUP_PREFIX),
					IsSeal:         true,
					Value:          float64(wrapper.packet.Value),
					RemainingValue: float64(wrapper.packet.Value),
					StockRackID:    wrapper.stockRack.ID,
				},
			}
			stockLookups = append(stockLookups, stockLookup)
			wrapper.stockTrxItemLookups = append(wrapper.stockTrxItemLookups, model.StockTrxItemLookupModel{
				StockTrxItemLookupEntity: model.StockTrxItemLookupEntity{
					StockLookupEntity: stockLookup.StockLookupEntity,
					Code:              stockLookup.Code,
					StockTrxItemID:    wrapper.stockTrxItemID,
				},
			})
		}
		wrapper.stockRack.TotalSeal += int64(wrapper.qtySeal)
		wrapper.stockRack.Total += int64(wrapper.qtySeal)
		wrapper.stock.TotalSeal += int64(wrapper.qtySeal)
		wrapper.stock.Total += int64(wrapper.qtySeal)

	} else {
		stockTrxItemLookups, err := u.Repo.StockTrxItemLookup.FindByIDs(ctx, wrapper.req.StockTrxItemLookupIDs, "")
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find stock trx item lookups by ids error")
		}

		for _, v2 := range stockTrxItemLookups {
			if v2.IsSeal {
				wrapper.qtySeal++
			} else {
				wrapper.qtyNotSeal++
			}

			stockLookup := model.StockLookupModel{
				StockLookupEntity: v2.StockLookupEntity,
			}
			stockLookup.Code = v2.Code
			stockLookups = append(stockLookups, stockLookup)
			wrapper.stockTrxItemLookups = append(wrapper.stockTrxItemLookups, model.StockTrxItemLookupModel{
				StockTrxItemLookupEntity: model.StockTrxItemLookupEntity{
					StockLookupEntity: stockLookup.StockLookupEntity,
					Code:              stockLookup.Code,
					StockTrxItemID:    wrapper.stockTrxItemID,
				},
			})
		}
		wrapper.stockRack.TotalSeal += int64(wrapper.qtySeal)
		wrapper.stockRack.TotalNotSeal += int64(wrapper.qtyNotSeal)
		wrapper.stockRack.Total += int64(wrapper.qtySeal) + int64(wrapper.qtyNotSeal)
		wrapper.stock.TotalSeal += int64(wrapper.qtySeal)
		wrapper.stock.TotalNotSeal += int64(wrapper.qtyNotSeal)
		wrapper.stock.Total += int64(wrapper.qtySeal) + int64(wrapper.qtyNotSeal)
	}

	_, err = u.Repo.StockLookup.Creates(ctx, stockLookups)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "create stock lookups error")
	}

	return nil
}

func (u *stock) transactionTypeOut(ctx context.Context, wrapper *transactionDataWrapper) (err error) {
	stockLookups, err := u.Repo.StockLookup.FindByIDs(ctx, wrapper.req.StockLookupIDs, "")
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find stock lookups by ids error")
	}

	for _, v2 := range stockLookups {
		if v2.IsSeal {
			wrapper.qtySeal++
		} else {
			wrapper.qtyNotSeal++
		}

		wrapper.stockTrxItemLookups = append(wrapper.stockTrxItemLookups, model.StockTrxItemLookupModel{
			StockTrxItemLookupEntity: model.StockTrxItemLookupEntity{
				StockLookupEntity: v2.StockLookupEntity,
				Code:              v2.StockLookupEntity.Code,
				StockTrxItemID:    wrapper.stockTrxItemID,
			},
		})
	}
	wrapper.stockRack.TotalSeal -= int64(wrapper.qtySeal)
	wrapper.stockRack.TotalNotSeal -= int64(wrapper.qtyNotSeal)
	wrapper.stockRack.Total -= int64(wrapper.qtySeal) + int64(wrapper.qtyNotSeal)
	wrapper.stock.TotalSeal -= int64(wrapper.qtySeal)
	wrapper.stock.TotalNotSeal -= int64(wrapper.qtyNotSeal)
	wrapper.stock.Total -= int64(wrapper.qtySeal) + int64(wrapper.qtyNotSeal)

	err = u.Repo.StockLookup.DeleteByIDs(ctx, wrapper.req.StockLookupIDs)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "delete stock lookups error")
	}
	return nil
}

type convertionDataWrapper struct {
	origin         *convertionProductWrapper
	destination    *convertionProductWrapper
	convertionUnit *model.ConvertionUnitModel
	trx            *convertionStockTrxWrapper
}

type convertionProductWrapper struct {
	product      *model.ProductModel
	stock        *model.StockModel
	stockRack    *model.StockRackModel
	packet       *model.PacketModel
	stockLookups []model.StockLookupModel
	totalValue   float64
}

type convertionStockTrxWrapper struct {
	stockTrx            *model.StockTrxModel
	stockTrxItems       []model.StockTrxItemModel
	stockTrxItemLookups []model.StockTrxItemLookupModel
}

func (u *stock) convertionValidateData(ctx context.Context, payload dto.ConvertionStockRequest, wrapper *convertionDataWrapper) error {
	productOrigin, err := u.Repo.Product.FindByID(ctx, payload.Origin.ProductID)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find product origin by id error")
	}
	stockOrigin, err := u.Repo.Stock.FindByProductID(ctx, productOrigin.ID)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find stock origin by product id error")
	}
	stockRackOrigin, err := u.Repo.StockRack.FindByStockAndRackID(ctx, stockOrigin.ID, payload.Origin.RackID)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find stock rack origin by product id error")
	}
	packetOrigin, err := u.Repo.Packet.FindByID(ctx, productOrigin.PacketID)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find packet origin by id error")
	}

	productDestination, err := u.Repo.Product.FindByBrandVariantPacketID(ctx, productOrigin.BrandID, productOrigin.VariantID, payload.Destination.PacketID)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find product destination by id error")
	}
	stockDestination, err := u.Repo.Stock.FindByProductID(ctx, productDestination.ID)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find stock destination by product id error")
	}
	stockRackDestination, err := u.upsertStockRack(ctx, stockDestination.ID, payload.Origin.RackID)
	if err != nil {
		return err
	}
	packetDestination, err := u.Repo.Packet.FindByID(ctx, payload.Destination.PacketID)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find packet destination by id error")
	}

	convertionUnit, err := u.Repo.ConvertionUnit.FindByUnitOriginDestinationID(ctx, packetOrigin.UnitID, packetDestination.UnitID)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find convertion unit by unit origin & destination id error")
	}

	wrapper.origin = &convertionProductWrapper{
		product:   productOrigin,
		stock:     stockOrigin,
		stockRack: stockRackOrigin,
		packet:    packetOrigin,
	}
	wrapper.destination = &convertionProductWrapper{
		product:   productDestination,
		stock:     stockDestination,
		stockRack: stockRackDestination,
		packet:    packetDestination,
	}
	wrapper.convertionUnit = convertionUnit

	return nil
}

func (u *stock) convertionValidateCalculation(ctx context.Context, payload dto.ConvertionStockRequest, wrapper *convertionDataWrapper) error {
	if wrapper.origin.packet.UnitID == wrapper.destination.packet.UnitID && wrapper.origin.packet.Value <= wrapper.destination.packet.Value {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, errors.New("packet destination must less than origin"))
	}
	sumStockLookup, err := u.Repo.StockLookup.SumByIDs(ctx, payload.Origin.StockLookupIDs)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "sum stock lookup by ids error")
	}
	if sumStockLookup.RemainingValue == 0 {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, errors.New("remaining value must be more than 0"))
	}
	totalValueOrigin := sumStockLookup.RemainingValue * wrapper.convertionUnit.ValueConvertion
	totalValueDestination := float64(payload.Destination.Total) * float64(wrapper.destination.packet.Value)
	if totalValueOrigin < totalValueDestination {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, errors.New("total value origin must be more than destination"))
	}

	wrapper.origin.totalValue = totalValueOrigin
	wrapper.destination.totalValue = totalValueDestination

	return nil
}

func (u *stock) convertionProcess(ctx context.Context, payload dto.ConvertionStockRequest, wrapper *convertionDataWrapper) (err error) {
	// prepare
	stockTrxID := uuid.New().String()
	stockTrx := model.StockTrxModel{
		Entity: model.Entity{
			ID: stockTrxID,
		},
		StockTrxEntity: model.StockTrxEntity{
			TrxType: constant.TRX_TYPE_CONVERT,
			Code:    str.GenCode(constant.CODE_STOCK_TRX_PREFIX),
		},
	}
	stockTrxItems := []model.StockTrxItemModel{
		{
			Entity: model.Entity{
				ID: uuid.New().String(),
			},
			StockTrxItemEntity: model.StockTrxItemEntity{
				ConvertType: constant.CONVERT_TYPE_ORIGIN,
				StockTrxID:  stockTrxID,
				StockID:     wrapper.origin.stock.ID,
				ProductID:   wrapper.origin.product.ID,
			},
		},
		{
			Entity: model.Entity{
				ID: uuid.New().String(),
			},
			StockTrxItemEntity: model.StockTrxItemEntity{
				ConvertType: constant.CONVERT_TYPE_DESTINATION,
				StockTrxID:  stockTrxID,
				StockID:     wrapper.destination.stock.ID,
				ProductID:   wrapper.destination.product.ID,
			},
		},
	}
	stockTrxItemLookups := []model.StockTrxItemLookupModel{}
	stockLookupOrigins, err := u.Repo.StockLookup.FindByIDs(ctx, payload.Origin.StockLookupIDs, "remaining_value asc")
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "find stock lookup by ids error")
	}
	stockLookupDestinations := []model.StockLookupModel{}
	decrementValue := float64(wrapper.destination.packet.Value) / wrapper.convertionUnit.ValueConvertion

	// process
	for i := 0; i < int(payload.Destination.Total); i++ {
		for j := range stockLookupOrigins {
			// prevention
			if stockLookupOrigins[j].RemainingValue <= 0 {
				continue
			}
			if stockLookupOrigins[j].RemainingValue < decrementValue && wrapper.origin.totalValue < wrapper.destination.totalValue {
				break
			}
			wrapper.origin.totalValue -= float64(wrapper.destination.packet.Value)

			// origin
			stockLookupOrigins[j].RemainingValueBefore = stockLookupOrigins[j].RemainingValue
			stockLookupOrigins[j].RemainingValue -= decrementValue
			if stockLookupOrigins[j].IsSeal {
				stockLookupOrigins[j].IsSeal = false
				wrapper.origin.stockRack.TotalNotSeal++
				wrapper.origin.stockRack.TotalSeal--
				wrapper.origin.stock.TotalNotSeal++
				wrapper.origin.stock.TotalSeal--
				stockTrxItems[0].TotalSeal++
			} else {
				stockTrxItems[0].TotalNotSeal++
			}
			stockTrxItemLookups = append(stockTrxItemLookups, model.StockTrxItemLookupModel{
				StockTrxItemLookupEntity: model.StockTrxItemLookupEntity{
					StockLookupEntity: stockLookupOrigins[j].StockLookupEntity,
					Code:              stockLookupOrigins[j].Code,
					StockTrxItemID:    stockTrxItems[0].ID,
				},
			})

			// destination
			stockLookupDestination := model.StockLookupModel{
				StockLookupEntity: model.StockLookupEntity{
					Code:           str.GenCode(constant.CODE_STOCK_LOOKUP_PREFIX),
					IsSeal:         true,
					Value:          float64(wrapper.destination.packet.Value),
					RemainingValue: float64(wrapper.destination.packet.Value),
					StockRackID:    wrapper.destination.stockRack.ID,
				},
			}
			stockLookupDestinations = append(stockLookupDestinations, stockLookupDestination)
			wrapper.destination.stockRack.TotalSeal++
			wrapper.destination.stockRack.Total++
			wrapper.destination.stock.TotalSeal++
			wrapper.destination.stock.Total++
			stockTrxItems[1].TotalSeal++
			stockTrxItemLookups = append(stockTrxItemLookups, model.StockTrxItemLookupModel{
				StockTrxItemLookupEntity: model.StockTrxItemLookupEntity{
					StockLookupEntity: stockLookupDestination.StockLookupEntity,
					Code:              stockLookupDestination.Code,
					StockTrxItemID:    stockTrxItems[1].ID,
				},
			})
			break
		}
	}

	wrapper.origin.stockLookups = stockLookupOrigins
	wrapper.destination.stockLookups = stockLookupDestinations
	wrapper.trx = &convertionStockTrxWrapper{
		stockTrx:            &stockTrx,
		stockTrxItems:       stockTrxItems,
		stockTrxItemLookups: stockTrxItemLookups,
	}

	return nil
}

func (u *stock) convertionMutation(ctx context.Context, wrapper *convertionDataWrapper) (err error) {
	// stock lookup
	for _, stockLookupOrigin := range wrapper.origin.stockLookups {
		_, err = u.Repo.StockLookup.UpdateByID(ctx, stockLookupOrigin.ID, stockLookupOrigin)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "update stock lookup origin by id error")
		}
	}
	_, err = u.Repo.StockLookup.Creates(ctx, wrapper.destination.stockLookups)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "creates stock lookups destination error")
	}

	// stock rack
	_, err = u.Repo.StockRack.UpdateByID(ctx, wrapper.origin.stockRack.ID, *wrapper.origin.stockRack)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "update stock rack origin by id error")
	}
	_, err = u.Repo.StockRack.UpdateByID(ctx, wrapper.destination.stockRack.ID, *wrapper.destination.stockRack)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "update stock rack destination by id error")
	}

	// stock
	_, err = u.Repo.Stock.UpdateByID(ctx, wrapper.origin.stock.ID, *wrapper.origin.stock)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "update stock origin by id error")
	}
	_, err = u.Repo.Stock.UpdateByID(ctx, wrapper.destination.stock.ID, *wrapper.destination.stock)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "update stock destination by id error")
	}

	// stock trx
	_, err = u.Repo.StockTrx.Create(ctx, *wrapper.trx.stockTrx)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "create stock trx error")
	}
	_, err = u.Repo.StockTrxItem.Creates(ctx, wrapper.trx.stockTrxItems)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "creates stock trx items error")
	}
	_, err = u.Repo.StockTrxItemLookup.Creates(ctx, wrapper.trx.stockTrxItemLookups)
	if err != nil {
		return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "creates stock trx item lookups error")
	}

	return nil
}
