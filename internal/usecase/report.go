package usecase

import (
	"context"

	"gitlab.com/stoqu/stoqu-be/internal/config"

	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// !TODO: whole report function still need to restructure
type Report interface {
	Order(ctx context.Context, filterParam abstraction.Filter, query dto.OrderReportQuery) ([]dto.OrderReportResponse, abstraction.PaginationInfo, error)
	Product(ctx context.Context, filterParam abstraction.Filter, query dto.ProductReportQuery) ([]dto.ProductReportResponse, abstraction.PaginationInfo, error)
}

type report struct {
	Cfg  *config.Configuration
	Repo repository.Factory
}

func NewReport(cfg *config.Configuration, f repository.Factory) Report {
	return &report{cfg, f}
}

func (u *report) Order(ctx context.Context, filterParam abstraction.Filter, query dto.OrderReportQuery) (result []dto.OrderReportResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := `
			order_trxs.created_at >= ? and order_trxs.created_at <= ? 
		`
		search = &abstraction.Search{
			Query: searchQuery,
			Args: []interface{}{
				query.StartDateTime,
				query.EndDateTime,
			},
		}
	}

	orders, info, err := u.Repo.OrderTrx.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, order := range orders {
		result = append(result, dto.OrderReportResponse{
			OrderView: order,
		})
	}

	return result, pagination, nil
}

func (u *report) Product(ctx context.Context, filterParam abstraction.Filter, query dto.ProductReportQuery) (result []dto.ProductReportResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := `
			order_trxs.created_at >= ? and order_trxs.created_at <= ? 
		`
		search = &abstraction.Search{
			Query: searchQuery,
			Args: []interface{}{
				query.StartDateTime,
				query.EndDateTime,
			},
		}
	}

	orders, info, err := u.Repo.OrderTrx.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, order := range orders {
		result = append(result, dto.ProductReportResponse{
			OrderView: order,
		})
	}

	return result, pagination, nil
}
