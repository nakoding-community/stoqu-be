package seed

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	"gitlab.com/stoqu/stoqu-be/pkg/util/str"
	"gorm.io/gorm"
)

type PacketSeed struct{}

func (s *PacketSeed) Run(conn *gorm.DB) error {
	trx := conn.Begin()

	// unit
	unitNames := []string{"KG", "L", "ML"}
	var units []entity.UnitModel
	for _, v := range unitNames {
		unit := entity.UnitModel{
			Entity: entity.Entity{
				ID: uuid.New().String(),
			},
			UnitEntity: entity.UnitEntity{
				Code: str.GenCode(constant.CODE_UNIT_PREFIX),
				Name: v,
			},
		}
		units = append(units, unit)
	}
	if err := trx.Create(&units).Error; err != nil {
		trx.Rollback()
		logrus.Error(err)
		return err
	}

	// convertion unit
	var convertionUnits = []entity.ConvertionUnitModel{
		{
			Entity: entity.Entity{
				ID: uuid.New().String(),
			},
			ConvertionUnitEntity: entity.ConvertionUnitEntity{
				Code:              str.GenCode(constant.CODE_CONVERTION_UNIT_PREFIX),
				Name:              "KG -> L",
				UnitOriginID:      units[0].ID,
				UnitDestinationID: units[1].ID,
				ValueConvertion:   1,
			},
		},
		{
			Entity: entity.Entity{
				ID: uuid.New().String(),
			},
			ConvertionUnitEntity: entity.ConvertionUnitEntity{
				Code:              str.GenCode(constant.CODE_CONVERTION_UNIT_PREFIX),
				Name:              "KG -> ML",
				UnitOriginID:      units[0].ID,
				UnitDestinationID: units[2].ID,
				ValueConvertion:   100,
			},
		},
		{
			Entity: entity.Entity{
				ID: uuid.New().String(),
			},
			ConvertionUnitEntity: entity.ConvertionUnitEntity{
				Code:              str.GenCode(constant.CODE_CONVERTION_UNIT_PREFIX),
				Name:              "L -> KG",
				UnitOriginID:      units[1].ID,
				UnitDestinationID: units[0].ID,
				ValueConvertion:   1,
			},
		},

		{
			Entity: entity.Entity{
				ID: uuid.New().String(),
			},
			ConvertionUnitEntity: entity.ConvertionUnitEntity{
				Code:              str.GenCode(constant.CODE_CONVERTION_UNIT_PREFIX),
				Name:              "L -> ML",
				UnitOriginID:      units[1].ID,
				UnitDestinationID: units[2].ID,
				ValueConvertion:   100,
			},
		},
	}
	if err := trx.Create(&convertionUnits).Error; err != nil {
		trx.Rollback()
		logrus.Error(err)
		return err
	}

	// packet
	pkgValues := []int{1, 5, 10, 25}
	var pkgs []entity.PacketModel
	for _, v := range pkgValues {
		for _, v2 := range units {
			pkg := entity.PacketModel{
				Entity: entity.Entity{
					ID: uuid.New().String(),
				},
				PacketEntity: entity.PacketEntity{
					Code:   str.GenCode(constant.CODE_PACKET_PREFIX),
					Name:   fmt.Sprintf(`%d %s`, v, v2.Name),
					UnitID: v2.ID,
					Value:  int64(v),
				},
			}
			pkgs = append(pkgs, pkg)
		}
	}
	if err := trx.Create(&pkgs).Error; err != nil {
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

func (s *PacketSeed) GetTag() string {
	return `packet_seed`
}
