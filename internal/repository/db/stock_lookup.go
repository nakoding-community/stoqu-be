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
		FindByID(ctx context.Context, id string) (*model.StockLookupModel, error)
		FindByIDs(ctx context.Context, ids []string, sortBy string) ([]model.StockLookupModel, error)
		FindByCode(ctx context.Context, code string) (*model.StockLookupModel, error)
		FindByName(ctx context.Context, name string) (*model.StockLookupModel, error)
		Create(ctx context.Context, data model.StockLookupModel) (model.StockLookupModel, error)
		Creates(ctx context.Context, data []model.StockLookupModel) ([]model.StockLookupModel, error)
		UpdateByID(ctx context.Context, id string, data model.StockLookupModel) (model.StockLookupModel, error)
		DeleteByID(ctx context.Context, id string) error
		DeleteByIDs(ctx context.Context, ids []string) error
		Count(ctx context.Context) (int64, error)

		// Custom
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.StockLookupView, *abstraction.PaginationInfo, error)
		SumByIDs(ctx context.Context, ids []string) (*model.StockLookupSum, error)
	}

	stockLookup struct {
		Base[model.StockLookupModel]

		entity     model.StockLookupModel
		entityName string
	}
)

func NewStockLookup(conn *gorm.DB) StockLookup {
	model := model.StockLookupModel{}
	base := NewBase(conn, model, model.TableName())
	return &stockLookup{
		Base:       base,
		entity:     model,
		entityName: model.TableName(),
	}
}

func (m *stockLookup) Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.StockLookupView, *abstraction.PaginationInfo, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			stock_lookups.*,
			stock_racks.rack_id as rack_id,
			stocks.product_id as product_id
		`).
		Joins(`join stock_racks on stock_racks.id = stock_lookups.stock_rack_id`).
		Joins(`join stocks on stocks.id = stock_racks.stock_id`)

	if search != nil {
		query = query.Where(search.Query, search.Args...)
	}

	m.BuildFilterSort(ctx, query, filterParam)
	info := m.BuildPagination(ctx, query, filterParam.Pagination)

	result := []model.StockLookupView{}
	err := query.WithContext(ctx).Find(&result).Error

	if err != nil {
		return nil, info, m.MaskError(err)
	}
	return result, info, nil
}

func (m *stockLookup) SumByIDs(ctx context.Context, ids []string) (*model.StockLookupSum, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			sum(remaining_value) as remaining_value 
		`).
		Where(`id in ?`, ids).
		Group(`id`)

	result := new(model.StockLookupSum)
	err := query.WithContext(ctx).First(result).Error
	if err != nil {
		return nil, m.MaskError(err)
	}
	return result, nil
}
