package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	Brand interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.BrandModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.BrandModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.BrandModel, error)
		FindByCode(ctx context.Context, code string) (*model.BrandModel, error)
		FindByName(ctx context.Context, name string) (*model.BrandModel, error)
		Create(ctx context.Context, data model.BrandModel) (model.BrandModel, error)
		Creates(ctx context.Context, data []model.BrandModel) ([]model.BrandModel, error)
		UpdateByID(ctx context.Context, id string, data model.BrandModel) (model.BrandModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)
	}

	brand struct {
		Base[model.BrandModel]
	}
)

func NewBrand(conn *gorm.DB) Brand {
	model := model.BrandModel{}
	base := NewBase(conn, model, model.TableName())
	return &brand{
		base,
	}
}
