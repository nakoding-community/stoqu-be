package usecase

import (
	"context"

	"gitlab.com/stoqu/stoqu-be/internal/config"

	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

type Unit interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.UnitResponse, abstraction.PaginationInfo, error)
}

type unit struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewUnit(cfg *config.Configuration, f repository.Factory) Unit {
	return &unit{f, cfg}
}

func (u *unit) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.UnitResponse, pagination abstraction.PaginationInfo, err error) {
	units, info, err := u.Repo.Unit.Find(ctx, filterParam, nil)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, unit := range units {
		result = append(result, dto.UnitResponse{
			UnitModel: unit,
		})
	}

	return result, pagination, nil
}
