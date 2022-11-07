package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	OrderTrxStatus interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.OrderTrxStatusModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.OrderTrxStatusModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.OrderTrxStatusModel, error)
		FindByIDs(ctx context.Context, ids []string, sortBy string) ([]model.OrderTrxStatusModel, error)
		FindByCode(ctx context.Context, code string) (*model.OrderTrxStatusModel, error)
		FindByName(ctx context.Context, name string) (*model.OrderTrxStatusModel, error)
		Create(ctx context.Context, data model.OrderTrxStatusModel) (model.OrderTrxStatusModel, error)
		Creates(ctx context.Context, data []model.OrderTrxStatusModel) ([]model.OrderTrxStatusModel, error)
		UpdateByID(ctx context.Context, id string, data model.OrderTrxStatusModel) (model.OrderTrxStatusModel, error)
		DeleteByID(ctx context.Context, id string) error
		DeleteByIDs(ctx context.Context, ids []string) error
		Count(ctx context.Context) (int64, error)
	}

	orderTrxStatus struct {
		Base[model.OrderTrxStatusModel]
	}
)

func NewOrderTrxStatus(conn *gorm.DB) OrderTrxStatus {
	model := model.OrderTrxStatusModel{}
	base := NewBase(conn, model, model.TableName())
	return &orderTrxStatus{
		base,
	}
}
