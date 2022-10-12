package seed

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	"gitlab.com/stoqu/stoqu-be/pkg/util/str"
	"gorm.io/gorm"
)

type ReminderStockSeed struct{}

func (s *ReminderStockSeed) Run(conn *gorm.DB) error {
	trx := conn.Begin()

	reminderStockNames := []string{"daily", "monthly"}
	var currencies []entity.ReminderStockModel
	for _, v := range reminderStockNames {
		reminderStock := entity.ReminderStockModel{
			ReminderStockEntity: entity.ReminderStockEntity{
				Code:     str.GenCode(constant.CODE_REMINDER_STOCK_PREFIX),
				Name:     v,
				MinStock: 10,
			},
		}
		currencies = append(currencies, reminderStock)
	}

	if err := trx.Create(&currencies).Error; err != nil {
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

func (s *ReminderStockSeed) GetTag() string {
	return `reminder_stock_seed`
}
