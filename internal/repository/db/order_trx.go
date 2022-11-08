package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	OrderTrx interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.OrderTrxModel]

		// Base
		FindByIDs(ctx context.Context, ids []string, sortBy string) ([]model.OrderTrxModel, error)
		FindByCode(ctx context.Context, code string) (*model.OrderTrxModel, error)
		FindByName(ctx context.Context, name string) (*model.OrderTrxModel, error)
		Create(ctx context.Context, data model.OrderTrxModel) (model.OrderTrxModel, error)
		Creates(ctx context.Context, data []model.OrderTrxModel) ([]model.OrderTrxModel, error)
		UpdateByID(ctx context.Context, id string, data model.OrderTrxModel) (model.OrderTrxModel, error)
		DeleteByID(ctx context.Context, id string) error
		DeleteByIDs(ctx context.Context, ids []string) error
		Count(ctx context.Context) (int64, error)

		// Custom
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.OrderView, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.OrderView, error)
	}

	orderTrx struct {
		Base[model.OrderTrxModel]

		entity     model.OrderTrxModel
		entityName string
	}
)

func NewOrderTrx(conn *gorm.DB) OrderTrx {
	model := model.OrderTrxModel{}
	base := NewBase(conn, model, model.TableName())
	return &orderTrx{
		Base:       base,
		entity:     model,
		entityName: model.TableName(),
	}
}

func (m *orderTrx) Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.OrderView, *abstraction.PaginationInfo, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			order_trxs.*, 
			customers.name as customer_name, 
			suppliers.name as supplier_name,
			pics.name as pic_name
		`).
		Joins(`join users as customers on customers.id = order_trxs.customer_id`).
		Joins(`join users as suppliers on suppliers.id = order_trxs.supplier_id`).
		Joins(`join users as pics on pic.id = order_trxs.pic_id`)

	if search != nil {
		query = query.Where(search.Query, search.Args...)
	}

	m.BuildFilterSort(ctx, query, filterParam)
	info := m.BuildPagination(ctx, query, filterParam.Pagination)

	result := []model.OrderView{}
	err := query.WithContext(ctx).Find(&result).Error

	if err != nil {
		return nil, info, err
	}
	return result, info, nil
}

func (m *orderTrx) FindByID(ctx context.Context, id string) (*model.OrderView, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			order_trxs.*, 
			customers.name as customer_name, 
			suppliers.name as supplier_name,
			pics.name as pic_name
		`).
		Joins(`join users as customers on customers.id = order_trxs.customer_id`).
		Joins(`join users as suppliers on suppliers.id = order_trxs.supplier_id`).
		Joins(`join users as pics on pic.id = order_trxs.pic_id`).
		Where("id", id)
	result := new(model.OrderView)
	err := query.WithContext(ctx).First(result).Error
	if err != nil {
		return nil, m.MaskError(err)
	}
	return result, nil
}
