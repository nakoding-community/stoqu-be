package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	OrderTrxItemLookup interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.OrderTrxItemLookupModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.OrderTrxItemLookupModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.OrderTrxItemLookupModel, error)
		FindByIDs(ctx context.Context, ids []string, sortBy string) ([]model.OrderTrxItemLookupModel, error)
		FindByCode(ctx context.Context, code string) (*model.OrderTrxItemLookupModel, error)
		FindByName(ctx context.Context, name string) (*model.OrderTrxItemLookupModel, error)
		Create(ctx context.Context, data model.OrderTrxItemLookupModel) (model.OrderTrxItemLookupModel, error)
		Creates(ctx context.Context, data []model.OrderTrxItemLookupModel) ([]model.OrderTrxItemLookupModel, error)
		UpdateByID(ctx context.Context, id string, data model.OrderTrxItemLookupModel) (model.OrderTrxItemLookupModel, error)
		DeleteByID(ctx context.Context, id string) error
		DeleteByIDs(ctx context.Context, ids []string) error
		Count(ctx context.Context) (int64, error)
	}

	orderTrxItemLookup struct {
		Base[model.OrderTrxItemLookupModel]
	}
)

func NewOrderTrxItemLookup(conn *gorm.DB) OrderTrxItemLookup {
	model := model.OrderTrxItemLookupModel{}
	base := NewBase(conn, model, model.TableName())
	return &orderTrxItemLookup{
		base,
	}
}
