package seed

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gorm.io/gorm"
)

type ReminderStockHistorySeed struct{}

func (s *ReminderStockHistorySeed) Run(conn *gorm.DB) error {
	trx := conn.Begin()

	var reminderStockHistories []entity.ReminderStockHistoryModel
	for i := 0; i < 5; i++ {
		reminderStock := entity.ReminderStockHistoryModel{
			ReminderStockHistoryEntity: entity.ReminderStockHistoryEntity{
				Title: fmt.Sprintf(`%s %d`, "This is title reminder stock history", i),
				Body:  fmt.Sprintf(`%s %d`, "This is body reminder stock history", i),
			},
		}
		reminderStockHistories = append(reminderStockHistories, reminderStock)
	}

	if err := trx.Create(&reminderStockHistories).Error; err != nil {
		trx.Rollback()
		logrus.Error(err)
		return err
	}

	if err := trx.Commit().Error; err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (s *ReminderStockHistorySeed) GetTag() string {
	return `reminder_stock_history_seed`
}
