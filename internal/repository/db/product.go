package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	Product interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.ProductModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.ProductModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.ProductModel, error)
		FindByCode(ctx context.Context, code string) (*model.ProductModel, error)
		FindByName(ctx context.Context, name string) (*model.ProductModel, error)
		Create(ctx context.Context, data model.ProductModel) (model.ProductModel, error)
		Creates(ctx context.Context, data []model.ProductModel) ([]model.ProductModel, error)
		UpdateByID(ctx context.Context, id string, data model.ProductModel) (model.ProductModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)
	}

	product struct {
		Base[model.ProductModel]
	}
)

func NewProduct(conn *gorm.DB) Product {
	model := model.ProductModel{}
	base := NewBase(conn, model, model.TableName())
	return &product{
		base,
	}
}
