package usecase

import (
	"context"
	"fmt"

	"gitlab.com/stoqu/stoqu-be/internal/config"

	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

type Report interface {
	FindOrder(ctx context.Context, filterParam abstraction.Filter) (result dto.OrderReportResponse, pagination abstraction.PaginationInfo, err error)
	FindOrderProduct(ctx context.Context, filterParam abstraction.Filter, query dto.ProductReportQuery) (result []dto.OrderProductReportResponse, pagination abstraction.PaginationInfo, err error)
}

type report struct {
	Cfg  *config.Configuration
	Repo repository.Factory
}

func NewReport(cfg *config.Configuration, f repository.Factory) Report {
	return &report{cfg, f}
}

func (u *report) FindOrder(ctx context.Context, filterParam abstraction.Filter) (result dto.OrderReportResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	orders, info, err := u.Repo.OrderTrx.Find(ctx, filterParam, search)
	if err != nil {
		return dto.OrderReportResponse{}, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, order := range orders {
		result.Orders = append(result.Orders, order)
	}

	totalOrder, err := u.Repo.OrderTrx.Count(ctx)
	if err != nil {
		return dto.OrderReportResponse{}, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	totalIncome, err := u.Repo.OrderTrx.CountIncome(ctx)
	if err != nil {
		return dto.OrderReportResponse{}, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	result.TotalIncome = totalIncome
	result.TotalOrder = totalOrder

	return result, pagination, nil
}

func (u *report) FindOrderProduct(ctx context.Context, filterParam abstraction.Filter, query dto.ProductReportQuery) (result []dto.OrderProductReportResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search

	var (
		orders []entity.OrderViewProduct //!TODO, why this was define as entity.OrderViewProduct not model.OrderViewProduct ?
		info   *abstraction.PaginationInfo
	)
	switch query.Group {
	case constant.GROUP_BY_VARIANT:
		fmt.Println("FindGroupByVariant")
		orders, info, err = u.Repo.OrderTrx.FindGroupByVariant(ctx, filterParam, search)
		break
	case constant.GROUP_BY_PACKET:
		fmt.Println("FindGroupByPacket")
		orders, info, err = u.Repo.OrderTrx.FindGroupByPacket(ctx, filterParam, search)
		break
	default: //constant.GROUP_BY_BRAND
		fmt.Println("FindGroupByBrand")
		orders, info, err = u.Repo.OrderTrx.FindGroupByBrand(ctx, filterParam, search)
		break
	}

	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, order := range orders {
		result = append(result, dto.OrderProductReportResponse{
			OrderViewProduct: order,
		})
	}

	return result, pagination, nil
}
