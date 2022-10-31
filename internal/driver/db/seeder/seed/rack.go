package seed

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	"gitlab.com/stoqu/stoqu-be/pkg/util/str"
	"gorm.io/gorm"
)

type RackSeed struct{}

func (s *RackSeed) Run(conn *gorm.DB) error {
	trx := conn.Begin()

	rackNames := []string{"Rack 1", "Rack 2", "Rack 3"}
	var racks []entity.RackModel
	for _, v := range rackNames {
		rack := entity.RackModel{
			RackEntity: entity.RackEntity{
				Code: str.GenCode(constant.CODE_RACK_PREFIX),
				Name: v,
			},
		}
		racks = append(racks, rack)
	}

	if err := trx.Create(&racks).Error; err != nil {
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

func (s *RackSeed) GetTag() string {
	return `rack_seed`
}
