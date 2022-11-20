package usecase

import (
	"context"

	"gitlab.com/stoqu/stoqu-be/internal/config"

	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
)

type Dashboard interface {
	Count(ctx context.Context) (dto.DashboardResponse, error)
}

type dashboard struct {
	Cfg  *config.Configuration
	Repo repository.Factory
}

func NewDashboard(cfg *config.Configuration, f repository.Factory) Dashboard {
	return &dashboard{cfg, f}
}

func (u *dashboard) Count(ctx context.Context) (result dto.DashboardResponse, err error) {
	countBrand, err := u.Repo.Brand.Count(ctx)
	if err != nil {
		return result, err
	}
	countProduct, err := u.Repo.Product.Count(ctx)
	if err != nil {
		return result, err
	}
	countStock, err := u.Repo.Stock.Count(ctx)
	if err != nil {
		return result, err
	}
	countOrder, err := u.Repo.OrderTrx.Count(ctx)
	if err != nil {
		return result, err
	}
	countOrderDaily, err := u.Repo.OrderTrx.CountLastWeek(ctx)
	if err != nil {
		return result, err
	}

	result = dto.DashboardResponse{
		TotalBrand:   countBrand,
		TotalProduct: countProduct,
		TotalStock:   countStock,
		TotalOrder:   countOrder,
		OrderDaily:   countOrderDaily,
	}

	return result, nil
}
