package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
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
		FindDetailByID(ctx context.Context, id string) (*model.OrderDetailView, error)
		CountLastWeek(ctx context.Context) ([]dto.DashboardOrderDailyResponse, error)
		CountIncome(ctx context.Context) (int64, error)
		FindGroupByBrand(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.OrderViewProduct, int64, *abstraction.PaginationInfo, error)
		FindGroupByVariant(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.OrderViewProduct, int64, *abstraction.PaginationInfo, error)
		FindGroupByPacket(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.OrderViewProduct, int64, *abstraction.PaginationInfo, error)
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

func (m *orderTrx) CountIncome(ctx context.Context) (int64, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			sum(order_trxs.price)
		`)

	var count int64
	err := query.WithContext(ctx).Scan(&count).Error

	return count, err
}

func (m *orderTrx) Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.OrderView, *abstraction.PaginationInfo, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			order_trxs.*, 
			customers.name as customer_name, 
			customer_profiles.phone as customer_phone, 
			suppliers.name as supplier_name,
			pics.name as pic_name
		`).
		Joins(`left join users as customers on customers.id = order_trxs.customer_id`).
		Joins(`left join user_profiles as customer_profiles on customer_profiles.user_id = customers.id`).
		Joins(`left join users as suppliers on suppliers.id = order_trxs.supplier_id`).
		Joins(`left join users as pics on pics.id = order_trxs.pic_id`)

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

func (m *orderTrx) FindGroupByBrand(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.OrderViewProduct, int64, *abstraction.PaginationInfo, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			b.id as brand_id, 
			b.name as brand_name, 
			pt.id as packet_id, 
			pt.name as packet_name, 
			SUM(oti.total) as count
		`).
		Joins(`join order_trx_items as oti on oti.order_trx_id = order_trxs.id`).
		Joins(`join products p on p.id = oti.product_id`).
		Joins("join packets pt on pt.id = p.packet_id").
		Joins("join brands b on b.id = p.brand_id").
		Group("pt.id, b.id")

	// count total
	var count int64
	err := query.WithContext(ctx).Count(&count).Error
	if err != nil {
		return nil, count, nil, err
	}

	// find data
	if search != nil {
		query = query.Where(search.Query, search.Args...)
	}

	m.BuildFilterSort(ctx, query, filterParam)
	info := m.BuildPagination(ctx, query, filterParam.Pagination)

	result := []model.OrderViewProduct{}
	err = query.WithContext(ctx).Find(&result).Error
	if err != nil {
		return nil, count, info, err
	}
	return result, count, info, nil
}

func (m *orderTrx) FindGroupByVariant(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.OrderViewProduct, int64, *abstraction.PaginationInfo, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			b.id as brand_id, 
			b.name as brand_name, 
			pt.id as packet_id, 
			pt.name as packet_name, 
			v.id as variant_id, 
			v.name as variant_name, 
			SUM(oti.total) as count
		`).
		Joins(`join order_trx_items as oti on oti.order_trx_id = order_trxs.id`).
		Joins(`join products p on p.id = oti.product_id`).
		Joins("join packets pt on pt.id = p.packet_id").
		Joins("join brands b on b.id = p.brand_id").
		Joins("join variants v on v.id = p.variant_id").
		Group("pt.id, b.id, v.id")

	// count total
	var count int64
	err := query.WithContext(ctx).Count(&count).Error
	if err != nil {
		return nil, count, nil, err
	}

	// find data
	if search != nil {
		query = query.Where(search.Query, search.Args...)
	}

	m.BuildFilterSort(ctx, query, filterParam)
	info := m.BuildPagination(ctx, query, filterParam.Pagination)

	result := []model.OrderViewProduct{}
	err = query.WithContext(ctx).Find(&result).Error

	if err != nil {
		return nil, count, info, err
	}
	return result, count, info, nil
}

func (m *orderTrx) FindGroupByPacket(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.OrderViewProduct, int64, *abstraction.PaginationInfo, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			pt.id as packet_id, 
			pt.name as packet_name, 
			SUM(oti.total) as count
		`).
		Joins(`join order_trx_items as oti on oti.order_trx_id = order_trxs.id`).
		Joins(`join products p on p.id = oti.product_id`).
		Joins("join packets pt on pt.id = p.packet_id").
		Group("pt.id")

	// count total
	var count int64
	err := query.WithContext(ctx).Count(&count).Error
	if err != nil {
		return nil, count, nil, err
	}

	// find data
	if search != nil {
		query = query.Where(search.Query, search.Args...)
	}

	m.BuildFilterSort(ctx, query, filterParam)
	info := m.BuildPagination(ctx, query, filterParam.Pagination)

	result := []model.OrderViewProduct{}
	err = query.WithContext(ctx).Find(&result).Error

	if err != nil {
		return nil, count, info, err
	}
	return result, count, info, nil
}

func (m *orderTrx) FindByID(ctx context.Context, id string) (*model.OrderView, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			order_trxs.*, 
			customers.name as customer_name, 
			suppliers.name as supplier_name,
			pics.name as pic_name
		`).
		Joins(`left join users as customers on customers.id = order_trxs.customer_id`).
		Joins(`left join users as suppliers on suppliers.id = order_trxs.supplier_id`).
		Joins(`left join users as pics on pics.id = order_trxs.pic_id`).
		Where("order_trxs.id = ?", id)
	result := new(model.OrderView)
	err := query.WithContext(ctx).First(result).Error
	if err != nil {
		return nil, m.MaskError(err)
	}
	return result, nil
}

func (m *orderTrx) FindDetailByID(ctx context.Context, id string) (*model.OrderDetailView, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			order_trxs.*, 
			customers.name as customer_name, 
			suppliers.name as supplier_name,
			pics.name as pic_name
		`).
		Joins(`left join users as customers on customers.id = order_trxs.customer_id`).
		Joins(`left join users as suppliers on suppliers.id = order_trxs.supplier_id`).
		Joins(`left join users as pics on pics.id = order_trxs.pic_id`).
		Where("order_trxs.id = ?", id)

	query = query.Preload("OrderTrxItems").Preload("OrderTrxReceipts").Preload("OrderTrxItems.OrderTrxItemLookups")

	result := new(model.OrderDetailView)
	err := query.WithContext(ctx).First(result).Error
	if err != nil {
		return nil, m.MaskError(err)
	}
	return result, nil
}

func (m *orderTrx) CountLastWeek(ctx context.Context) ([]dto.DashboardOrderDailyResponse, error) {
	results := []dto.DashboardOrderDailyResponse{}
	err := m.GetConn(ctx).Model(m.entity).
		Select(`lower( to_char(order_trxs.created_at, 'Day')) as day, count(order_trxs.created_at) as total`).
		Where(`order_trxs.created_at BETWEEN NOW() - INTERVAL '1 week' AND NOW() `).
		Group(`to_char(order_trxs.created_at, 'Day')`).
		Find(&results).Error

	return results, err
}
