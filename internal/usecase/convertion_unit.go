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

type ConvertionUnit interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.ConvertionUnitResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.ConvertionUnitResponse, error)
	Create(ctx context.Context, payload dto.CreateConvertionUnitRequest) (dto.ConvertionUnitResponse, error)
	Update(ctx context.Context, payload dto.UpdateConvertionUnitRequest) (dto.ConvertionUnitResponse, error)
	Delete(ctx context.Context, payload dto.ByIDRequest) (dto.ConvertionUnitResponse, error)
}

type convertionUnit struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewConvertionUnit(cfg *config.Configuration, f repository.Factory) ConvertionUnit {
	return &convertionUnit{f, cfg}
}

func (u *convertionUnit) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.ConvertionUnitResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := "lower(code) LIKE ? OR lower(name) LIKE ? OR value_convertion = ?"
		searchVal := "%" + strings.ToLower(filterParam.Search) + "%"
		searchValFloat, _ := strconv.ParseFloat(filterParam.Search, 64)
		search = &abstraction.Search{
			Query: searchQuery,
			Args:  []interface{}{searchVal, searchVal, searchValFloat},
		}
	}

	convertionUnits, info, err := u.Repo.ConvertionUnit.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, convertionUnit := range convertionUnits {
		unitOrigin, err := u.Repo.Unit.FindByID(ctx, convertionUnit.UnitOriginID)
		if err != nil {
			return nil, pagination, err
		}

		unitDestination, err := u.Repo.Unit.FindByID(ctx, convertionUnit.UnitDestinationID)
		if err != nil {
			return nil, pagination, err
		}

		result = append(result, dto.ConvertionUnitResponse{
			ID:          convertionUnit.ID,
			Origin:      unitOrigin.Name,
			Destination: unitDestination.Name,
			Total:       convertionUnit.ValueConvertion,
		})
	}

	return result, pagination, nil
}

func (u *convertionUnit) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.ConvertionUnitResponse, error) {
	var result dto.ConvertionUnitResponse

	convertionUnit, err := u.Repo.ConvertionUnit.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	unitOrigin, err := u.Repo.Unit.FindByID(ctx, convertionUnit.UnitOriginID)
	if err != nil {
		return result, err
	}

	unitDestination, err := u.Repo.Unit.FindByID(ctx, convertionUnit.UnitDestinationID)
	if err != nil {
		return result, err
	}

	result = dto.ConvertionUnitResponse{
		ID:          convertionUnit.ID,
		Origin:      unitOrigin.Name,
		Destination: unitDestination.Name,
		Total:       convertionUnit.ValueConvertion,
	}

	return result, nil
}

func (u *convertionUnit) Create(ctx context.Context, payload dto.CreateConvertionUnitRequest) (result dto.ConvertionUnitResponse, err error) {
	var (
		convertionUnitID = uuid.New().String()
		convertionUnit   = model.ConvertionUnitModel{
			Entity: model.Entity{
				ID: convertionUnitID,
			},
			ConvertionUnitEntity: model.ConvertionUnitEntity{
				Code:              str.GenCode(constant.CODE_CONVERTION_UNIT_PREFIX),
				UnitOriginID:      payload.Origin,
				UnitDestinationID: payload.Destination,
				ValueConvertion:   payload.Total,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.ConvertionUnit.Create(ctx, convertionUnit)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.ConvertionUnitResponse{
		ID:          convertionUnitID,
		Origin:      payload.Origin,
		Destination: payload.Destination,
		Total:       payload.Total,
	}

	return result, nil
}

func (u *convertionUnit) Update(ctx context.Context, payload dto.UpdateConvertionUnitRequest) (result dto.ConvertionUnitResponse, err error) {
	var (
		convertionUnitData = &model.ConvertionUnitModel{
			ConvertionUnitEntity: model.ConvertionUnitEntity{
				UnitOriginID:      payload.Origin,
				UnitDestinationID: payload.Destination,
				ValueConvertion:   payload.Total,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.ConvertionUnit.UpdateByID(ctx, payload.ID, *convertionUnitData)
		if err != nil {
			return err
		}

		convertionUnitData, err = u.Repo.ConvertionUnit.FindByID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.ConvertionUnitResponse{
		ID:          convertionUnitData.ID,
		Origin:      convertionUnitData.UnitOriginID,
		Destination: convertionUnitData.UnitDestinationID,
		Total:       convertionUnitData.ValueConvertion,
	}

	return result, nil
}

func (u *convertionUnit) Delete(ctx context.Context, payload dto.ByIDRequest) (result dto.ConvertionUnitResponse, err error) {
	var data *model.ConvertionUnitModel

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		data, err = u.Repo.ConvertionUnit.FindByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.ConvertionUnit.DeleteByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.ConvertionUnitResponse{
		ID:          data.ID,
		Origin:      data.UnitOriginID,
		Destination: data.UnitDestinationID,
		Total:       data.ValueConvertion,
	}

	return result, nil
}
