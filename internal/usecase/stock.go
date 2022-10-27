package usecase

import (
	"context"
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
	"gitlab.com/stoqu/stoqu-be/pkg/util/str"
	"gitlab.com/stoqu/stoqu-be/pkg/util/trxmanager"
)

type Stock interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.StockViewResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.StockResponse, error)
	Transaction(ctx context.Context, payload dto.TransactionStockRequest) (dto.StockTransactionResponse, error)
}

type stock struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewStock(cfg *config.Configuration, f repository.Factory) Stock {
	return &stock{f, cfg}
}

func (u *stock) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.StockViewResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := `
			lower(products.code) LIKE ? OR 
			lower(products.name) LIKE ? OR 
			products.price_usd = ? OR 
			products.price_final = ? OR 
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
				searchValFloat,
				searchValFloat,
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
		stockTrxID := uuid.New().String()
		stockTrx := model.StockTrxModel{
			Entity: model.Entity{
				ID: stockTrxID,
			},
			StockTrxEntity: model.StockTrxEntity{
				TrxType: payload.TrxType,
				Code:    str.GenCode(constant.CODE_STOCK_TRX_PREFIX),
			},
		}
		stockTrxItems := []model.StockTrxItemModel{}
		stockTrxItemLookups := []model.StockTrxItemLookupModel{}

		for _, v := range payload.Products {
			var qtySeal, qtyNotSeal int
			stockTrxItemID := uuid.New().String()

			stock, err := u.Repo.Stock.FindByProductID(ctx, v.ID)
			if err != nil {
				return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "stock not found")
			}
			packet, err := u.Repo.Packet.FindByID(ctx, stock.PacketID)
			if err != nil {
				return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "packet not found")
			}

			if payload.TrxType == constant.TRX_TYPE_IN {
				qtySeal = v.Quantity
				stock.TotalSeal += int64(qtySeal)
				stock.Total += int64(qtySeal)

				stockLookups := []model.StockLookupModel{}
				for i := 0; i < v.Quantity; i++ {
					stockLookup := model.StockLookupModel{
						StockLookupEntity: model.StockLookupEntity{
							Code:           str.GenCode(constant.CODE_STOCK_LOOKUP_PREFIX),
							IsSeal:         true,
							Value:          float64(packet.Value),
							RemainingValue: float64(packet.Value),
							StockID:        stock.ID,
						},
					}
					stockLookups = append(stockLookups, stockLookup)
					stockTrxItemLookups = append(stockTrxItemLookups, model.StockTrxItemLookupModel{
						StockTrxItemLookupEntity: model.StockTrxItemLookupEntity{
							StockLookupEntity: stockLookup.StockLookupEntity,
							StockTrxItemID:    stockTrxItemID,
						},
					})
				}

				_, err = u.Repo.StockLookup.Creates(ctx, stockLookups)
				if err != nil {
					return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "failed create stock lookups")
				}
			} else {
				stockLookups, err := u.Repo.StockLookup.FindByIDs(ctx, v.StockLookupIDs)
				if err != nil {
					return res.ErrorBuilder(res.Constant.Error.BadRequest, err, "failed find stock lookups")
				}

				for _, v2 := range stockLookups {
					if v2.IsSeal {
						qtySeal++
					} else {
						qtyNotSeal++
					}

					stockTrxItemLookups = append(stockTrxItemLookups, model.StockTrxItemLookupModel{
						StockTrxItemLookupEntity: model.StockTrxItemLookupEntity{
							StockLookupEntity: v2.StockLookupEntity,
							StockTrxItemID:    stockTrxItemID,
						},
					})
				}
				stock.TotalSeal -= int64(qtySeal)
				stock.TotalNotSeal -= int64(qtyNotSeal)
				stock.Total -= int64(qtySeal) + int64(qtyNotSeal)

				err = u.Repo.StockLookup.DeleteByIDs(ctx, v.StockLookupIDs)
				if err != nil {
					return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "failed delete stock lookups")
				}
			}

			stockTrxItems = append(stockTrxItems, model.StockTrxItemModel{
				Entity: model.Entity{
					ID: stockTrxItemID,
				},
				StockTrxItemEntity: model.StockTrxItemEntity{
					TotalSeal:    qtySeal,
					TotalNotSeal: qtyNotSeal,
					StockTrxID:   stockTrx.Entity.ID,
					StockID:      stock.ID,
					ProductID:    stock.ProductID,
				},
			})

			_, err = u.Repo.Stock.UpdateByID(ctx, stock.ID, *stock)
			if err != nil {
				return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "failed update stock")
			}
		}

		_, err = u.Repo.StockTrx.Create(ctx, stockTrx)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "failed create stock trx")
		}
		_, err = u.Repo.StockTrxItem.Creates(ctx, stockTrxItems)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "failed create stock trx items")
		}
		_, err = u.Repo.StockTrxItemLookup.Creates(ctx, stockTrxItemLookups)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err, "failed create stock trx item lookups")
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.StockTransactionResponse{
		Status: constant.STATUS_SUCCESS,
	}

	return result, nil
}
