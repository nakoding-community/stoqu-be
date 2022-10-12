package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	ProductLookup interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.ProductLookupModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.ProductLookupModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.ProductLookupModel, error)
		FindByCode(ctx context.Context, code string) (*model.ProductLookupModel, error)
		FindByName(ctx context.Context, name string) (*model.ProductLookupModel, error)
		Create(ctx context.Context, data model.ProductLookupModel) (model.ProductLookupModel, error)
		Creates(ctx context.Context, data []model.ProductLookupModel) ([]model.ProductLookupModel, error)
		UpdateByID(ctx context.Context, id string, data model.ProductLookupModel) (model.ProductLookupModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)
	}

	productLookup struct {
		Base[model.ProductLookupModel]
	}
)

func NewProductLookup(conn *gorm.DB) ProductLookup {
	model := model.ProductLookupModel{}
	base := NewBase(conn, model, model.TableName())
	return &productLookup{
		base,
	}
}
