package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	Unit interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.UnitModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.UnitModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.UnitModel, error)
		FindByCode(ctx context.Context, code string) (*model.UnitModel, error)
		FindByName(ctx context.Context, name string) (*model.UnitModel, error)
		Create(ctx context.Context, data model.UnitModel) (model.UnitModel, error)
		Creates(ctx context.Context, data []model.UnitModel) ([]model.UnitModel, error)
		UpdateByID(ctx context.Context, id string, data model.UnitModel) (model.UnitModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)
	}

	unit struct {
		Base[model.UnitModel]
	}
)

func NewUnit(conn *gorm.DB) Unit {
	model := model.UnitModel{}
	base := NewBase(conn, model, model.TableName())
	return &unit{
		base,
	}
}
