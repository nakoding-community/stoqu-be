package usecase

import (
	"context"
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

type Unit interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.UnitResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.UnitResponse, error)
	Create(ctx context.Context, payload dto.CreateUnitRequest) (dto.UnitResponse, error)
	Update(ctx context.Context, payload dto.UpdateUnitRequest) (dto.UnitResponse, error)
	Delete(ctx context.Context, payload dto.ByIDRequest) (dto.UnitResponse, error)
}

type unit struct {
	Cfg  *config.Configuration
	Repo repository.Factory
}

func NewUnit(cfg *config.Configuration, f repository.Factory) Unit {
	return &unit{cfg, f}
}

func (u *unit) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.UnitResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := "lower(code) LIKE ? OR lower(name) LIKE ?"
		searchVal := "%" + strings.ToLower(filterParam.Search) + "%"
		search = &abstraction.Search{
			Query: searchQuery,
			Args:  []interface{}{searchVal, searchVal},
		}
	}

	units, info, err := u.Repo.Unit.Find(ctx, filterParam, search)
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

func (u *unit) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.UnitResponse, error) {
	var result dto.UnitResponse

	unit, err := u.Repo.Unit.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.UnitResponse{
		UnitModel: *unit,
	}

	return result, nil
}

func (u *unit) Create(ctx context.Context, payload dto.CreateUnitRequest) (result dto.UnitResponse, err error) {
	var (
		unitID = uuid.New().String()
		unit   = model.UnitModel{
			Entity: model.Entity{
				ID: unitID,
			},
			UnitEntity: model.UnitEntity{
				Code: str.GenCode(constant.CODE_PACKET_PREFIX),
				Name: payload.Name,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Unit.Create(ctx, unit)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.UnitResponse{
		UnitModel: unit,
	}

	return result, nil
}

func (u *unit) Update(ctx context.Context, payload dto.UpdateUnitRequest) (result dto.UnitResponse, err error) {
	var (
		unitData = &model.UnitModel{
			UnitEntity: model.UnitEntity{
				Name: payload.Name,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Unit.UpdateByID(ctx, payload.ID, *unitData)
		if err != nil {
			return err
		}

		unitData, err = u.Repo.Unit.FindByID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.UnitResponse{
		UnitModel: *unitData,
	}

	return result, nil
}

func (u *unit) Delete(ctx context.Context, payload dto.ByIDRequest) (result dto.UnitResponse, err error) {
	var data *model.UnitModel

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		data, err = u.Repo.Unit.FindByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.Unit.DeleteByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.UnitResponse{
		UnitModel: *data,
	}

	return result, nil
}
