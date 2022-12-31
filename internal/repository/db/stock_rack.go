package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	StockRack interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.StockRackModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.StockRackModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.StockRackModel, error)
		FindByIDs(ctx context.Context, ids []string, sortBy string) ([]model.StockRackModel, error)
		FindByCode(ctx context.Context, code string) (*model.StockRackModel, error)
		FindByName(ctx context.Context, name string) (*model.StockRackModel, error)
		Create(ctx context.Context, data model.StockRackModel) (model.StockRackModel, error)
		Creates(ctx context.Context, data []model.StockRackModel) ([]model.StockRackModel, error)
		UpdateByID(ctx context.Context, id string, data model.StockRackModel) (model.StockRackModel, error)
		DeleteByID(ctx context.Context, id string) error
		DeleteByIDs(ctx context.Context, ids []string) error
		Count(ctx context.Context) (int64, error)

		// Custom
		FindByStockAndRackID(ctx context.Context, stockID, rackID string) (*model.StockRackModel, error)
		FindByStockID(ctx context.Context, stockID string) ([]model.StockRackModel, error)
	}

	stockRack struct {
		Base[model.StockRackModel]

		entity     model.StockRackModel
		entityName string
	}
)

func NewStockRack(conn *gorm.DB) StockRack {
	model := model.StockRackModel{}
	base := NewBase(conn, model, model.TableName())
	return &stockRack{
		Base:       base,
		entity:     model,
		entityName: model.TableName(),
	}
}

func (m *stockRack) FindByStockAndRackID(ctx context.Context, stockID, rackID string) (*model.StockRackModel, error) {
	query := m.GetConn(ctx).Model(m.entity)
	result := new(model.StockRackModel)
	err := query.WithContext(ctx).Where("stock_id = ? AND rack_id = ?", stockID, rackID).First(result).Error
	if err != nil {
		return nil, m.MaskError(err)
	}
	return result, nil
}

func (m *stockRack) FindByStockID(ctx context.Context, stockID string) ([]model.StockRackModel, error) {
	query := m.GetConn(ctx).Model(m.entity)
	results := []model.StockRackModel{}
	err := query.WithContext(ctx).Where("stock_id = ?", stockID).Find(&results).Error
	if err != nil {
		return nil, m.MaskError(err)
	}

	return results, nil
}
