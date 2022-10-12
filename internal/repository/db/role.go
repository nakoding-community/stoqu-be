package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	Role interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.RoleModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.RoleModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.RoleModel, error)
		FindByCode(ctx context.Context, code string) (*model.RoleModel, error)
		FindByName(ctx context.Context, name string) (*model.RoleModel, error)
		Create(ctx context.Context, data model.RoleModel) (model.RoleModel, error)
		Creates(ctx context.Context, data []model.RoleModel) ([]model.RoleModel, error)
		UpdateByID(ctx context.Context, id string, data model.RoleModel) (model.RoleModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)
	}

	role struct {
		Base[model.RoleModel]
	}
)

func NewRole(conn *gorm.DB) Role {
	model := model.RoleModel{}
	base := NewBase(conn, model, model.TableName())
	return &role{
		base,
	}
}
