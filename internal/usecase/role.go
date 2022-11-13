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

type Role interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.RoleResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.RoleResponse, error)
	Create(ctx context.Context, payload dto.CreateRoleRequest) (dto.RoleResponse, error)
	Update(ctx context.Context, payload dto.UpdateRoleRequest) (dto.RoleResponse, error)
	Delete(ctx context.Context, payload dto.ByIDRequest) (dto.RoleResponse, error)
}

type role struct {
	Cfg  *config.Configuration
	Repo repository.Factory
}

func NewRole(cfg *config.Configuration, f repository.Factory) Role {
	return &role{cfg, f}
}

func (u *role) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.RoleResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := "lower(code) LIKE ? OR lower(name) LIKE ?"
		searchVal := "%" + strings.ToLower(filterParam.Search) + "%"
		search = &abstraction.Search{
			Query: searchQuery,
			Args:  []interface{}{searchVal, searchVal},
		}
	}

	roles, info, err := u.Repo.Role.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, role := range roles {
		result = append(result, dto.RoleResponse{
			RoleModel: role,
		})
	}

	return result, pagination, nil
}

func (u *role) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.RoleResponse, error) {
	var result dto.RoleResponse

	role, err := u.Repo.Role.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.RoleResponse{
		RoleModel: *role,
	}

	return result, nil
}

func (u *role) Create(ctx context.Context, payload dto.CreateRoleRequest) (result dto.RoleResponse, err error) {
	var (
		roleID = uuid.New().String()
		role   = model.RoleModel{
			Entity: model.Entity{
				ID: roleID,
			},
			RoleEntity: model.RoleEntity{
				Code: str.GenCode(constant.CODE_PACKET_PREFIX),
				Name: payload.Name,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Role.Create(ctx, role)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.RoleResponse{
		RoleModel: role,
	}

	return result, nil
}

func (u *role) Update(ctx context.Context, payload dto.UpdateRoleRequest) (result dto.RoleResponse, err error) {
	var (
		roleData = &model.RoleModel{
			RoleEntity: model.RoleEntity{
				Name: payload.Name,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Role.UpdateByID(ctx, payload.ID, *roleData)
		if err != nil {
			return err
		}

		roleData, err = u.Repo.Role.FindByID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.RoleResponse{
		RoleModel: *roleData,
	}

	return result, nil
}

func (u *role) Delete(ctx context.Context, payload dto.ByIDRequest) (result dto.RoleResponse, err error) {
	var data *model.RoleModel

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		data, err = u.Repo.Role.FindByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.Role.DeleteByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.RoleResponse{
		RoleModel: *data,
	}

	return result, nil
}
