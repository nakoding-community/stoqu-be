package db

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	Role interface {
		Base[model.RoleModel]
	}

	role struct {
		Base[model.RoleModel]
	}
)

func NewRole(conn *gorm.DB) Role {
	model := model.RoleModel{}
	base := NewBase(conn, model, model.TableName())
	return &role{
		base,
	}
}
