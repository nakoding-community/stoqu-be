package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	ReminderStock interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.ReminderStockModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.ReminderStockModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.ReminderStockModel, error)
		FindByCode(ctx context.Context, code string) (*model.ReminderStockModel, error)
		FindByName(ctx context.Context, name string) (*model.ReminderStockModel, error)
		Create(ctx context.Context, data model.ReminderStockModel) (model.ReminderStockModel, error)
		Creates(ctx context.Context, data []model.ReminderStockModel) ([]model.ReminderStockModel, error)
		UpdateByID(ctx context.Context, id string, data model.ReminderStockModel) (model.ReminderStockModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)
	}

	reminderStock struct {
		Base[model.ReminderStockModel]
	}
)

func NewReminderStock(conn *gorm.DB) ReminderStock {
	model := model.ReminderStockModel{}
	base := NewBase(conn, model, model.TableName())
	return &reminderStock{
		base,
	}
}
