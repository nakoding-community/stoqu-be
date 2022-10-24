package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	StockTrxItemLookup interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.StockTrxItemLookupModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.StockTrxItemLookupModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.StockTrxItemLookupModel, error)
		FindByIDs(ctx context.Context, ids []string) ([]model.StockTrxItemLookupModel, error)
		FindByCode(ctx context.Context, code string) (*model.StockTrxItemLookupModel, error)
		FindByName(ctx context.Context, name string) (*model.StockTrxItemLookupModel, error)
		Create(ctx context.Context, data model.StockTrxItemLookupModel) (model.StockTrxItemLookupModel, error)
		Creates(ctx context.Context, data []model.StockTrxItemLookupModel) ([]model.StockTrxItemLookupModel, error)
		UpdateByID(ctx context.Context, id string, data model.StockTrxItemLookupModel) (model.StockTrxItemLookupModel, error)
		DeleteByID(ctx context.Context, id string) error
		DeleteByIDs(ctx context.Context, ids []string) error
		Count(ctx context.Context) (int64, error)
	}

	stockTrxItemLookup struct {
		Base[model.StockTrxItemLookupModel]
	}
)

func NewStockTrxItemLookup(conn *gorm.DB) StockTrxItemLookup {
	model := model.StockTrxItemLookupModel{}
	base := NewBase(conn, model, model.TableName())
	return &stockTrxItemLookup{
		base,
	}
}
