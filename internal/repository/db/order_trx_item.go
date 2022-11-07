package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	OrderTrxItem interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.OrderTrxItemModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.OrderTrxItemModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.OrderTrxItemModel, error)
		FindByIDs(ctx context.Context, ids []string, sortBy string) ([]model.OrderTrxItemModel, error)
		FindByCode(ctx context.Context, code string) (*model.OrderTrxItemModel, error)
		FindByName(ctx context.Context, name string) (*model.OrderTrxItemModel, error)
		Create(ctx context.Context, data model.OrderTrxItemModel) (model.OrderTrxItemModel, error)
		Creates(ctx context.Context, data []model.OrderTrxItemModel) ([]model.OrderTrxItemModel, error)
		UpdateByID(ctx context.Context, id string, data model.OrderTrxItemModel) (model.OrderTrxItemModel, error)
		DeleteByID(ctx context.Context, id string) error
		DeleteByIDs(ctx context.Context, ids []string) error
		Count(ctx context.Context) (int64, error)
	}

	orderTrxItem struct {
		Base[model.OrderTrxItemModel]
	}
)

func NewOrderTrxItem(conn *gorm.DB) OrderTrxItem {
	model := model.OrderTrxItemModel{}
	base := NewBase(conn, model, model.TableName())
	return &orderTrxItem{
		base,
	}
}
