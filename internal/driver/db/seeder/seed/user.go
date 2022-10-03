package seed

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gorm.io/gorm"
)

type UserSeed struct{}

func (s *UserSeed) Run(conn *gorm.DB) error {
	trx := conn.Begin()

	var roles []entity.RoleModel
	if err := trx.Model(&entity.RoleModel{}).Find(&roles).Error; err != nil {
		log.Println(err)
		return err
	}

	for _, role := range roles {
		user := entity.UserModel{
			Entity: entity.Entity{
				ID: uuid.New().String(),
			},
			UserEntity: entity.UserEntity{
				Name:     role.Name,
				Email:    fmt.Sprintf(`%s@gmail.com`, role.Name),
				Password: role.Name,
				RoleID:   role.ID,
			},
		}
		if err := trx.Create(&user).Error; err != nil {
			trx.Rollback()
			logrus.Error(err)
			return err
		}

		userProfile := entity.UserProfileModel{
			UserProfileEntity: entity.UserProfileEntity{
				Address: "xxx",
				Phone:   "021",
				UserID:  user.ID,
			},
		}
		if err := trx.Create(&userProfile).Error; err != nil {
			trx.Rollback()
			logrus.Error(err)
			return err
		}
	}

	if err := trx.Commit().Error; err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (s *UserSeed) GetTag() string {
	return `user_seed`
}
