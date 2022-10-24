package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	StockTrxItem interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.StockTrxItemModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.StockTrxItemModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.StockTrxItemModel, error)
		FindByIDs(ctx context.Context, ids []string) ([]model.StockTrxItemModel, error)
		FindByCode(ctx context.Context, code string) (*model.StockTrxItemModel, error)
		FindByName(ctx context.Context, name string) (*model.StockTrxItemModel, error)
		Create(ctx context.Context, data model.StockTrxItemModel) (model.StockTrxItemModel, error)
		Creates(ctx context.Context, data []model.StockTrxItemModel) ([]model.StockTrxItemModel, error)
		UpdateByID(ctx context.Context, id string, data model.StockTrxItemModel) (model.StockTrxItemModel, error)
		DeleteByID(ctx context.Context, id string) error
		DeleteByIDs(ctx context.Context, ids []string) error
		Count(ctx context.Context) (int64, error)
	}

	stockTrxItem struct {
		Base[model.StockTrxItemModel]
	}
)

func NewStockTrxItem(conn *gorm.DB) StockTrxItem {
	model := model.StockTrxItemModel{}
	base := NewBase(conn, model, model.TableName())
	return &stockTrxItem{
		base,
	}
}
