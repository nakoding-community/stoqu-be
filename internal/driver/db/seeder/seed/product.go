package seed

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	"gitlab.com/stoqu/stoqu-be/pkg/util/str"
	"gorm.io/gorm"
)

type ProductSeed struct{}

func (s *ProductSeed) Run(conn *gorm.DB) error {
	trx := conn.Begin()

	// brand
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

	// variant
	variantNames := []string{"chocolate", "apple", "orange", "grape"}
	var variants []entity.VariantModel
	for _, v := range brands {
		for _, v2 := range variantNames {
			variant := entity.VariantModel{
				VariantEntity: entity.VariantEntity{
					Code:       str.GenCode(constant.CODE_VARIANT_PREFIX),
					Name:       v2,
					ITL:        v2,
					UniqueCode: "",
					BrandID:    v.ID,
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

	// product
	var packets []entity.PacketModel
	if err := trx.Model(&entity.PacketModel{}).Find(&packets).Error; err != nil {
		return err
	}
	var products []entity.ProductModel
	for _, v := range variants {
		for _, v2 := range packets {
			product := entity.ProductModel{
				ProductEntity: entity.ProductEntity{
					Code:      str.GenCode(constant.CODE_PACKET_PREFIX),
					PriceUSD:  1,
					PriceIDR:  15000,
					BrandID:   v.BrandID,
					VariantID: v.ID,
					PacketID:  v2.ID,
				},
			}
			products = append(products, product)
		}
	}
	if err := trx.Create(&products).Error; err != nil {
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

func (s *ProductSeed) GetTag() string {
	return `product_seed`
}
