package automigration

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.com/stoqu/stoqu-be/internal/config"
	"gitlab.com/stoqu/stoqu-be/internal/driver/db"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gorm.io/gorm"
)

type AutoMigration interface {
	Run()
	SetDb(*gorm.DB)
}

type automigration struct {
	Db       *gorm.DB
	DbModels *[]interface{}
}

func Init(cfg *config.Configuration) {
	var mgConfigurations = map[string]AutoMigration{}
	for _, v := range cfg.Databases {
		if v.DBAutomigrate {
			mgConfigurations[v.DBName] = &automigration{
				DbModels: &[]interface{}{
					&entity.RoleModel{},
					&entity.UserModel{},
					&entity.UserProfileModel{},
					&entity.ReminderStockModel{},
					&entity.ReminderStockHistoryModel{},
					&entity.ConvertionUnitModel{},
					&entity.UnitModel{},
					&entity.PacketModel{},
					&entity.CurrencyModel{},
					&entity.VariantModel{},
					&entity.BrandModel{},
					&entity.ProductModel{},
					&entity.StockModel{},
					&entity.StockLookupModel{},
					&entity.RackModel{},
				},
			}
		}
	}

	for k, v := range mgConfigurations {
		dbConnection, err := db.GetConnection(k)
		if err != nil {
			logrus.Error(fmt.Sprintf("failed to run automigration, database not found %s", k))
		} else {
			v.SetDb(dbConnection)
			v.Run()
			logrus.Info(fmt.Sprintf("successfully run automigration for database %s", k))
		}
	}
}

func (m *automigration) Run() {
	m.Db.AutoMigrate(*m.DbModels...)
}

func (m *automigration) SetDb(db *gorm.DB) {
	m.Db = db
}
