package usecase

import (
	"context"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.com/stoqu/stoqu-be/internal/config"

	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
	"gitlab.com/stoqu/stoqu-be/pkg/util/trxmanager"
)

type Order interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.OrderViewResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.OrderViewResponse, error)
	Upsert(ctx context.Context, payload dto.UpsertOrderRequest) (dto.OrderUpsertResponse, error)
}

type order struct {
	Cfg     *config.Configuration
	Repo    repository.Factory
	stockUC Stock
}

func NewOrder(cfg *config.Configuration, f repository.Factory, stockUC Stock) Order {
	return &order{cfg, f, stockUC}
}

func (u *order) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.OrderViewResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := `
			lower(order_trxs.code) LIKE ? OR 
			lower(order_trxs.trx_type) LIKE ? OR

			order_trxs.price = ? OR
			order_trxs.final_price = ? OR
			
			lower(order_trxs.stock_status) LIKE ? OR
			lower(order_trxs.payment_status) LIKE ? OR
			lower(order_trxs.status) LIKE ? OR

			lower(customers.name) LIKE ? OR
			lower(suppliers.name) LIKE ? OR
			lower(pics.name) LIKE ?
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
				searchVal,
				searchVal,
			},
		}
	}

	orders, info, err := u.Repo.OrderTrx.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, order := range orders {
		result = append(result, dto.OrderViewResponse{
			OrderView: order,
		})
	}

	return result, pagination, nil
}

func (u *order) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.OrderViewResponse, error) {
	var result dto.OrderViewResponse

	order, err := u.Repo.OrderTrx.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.OrderViewResponse{
		OrderView: *order,
	}

	return result, nil
}

func (u *order) Upsert(ctx context.Context, payload dto.UpsertOrderRequest) (result dto.OrderUpsertResponse, err error) {
	mapStockLookups, err := u.buildMapStockLookups(ctx, payload)
	if err != nil {
		return result, err
	}
	orderTrx := payload.ToOrderTrx(mapStockLookups)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(_ctx context.Context) error {
		if payload.ID == "" {
			err := u.create(_ctx, payload, orderTrx)
			if err != nil {
				return err
			}
		} else {
			err := u.update(_ctx, payload, orderTrx)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return result, err
	}

	orderView, err := u.Repo.OrderTrx.FindByID(ctx, orderTrx.ID)
	if err != nil {
		return result, err
	}

	// sync firestore
	orderTrxFs := orderTrx.ToOrderTrxFs(*orderView)
	if payload.ID == "" {
		if err = u.Repo.OrderFs.Add(constant.FIRESTORE_COLLECTION_DASHBOARD_ORDER, orderTrx.ID, orderTrxFs); err != nil {
			logrus.Error("error add order fs, in collection `dashboard-order`", err)
		}
	} else {
		if err = u.Repo.OrderFs.Update(constant.FIRESTORE_COLLECTION_DASHBOARD_ORDER, orderTrx.ID, orderTrxFs); err != nil {
			logrus.Error("error update order fs, in collection `dashboard-order`", err)
		}
	}

	count, err := u.Repo.OrderTrx.Count(ctx)
	if err != nil {
		return result, err
	}
	if err = u.Repo.OrderFs.UpdateTotal(constant.FIRESTORE_COLLECTION_TOTAL_ORDER, "1", entity.OrderTrxTotalFs{
		ID:         "1",
		TotalOrder: count,
	}); err != nil {
		logrus.Error("error update total order fs, in collection `total-order`", err)
	}

	result = dto.OrderUpsertResponse{
		Status: constant.STATUS_SUCCESS,
	}

	return result, nil
}

// helper
func (u *order) buildMapStockLookups(ctx context.Context, payload dto.UpsertOrderRequest) (map[string]entity.StockLookupModel, error) {
	mapStockLookups := make(map[string]entity.StockLookupModel)
	ids := []string{}
	for _, item := range payload.Items {
		for _, lookup := range item.StockLookups {
			ids = append(ids, lookup.ID)
		}
	}

	stockLookups, err := u.Repo.StockLookup.FindByIDs(ctx, ids, "")
	if err != nil {
		return mapStockLookups, err
	}
	for _, stockLookup := range stockLookups {
		mapStockLookups[stockLookup.ID] = stockLookup
	}

	return mapStockLookups, nil
}

func (u *order) create(ctx context.Context, payload dto.UpsertOrderRequest, orderTrx entity.OrderTrxModel) error {
	// create order
	_, err := u.Repo.OrderTrx.Create(ctx, orderTrx)
	if err != nil {
		return err
	}
	_, err = u.Repo.OrderTrxItem.Creates(ctx, orderTrx.OrderTrxItems)
	if err != nil {
		return err
	}

	trxStockProducts := []dto.TransactionStockProductRequest{}
	for _, item := range orderTrx.OrderTrxItems {
		_, err = u.Repo.OrderTrxItemLookup.Creates(ctx, item.OrderTrxItemLookups)
		if err != nil {
			return err
		}

		if len(item.OrderTrxItemLookups) > 0 {
			trxStockLookupIDs := []string{}
			for _, lookup := range item.OrderTrxItemLookups {
				trxStockLookupIDs = append(trxStockLookupIDs, lookup.ID)
			}
			trxStockProduct := dto.TransactionStockProductRequest{
				ID:             item.ProductID,
				Quantity:       len(item.OrderTrxItemLookups),
				RackID:         item.RackID,
				StockLookupIDs: trxStockLookupIDs,
			}
			trxStockProducts = append(trxStockProducts, trxStockProduct)
		}
	}
	_, err = u.Repo.OrderTrxReceipt.Creates(ctx, orderTrx.OrderTrxReceipts)
	if err != nil {
		return err
	}

	// transaction stock
	if len(trxStockProducts) > 0 {
		_, err = u.stockUC.TransactionProcess(ctx, dto.TransactionStockRequest{
			TrxType:    constant.TRX_TYPE_OUT,
			OrderTrxID: orderTrx.ID,
			Products:   trxStockProducts,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (u *order) update(ctx context.Context, payload dto.UpsertOrderRequest, orderTrx entity.OrderTrxModel) error {
	// upsert order
	_, err := u.Repo.OrderTrx.UpdateByID(ctx, orderTrx.ID, orderTrx)
	if err != nil {
		return err
	}

	trxStockProductsIn := []dto.TransactionStockProductRequest{}
	trxStockProductsOut := []dto.TransactionStockProductRequest{}
	for _, item := range orderTrx.OrderTrxItems {
		trxStockProduct := dto.TransactionStockProductRequest{
			ID:     item.ProductID,
			RackID: item.RackID,
		}
		trxStockProductIn := trxStockProduct
		trxStockProductOut := trxStockProduct

		if item.Action == constant.ACTION_DELETE {
			err = u.Repo.OrderTrxItem.DeleteByID(ctx, item.ID)
			if err != nil {
				return err
			}

			trxStockLookupIDsIn := []string{}
			for _, lookup := range item.OrderTrxItemLookups {
				err = u.Repo.OrderTrxItemLookup.DeleteByID(ctx, lookup.ID)
				if err != nil {
					return err
				}
				trxStockLookupIDsIn = append(trxStockLookupIDsIn, lookup.ID)
				trxStockProductIn.Quantity++
			}
			trxStockProductIn.StockLookupIDs = trxStockLookupIDsIn
		} else {
			if item.Action == constant.ACTION_UPDATE {
				_, err = u.Repo.OrderTrxItem.UpdateByID(ctx, item.ID, item)
				if err != nil {
					return err
				}
			}
			if item.Action == constant.ACTION_INSERT {
				_, err = u.Repo.OrderTrxItem.Create(ctx, item)
				if err != nil {
					return err
				}
			}

			trxStockLookupIDsOut := []string{}
			for _, lookup := range item.OrderTrxItemLookups {
				if lookup.Action == constant.ACTION_DELETE {
					err = u.Repo.OrderTrxItemLookup.DeleteByID(ctx, lookup.ID)
					if err != nil {
						return err
					}
					trxStockProductIn.Quantity++
				}
				if lookup.Action == constant.ACTION_UPDATE {
					_, err = u.Repo.OrderTrxItemLookup.UpdateByID(ctx, lookup.ID, lookup)
					if err != nil {
						return err
					}
				}
				if lookup.Action == constant.ACTION_INSERT {
					_, err = u.Repo.OrderTrxItemLookup.Create(ctx, lookup)
					if err != nil {
						return err
					}
					trxStockLookupIDsOut = append(trxStockLookupIDsOut, lookup.ID)
					trxStockProductOut.Quantity++
				}
			}
			trxStockProductOut.StockLookupIDs = trxStockLookupIDsOut
		}

		if trxStockProductIn.Quantity > 0 {
			trxStockProductsIn = append(trxStockProductsIn, trxStockProduct)
		}
		if trxStockProductOut.Quantity > 0 {
			trxStockProductsOut = append(trxStockProductsOut, trxStockProduct)
		}
	}

	// transaction stock
	if len(trxStockProductsIn) > 0 {
		_, err := u.stockUC.TransactionProcess(ctx, dto.TransactionStockRequest{
			TrxType:    constant.TRX_TYPE_OUT,
			OrderTrxID: orderTrx.ID,
			Products:   trxStockProductsIn,
		})
		if err != nil {
			return err
		}
	}
	if len(trxStockProductsOut) > 0 {
		_, err := u.stockUC.TransactionProcess(ctx, dto.TransactionStockRequest{
			TrxType:    constant.TRX_TYPE_IN,
			OrderTrxID: orderTrx.ID,
			Products:   trxStockProductsOut,
		})
		if err != nil {
			return err
		}
	}

	// upsert receipt
	for _, receipt := range orderTrx.OrderTrxReceipts {
		if receipt.Action == constant.ACTION_DELETE {
			err = u.Repo.OrderTrxReceipt.DeleteByID(ctx, receipt.ID)
			if err != nil {
				return err
			}
		}
		if receipt.Action == constant.ACTION_UPDATE {
			_, err = u.Repo.OrderTrxReceipt.UpdateByID(ctx, receipt.ID, receipt)
			if err != nil {
				return err
			}
		}
		if receipt.Action == constant.ACTION_INSERT {
			_, err = u.Repo.OrderTrxReceipt.Create(ctx, receipt)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
