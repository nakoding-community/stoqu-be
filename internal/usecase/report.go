package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
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
	FindOrderProduct(ctx context.Context, filterParam abstraction.Filter, query dto.ProductReportQuery) (result dto.OrderProductReportResponse, pagination abstraction.PaginationInfo, err error)
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
		return result, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, order := range orders {
		result.Orders = append(result.Orders, order)
	}

	totalOrder, err := u.Repo.OrderTrx.Count(ctx)
	if err != nil {
		return result, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	totalIncome, err := u.Repo.OrderTrx.CountIncome(ctx)
	if err != nil {
		return result, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	result.TotalIncome = totalIncome
	result.TotalOrder = totalOrder

	return result, pagination, nil
}

func (u *report) FindOrderProduct(ctx context.Context, filterParam abstraction.Filter, query dto.ProductReportQuery) (result dto.OrderProductReportResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search

	var (
		orders []entity.OrderViewProduct
		info   *abstraction.PaginationInfo
		count  int64
	)
	switch query.Group {
	case constant.GROUP_BY_VARIANT:
		logrus.Info("FindGroupByVariant")
		orders, count, info, err = u.Repo.OrderTrx.FindGroupByVariant(ctx, filterParam, search)
	case constant.GROUP_BY_PACKET:
		logrus.Info("FindGroupByPacket")
		orders, count, info, err = u.Repo.OrderTrx.FindGroupByPacket(ctx, filterParam, search)
	default: //constant.GROUP_BY_BRAND
		logrus.Info("FindGroupByBrand")
		orders, count, info, err = u.Repo.OrderTrx.FindGroupByBrand(ctx, filterParam, search)
	}

	if err != nil {
		return result, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, order := range orders {
		result.Orders = append(result.Orders, order)
	}

	result.Total = count

	return result, pagination, nil
}
