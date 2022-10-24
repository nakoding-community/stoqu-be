package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	StockLookup interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.StockLookupModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.StockLookupModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.StockLookupModel, error)
		FindByIDs(ctx context.Context, ids []string) ([]model.StockLookupModel, error)
		FindByCode(ctx context.Context, code string) (*model.StockLookupModel, error)
		FindByName(ctx context.Context, name string) (*model.StockLookupModel, error)
		Create(ctx context.Context, data model.StockLookupModel) (model.StockLookupModel, error)
		Creates(ctx context.Context, data []model.StockLookupModel) ([]model.StockLookupModel, error)
		UpdateByID(ctx context.Context, id string, data model.StockLookupModel) (model.StockLookupModel, error)
		DeleteByID(ctx context.Context, id string) error
		DeleteByIDs(ctx context.Context, ids []string) error
		Count(ctx context.Context) (int64, error)
	}

	productLookup struct {
		Base[model.StockLookupModel]
	}
)

func NewStockLookup(conn *gorm.DB) StockLookup {
	model := model.StockLookupModel{}
	base := NewBase(conn, model, model.TableName())
	return &productLookup{
		base,
	}
}
