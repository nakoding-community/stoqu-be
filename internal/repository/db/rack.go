package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	Rack interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.RackModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.RackModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.RackModel, error)
		FindByCode(ctx context.Context, code string) (*model.RackModel, error)
		FindByName(ctx context.Context, name string) (*model.RackModel, error)
		Create(ctx context.Context, data model.RackModel) (model.RackModel, error)
		Creates(ctx context.Context, data []model.RackModel) ([]model.RackModel, error)
		UpdateByID(ctx context.Context, id string, data model.RackModel) (model.RackModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)
	}

	rack struct {
		Base[model.RackModel]
	}
)

func NewRack(conn *gorm.DB) Rack {
	model := model.RackModel{}
	base := NewBase(conn, model, model.TableName())
	return &rack{
		base,
	}
}
