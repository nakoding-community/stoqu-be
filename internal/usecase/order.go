package usecase

import (
	"context"
	"strconv"
	"strings"

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
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewOrder(cfg *config.Configuration, f repository.Factory) Order {
	return &order{f, cfg}
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
	//!TODO:need calculation stock in whole function
	mapStockLookups, err := u.buildMapStockLookups(ctx, payload)
	if err != nil {
		return result, err
	}
	orderTrx := payload.ToOrderTrx(mapStockLookups)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		if payload.ID == "" {
			err := u.create(ctx, payload, orderTrx)
			if err != nil {
				return err
			}
		} else {
			err := u.update(ctx, payload, orderTrx)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.OrderUpsertResponse{
		Status: constant.STATUS_SUCCESS,
	}

	return result, nil
}

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
	_, err := u.Repo.OrderTrx.Create(ctx, orderTrx)
	if err != nil {
		return err
	}
	_, err = u.Repo.OrderTrxItem.Creates(ctx, orderTrx.OrderTrxItems)
	if err != nil {
		return err
	}
	for _, item := range orderTrx.OrderTrxItems {
		_, err = u.Repo.OrderTrxItemLookup.Creates(ctx, item.OrderTrxItemLookups)
		if err != nil {
			return err
		}
	}
	_, err = u.Repo.OrderTrxReceipt.Creates(ctx, orderTrx.OrderTrxReceipts)
	if err != nil {
		return err
	}

	return nil
}

func (u *order) update(ctx context.Context, payload dto.UpsertOrderRequest, orderTrx entity.OrderTrxModel) error {
	_, err := u.Repo.OrderTrx.UpdateByID(ctx, orderTrx.ID, orderTrx)
	if err != nil {
		return err
	}

	for _, item := range orderTrx.OrderTrxItems {
		if item.Action == constant.ACTION_DELETE {
			err = u.Repo.OrderTrxItem.DeleteByID(ctx, item.ID)
			if err != nil {
				return err
			}

			for _, lookup := range item.OrderTrxItemLookups {
				err = u.Repo.OrderTrxItemLookup.DeleteByID(ctx, lookup.ID)
				if err != nil {
					return err
				}
			}
		}
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

		for _, lookup := range item.OrderTrxItemLookups {
			if lookup.Action == constant.ACTION_DELETE {
				err = u.Repo.OrderTrxItemLookup.DeleteByID(ctx, lookup.ID)
				if err != nil {
					return err
				}
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
			}
		}
	}

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
