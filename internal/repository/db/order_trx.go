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
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.OrderTrxModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.OrderTrxModel, error)
		FindByIDs(ctx context.Context, ids []string, sortBy string) ([]model.OrderTrxModel, error)
		FindByCode(ctx context.Context, code string) (*model.OrderTrxModel, error)
		FindByName(ctx context.Context, name string) (*model.OrderTrxModel, error)
		Create(ctx context.Context, data model.OrderTrxModel) (model.OrderTrxModel, error)
		Creates(ctx context.Context, data []model.OrderTrxModel) ([]model.OrderTrxModel, error)
		UpdateByID(ctx context.Context, id string, data model.OrderTrxModel) (model.OrderTrxModel, error)
		DeleteByID(ctx context.Context, id string) error
		DeleteByIDs(ctx context.Context, ids []string) error
		Count(ctx context.Context) (int64, error)
	}

	orderTrx struct {
		Base[model.OrderTrxModel]
	}
)

func NewOrderTrx(conn *gorm.DB) OrderTrx {
	model := model.OrderTrxModel{}
	base := NewBase(conn, model, model.TableName())
	return &orderTrx{
		base,
	}
}
