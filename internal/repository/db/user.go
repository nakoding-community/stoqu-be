package db

import (
	"context"

	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	User interface {
		Base[model.UserModel]
		FindByEmail(ctx context.Context, email string) (*model.UserModel, error)
	}

	user struct {
		Base[model.UserModel]
	}
)

func NewUser(conn *gorm.DB) User {
	model := model.UserModel{}
	base := NewBase(conn, model, model.TableName())
	return &user{
		base,
	}
}

func (m *user) FindByEmail(ctx context.Context, email string) (*model.UserModel, error) {
	query := m.getConn(ctx).Model(model.UserModel{})
	result := new(model.UserModel)
	err := query.Where("email", email).First(result).Error
	if err != nil {
		return nil, m.maskError(err)
	}
	return result, nil
}
