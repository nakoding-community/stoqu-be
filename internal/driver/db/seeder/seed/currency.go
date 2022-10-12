package seed

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	"gitlab.com/stoqu/stoqu-be/pkg/util/str"
	"gorm.io/gorm"
)

type CurrencySeed struct{}

func (s *CurrencySeed) Run(conn *gorm.DB) error {
	trx := conn.Begin()

	currencyNames := []string{"IDR"}
	var currencies []entity.CurrencyModel
	for _, v := range currencyNames {
		currency := entity.CurrencyModel{
			CurrencyEntity: entity.CurrencyEntity{
				Code:   str.GenCode(constant.CODE_CURRENCY_PREFIX),
				Name:   v,
				IsAuto: false,
				Value:  15000,
			},
		}
		currencies = append(currencies, currency)
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

func (s *CurrencySeed) GetTag() string {
	return `currency_seed`
}
