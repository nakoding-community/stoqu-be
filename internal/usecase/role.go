package usecase

import (
	"context"

	"gitlab.com/stoqu/stoqu-be/internal/config"

	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

type Role interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.RoleResponse, abstraction.PaginationInfo, error)
}

type role struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewRole(cfg *config.Configuration, f repository.Factory) Role {
	return &role{f, cfg}
}

func (u *role) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.RoleResponse, pagination abstraction.PaginationInfo, err error) {
	roles, info, err := u.Repo.Role.Find(ctx, filterParam, nil)
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
