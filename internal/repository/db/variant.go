package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	Variant interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.VariantModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.VariantModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.VariantModel, error)
		FindByCode(ctx context.Context, code string) (*model.VariantModel, error)
		FindByName(ctx context.Context, name string) (*model.VariantModel, error)
		Create(ctx context.Context, data model.VariantModel) (model.VariantModel, error)
		Creates(ctx context.Context, data []model.VariantModel) ([]model.VariantModel, error)
		UpdateByID(ctx context.Context, id string, data model.VariantModel) (model.VariantModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)
	}

	variant struct {
		Base[model.VariantModel]
	}
)

func NewVariant(conn *gorm.DB) Variant {
	model := model.VariantModel{}
	base := NewBase(conn, model, model.TableName())
	return &variant{
		base,
	}
}
