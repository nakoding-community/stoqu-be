package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	StockTrx interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.StockTrxModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.StockTrxModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.StockTrxModel, error)
		FindByIDs(ctx context.Context, ids []string) ([]model.StockTrxModel, error)
		FindByCode(ctx context.Context, code string) (*model.StockTrxModel, error)
		FindByName(ctx context.Context, name string) (*model.StockTrxModel, error)
		Create(ctx context.Context, data model.StockTrxModel) (model.StockTrxModel, error)
		Creates(ctx context.Context, data []model.StockTrxModel) ([]model.StockTrxModel, error)
		UpdateByID(ctx context.Context, id string, data model.StockTrxModel) (model.StockTrxModel, error)
		DeleteByID(ctx context.Context, id string) error
		DeleteByIDs(ctx context.Context, ids []string) error
		Count(ctx context.Context) (int64, error)
	}

	stockTrx struct {
		Base[model.StockTrxModel]
	}
)

func NewStockTrx(conn *gorm.DB) StockTrx {
	model := model.StockTrxModel{}
	base := NewBase(conn, model, model.TableName())
	return &stockTrx{
		base,
	}
}
