package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	OrderTrxReceipt interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.OrderTrxReceiptModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.OrderTrxReceiptModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.OrderTrxReceiptModel, error)
		FindByIDs(ctx context.Context, ids []string, sortBy string) ([]model.OrderTrxReceiptModel, error)
		FindByCode(ctx context.Context, code string) (*model.OrderTrxReceiptModel, error)
		FindByName(ctx context.Context, name string) (*model.OrderTrxReceiptModel, error)
		Create(ctx context.Context, data model.OrderTrxReceiptModel) (model.OrderTrxReceiptModel, error)
		Creates(ctx context.Context, data []model.OrderTrxReceiptModel) ([]model.OrderTrxReceiptModel, error)
		UpdateByID(ctx context.Context, id string, data model.OrderTrxReceiptModel) (model.OrderTrxReceiptModel, error)
		DeleteByID(ctx context.Context, id string) error
		DeleteByIDs(ctx context.Context, ids []string) error
		Count(ctx context.Context) (int64, error)
	}

	orderTrxReceipt struct {
		Base[model.OrderTrxReceiptModel]
	}
)

func NewOrderTrxReceipt(conn *gorm.DB) OrderTrxReceipt {
	model := model.OrderTrxReceiptModel{}
	base := NewBase(conn, model, model.TableName())
	return &orderTrxReceipt{
		base,
	}
}
