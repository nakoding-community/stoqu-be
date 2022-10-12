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

	roleNames := []string{"admin", "admin-stock", "admin-order", "customer", "supplier"}
	var roles []entity.RoleModel
	for _, v := range roleNames {
		role := entity.RoleModel{
			RoleEntity: entity.RoleEntity{
				Code: str.GenCode(constant.CODE_ROLE_PREFIX),
				Name: v,
			},
		}
		roles = append(roles, role)
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
