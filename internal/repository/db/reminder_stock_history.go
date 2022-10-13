package db

import (
	"context"
	"time"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	"gitlab.com/stoqu/stoqu-be/pkg/util/ctxval"

	"gorm.io/gorm"
)

type (
	ReminderStockHistory interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.ReminderStockHistoryModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.ReminderStockHistoryModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.ReminderStockHistoryModel, error)
		FindByCode(ctx context.Context, code string) (*model.ReminderStockHistoryModel, error)
		FindByName(ctx context.Context, name string) (*model.ReminderStockHistoryModel, error)
		Create(ctx context.Context, data model.ReminderStockHistoryModel) (model.ReminderStockHistoryModel, error)
		Creates(ctx context.Context, data []model.ReminderStockHistoryModel) ([]model.ReminderStockHistoryModel, error)
		UpdateByID(ctx context.Context, id string, data model.ReminderStockHistoryModel) (model.ReminderStockHistoryModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)

		// Custom
		UpdateBulkRead(ctx context.Context, ids []string) error
		CountUnread(ctx context.Context) (count int64, err error)
	}

	reminderStockHistory struct {
		Base[model.ReminderStockHistoryModel]

		entity     model.ReminderStockHistoryModel
		entityName string
	}
)

func NewReminderStockHistory(conn *gorm.DB) ReminderStockHistory {
	model := model.ReminderStockHistoryModel{}
	base := NewBase(conn, model, model.TableName())
	return &reminderStockHistory{
		Base:       base,
		entity:     model,
		entityName: model.TableName(),
	}
}

func (m *reminderStockHistory) UpdateBulkRead(ctx context.Context, ids []string) error {
	authCtx := ctxval.GetAuthValue(ctx)
	modifiedBy := constant.DB_DEFAULT_SYSTEM
	if authCtx != nil {
		modifiedBy = authCtx.Name
	}

	err := m.GetConn(ctx).Model(&model.ReminderStockHistoryModel{}).Where("id IN ?", ids).Updates(map[string]interface{}{
		"is_read":     true,
		"modified_at": time.Now(),
		"modified_by": modifiedBy,
	}).Error

	return err
}

func (m *reminderStockHistory) CountUnread(ctx context.Context) (count int64, err error) {
	err = m.GetConn(ctx).Model(m.entity).Where("is_read = ?", false).Count(&count).Error
	return
}
