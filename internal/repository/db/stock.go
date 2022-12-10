package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	Stock interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.StockModel]

		// Base
		FindByID(ctx context.Context, id string) (*model.StockModel, error)
		FindByCode(ctx context.Context, code string) (*model.StockModel, error)
		FindByName(ctx context.Context, name string) (*model.StockModel, error)
		Create(ctx context.Context, data model.StockModel) (model.StockModel, error)
		Creates(ctx context.Context, data []model.StockModel) ([]model.StockModel, error)
		UpdateByID(ctx context.Context, id string, data model.StockModel) (model.StockModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)

		// Custom
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.StockView, *abstraction.PaginationInfo, error)
		FindByProductID(ctx context.Context, productID string) (*model.StockModel, error)
		FindByTotalLessThan(ctx context.Context, total int64) ([]model.StockView, error)
	}

	stock struct {
		Base[model.StockModel]

		entity     model.StockModel
		entityName string
	}
)

func NewStock(conn *gorm.DB) Stock {
	model := model.StockModel{}
	base := NewBase(conn, model, model.TableName())
	return &stock{
		Base:       base,
		entity:     model,
		entityName: model.TableName(),
	}
}

func (m *stock) Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.StockView, *abstraction.PaginationInfo, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			stocks.*,
			products.code as product_code,
			products.name as product_name,
			brands.name as brand_name,
			users.name as supplier_name,
			variants.name as variant_name,
			packets.value as packet_value,
			units.name as unit_name
		`).
		Joins(`join products on products.id = stocks.product_id`).
		Joins(`join brands on brands.id = products.brand_id`).
		Joins(`join users on users.id = brands.supplier_id`).
		Joins(`join variants on variants.id = products.variant_id`).
		Joins(`join packets on packets.id = products.packet_id`).
		Joins(`join units on units.id = packets.unit_id`)

	if search != nil {
		query = query.Where(search.Query, search.Args...)
	}

	m.BuildFilterSort(ctx, query, filterParam)
	info := m.BuildPagination(ctx, query, filterParam.Pagination)

	result := []model.StockView{}
	err := query.WithContext(ctx).Find(&result).Error

	if err != nil {
		return nil, info, err
	}
	return result, info, nil
}

func (m *stock) FindByProductID(ctx context.Context, productID string) (*model.StockModel, error) {
	query := m.GetConn(ctx).Model(m.entity)
	result := new(model.StockModel)
	err := query.WithContext(ctx).Where("product_id", productID).First(result).Error
	if err != nil {
		return nil, m.MaskError(err)
	}
	return result, nil
}

func (m *stock) FindByTotalLessThan(ctx context.Context, total int64) ([]model.StockView, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			stocks.*,
			products.code as product_code,
			products.name as product_name,
			brands.name as brand_name,
			users.name as supplier_name,
			variants.name as variant_name,
			packets.value as packet_value,
			units.name as unit_name
		`).
		Joins(`join products on products.id = stocks.product_id`).
		Joins(`join brands on brands.id = products.brand_id`).
		Joins(`join users on users.id = brands.supplier_id`).
		Joins(`join variants on variants.id = products.variant_id`).
		Joins(`join packets on packets.id = products.packet_id`).
		Joins(`join units on units.id = packets.unit_id`)
	result := []model.StockView{}
	err := query.WithContext(ctx).Where("total < ?", total).Find(&result).Error
	if err != nil {
		return nil, m.MaskError(err)
	}
	return result, nil
}

func (m *stock) Count(ctx context.Context) (int64, error) {
	var count int64
	err := m.GetConn(ctx).Model(m.entity).Raw(`select coalesce(sum(stocks.total_seal+stocks.total_not_seal), 0) as total from stocks`).Scan(&count).Error
	return count, err
}
