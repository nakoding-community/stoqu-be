package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	Product interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.ProductModel]

		// Base
		FindByID(ctx context.Context, id string) (*model.ProductModel, error)
		FindByCode(ctx context.Context, code string) (*model.ProductModel, error)
		FindByName(ctx context.Context, name string) (*model.ProductModel, error)
		Create(ctx context.Context, data model.ProductModel) (model.ProductModel, error)
		Creates(ctx context.Context, data []model.ProductModel) ([]model.ProductModel, error)
		UpdateByID(ctx context.Context, id string, data model.ProductModel) (model.ProductModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)

		// Custom
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.ProductView, *abstraction.PaginationInfo, error)
		FindByBrandVariantPacketID(ctx context.Context, brandID, variantID, packetID string) (*model.ProductModel, error)
	}

	product struct {
		Base[model.ProductModel]

		entity     model.ProductModel
		entityName string
	}
)

func NewProduct(conn *gorm.DB) Product {
	model := model.ProductModel{}
	base := NewBase(conn, model, model.TableName())
	return &product{
		Base:       base,
		entity:     model,
		entityName: model.TableName(),
	}
}

func (m *product) Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.ProductView, *abstraction.PaginationInfo, error) {
	query := m.GetConn(ctx).Model(m.entity).
		Select(`
			products.*, 
			brands.name as brand_name, 
			users.name as supplier_name,
			variants.name as variant_name,
			packets.value as packet_value,
			units.name as unit_name
		`).
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

	result := []model.ProductView{}
	err := query.WithContext(ctx).Find(&result).Error

	if err != nil {
		return nil, info, err
	}
	return result, info, nil
}

func (m *product) FindByBrandVariantPacketID(ctx context.Context, brandID, variantID, packetID string) (*model.ProductModel, error) {
	query := m.GetConn(ctx).Model(m.entity)
	result := new(model.ProductModel)
	err := query.WithContext(ctx).Where("brand_id = ? AND variant_id = ? AND packet_id = ?", brandID, variantID, packetID).First(result).Error
	if err != nil {
		return nil, m.MaskError(err)
	}
	return result, nil
}
