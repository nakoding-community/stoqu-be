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

type Rack interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.RackResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.RackResponse, error)
	Create(ctx context.Context, payload dto.CreateRackRequest) (dto.RackResponse, error)
	Update(ctx context.Context, payload dto.UpdateRackRequest) (dto.RackResponse, error)
	Delete(ctx context.Context, payload dto.ByIDRequest) (dto.RackResponse, error)
}

type rack struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewRack(cfg *config.Configuration, f repository.Factory) Rack {
	return &rack{f, cfg}
}

func (u *rack) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.RackResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := "lower(code) LIKE ? OR lower(name) LIKE ?"
		searchVal := "%" + strings.ToLower(filterParam.Search) + "%"
		search = &abstraction.Search{
			Query: searchQuery,
			Args:  []interface{}{searchVal, searchVal},
		}
	}

	racks, info, err := u.Repo.Rack.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, rack := range racks {
		result = append(result, dto.RackResponse{
			RackModel: rack,
		})
	}

	return result, pagination, nil
}

func (u *rack) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.RackResponse, error) {
	var result dto.RackResponse

	rack, err := u.Repo.Rack.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.RackResponse{
		RackModel: *rack,
	}

	return result, nil
}

func (u *rack) Create(ctx context.Context, payload dto.CreateRackRequest) (result dto.RackResponse, err error) {
	var (
		rackID = uuid.New().String()
		rack   = model.RackModel{
			Entity: model.Entity{
				ID: rackID,
			},
			RackEntity: model.RackEntity{
				Code: str.GenCode(constant.CODE_RACK_PREFIX),
				Name: payload.Name,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Rack.Create(ctx, rack)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.RackResponse{
		RackModel: rack,
	}

	return result, nil
}

func (u *rack) Update(ctx context.Context, payload dto.UpdateRackRequest) (result dto.RackResponse, err error) {
	var (
		rackData = &model.RackModel{
			RackEntity: model.RackEntity{
				Name: payload.Name,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Rack.UpdateByID(ctx, payload.ID, *rackData)
		if err != nil {
			return err
		}

		rackData, err = u.Repo.Rack.FindByID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.RackResponse{
		RackModel: *rackData,
	}

	return result, nil
}

func (u *rack) Delete(ctx context.Context, payload dto.ByIDRequest) (result dto.RackResponse, err error) {
	var data *model.RackModel

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		data, err = u.Repo.Rack.FindByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.Rack.DeleteByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.RackResponse{
		RackModel: *data,
	}

	return result, nil
}
