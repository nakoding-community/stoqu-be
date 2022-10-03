package seed

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	"gitlab.com/stoqu/stoqu-be/pkg/util/str"
	"gorm.io/gorm"
)

type RoleSeed struct{}

func (s *RoleSeed) Run(conn *gorm.DB) error {
	trx := conn.Begin()

	var roles = []entity.RoleModel{
		{
			RoleEntity: entity.RoleEntity{
				Code: constant.CODE_ROLE_PREFIX + str.GenRandStr(constant.LENGTH_CODE),
				Name: "admin",
			},
		},
		{
			RoleEntity: entity.RoleEntity{
				Code: constant.CODE_ROLE_PREFIX + str.GenRandStr(constant.LENGTH_CODE),
				Name: "admin-stock",
			},
		},
		{
			RoleEntity: entity.RoleEntity{
				Code: constant.CODE_ROLE_PREFIX + str.GenRandStr(constant.LENGTH_CODE),
				Name: "admin-order",
			},
		},
		{
			RoleEntity: entity.RoleEntity{
				Code: constant.CODE_ROLE_PREFIX + str.GenRandStr(constant.LENGTH_CODE),
				Name: "customer",
			},
		},
		{
			RoleEntity: entity.RoleEntity{
				Code: constant.CODE_ROLE_PREFIX + str.GenRandStr(constant.LENGTH_CODE),
				Name: "supplier",
			},
		},
	}

	if err := trx.Create(&roles).Error; err != nil {
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

func (s *RoleSeed) GetTag() string {
	return `role_seed`
}
