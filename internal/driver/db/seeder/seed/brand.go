package seed

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	"gitlab.com/stoqu/stoqu-be/pkg/util/str"
	"gorm.io/gorm"
)

type BrandSeed struct{}

func (s *BrandSeed) Run(conn *gorm.DB) error {
	trx := conn.Begin()

	brandNames := []string{"maestro", "kahf", "axe"}
	var brands []entity.BrandModel
	for _, v := range brandNames {
		var user entity.UserModel
		if err := trx.Model(&entity.UserModel{}).Joins(`join roles on roles.id = users.role_id`).Find(&user).Where("roles.name = 'supplier'").Error; err != nil {
			return err
		}

		brand := entity.BrandModel{
			BrandEntity: entity.BrandEntity{
				Code:       str.GenCode(constant.CODE_BRAND_PREFIX),
				Name:       v,
				SupplierID: user.ID,
			},
		}
		brands = append(brands, brand)
	}
	if err := trx.Create(&brands).Error; err != nil {
		trx.Rollback()
		logrus.Error(err)
		return err
	}

	variantNames := []string{"chocolate", "apple", "orange", "grape"}
	var variants []entity.VariantModel
	for _, v := range variantNames {
		for _, v2 := range brands {
			variant := entity.VariantModel{
				VariantEntity: entity.VariantEntity{
					Code:       str.GenCode(constant.CODE_VARIANT_PREFIX),
					Name:       v,
					ITL:        v,
					UniqueCode: "",
					BrandID:    v2.ID,
				},
			}
			variants = append(variants, variant)
		}
	}
	if err := trx.Create(&variants).Error; err != nil {
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

func (s *BrandSeed) GetTag() string {
	return `brand_seed`
}
