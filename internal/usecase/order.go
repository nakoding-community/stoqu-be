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
	var (
		orderID = uuid.New().String()
		order   = model.OrderTrxModel{
			Entity: model.Entity{
				ID: orderID,
			},
			//!TODO: payload mapping here
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.OrderTrx.Create(ctx, order)
		if err != nil {
			return err
		}

		//!TODO: logic here

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.OrderUpsertResponse{
		Status: constant.STATUS_SUCCESS,
	}

	return result, nil
}
