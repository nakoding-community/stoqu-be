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
		FindByID(ctx context.Context, id string) (*model.StockTrxModel, error)
		FindByIDs(ctx context.Context, ids []string, sortBy string) ([]model.StockTrxModel, error)
		FindByCode(ctx context.Context, code string) (*model.StockTrxModel, error)
		FindByName(ctx context.Context, name string) (*model.StockTrxModel, error)
		Create(ctx context.Context, data model.StockTrxModel) (model.StockTrxModel, error)
		Creates(ctx context.Context, data []model.StockTrxModel) ([]model.StockTrxModel, error)
		UpdateByID(ctx context.Context, id string, data model.StockTrxModel) (model.StockTrxModel, error)
		DeleteByID(ctx context.Context, id string) error
		DeleteByIDs(ctx context.Context, ids []string) error
		Count(ctx context.Context) (int64, error)

		// Custom
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.StockTrxView, *abstraction.PaginationInfo, error)
	}

	stockTrx struct {
		Base[model.StockTrxModel]

		entity     model.StockTrxModel
		entityName string
	}
)

func NewStockTrx(conn *gorm.DB) StockTrx {
	model := model.StockTrxModel{}
	base := NewBase(conn, model, model.TableName())
	return &stockTrx{
		Base:       base,
		entity:     model,
		entityName: model.TableName(),
	}
}

func (m *stockTrx) Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.StockTrxView, *abstraction.PaginationInfo, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			stock_trxs.*,
			order_trxs.code as order_code
		`).
		Joins(`left join order_trxs on order_trxs.id = stock_trxs.order_trx_id`)

	if search != nil {
		query = query.Where(search.Query, search.Args...)
	}

	m.BuildFilterSort(ctx, query, filterParam)
	info := m.BuildPagination(ctx, query, filterParam.Pagination)

	result := []model.StockTrxView{}
	err := query.WithContext(ctx).Find(&result).Error

	if err != nil {
		return nil, info, err
	}
	return result, info, nil
}
