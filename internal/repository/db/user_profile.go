package db

import (
	"context"

	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	UserProfile interface {
		Base[model.UserProfileModel]
		FindByUserID(ctx context.Context, userID string) (*model.UserProfileModel, error)
		UpdateByUserID(ctx context.Context, userID string, data model.UserProfileModel) (model.UserProfileModel, error)
		DeleteByUserID(ctx context.Context, userID string) error
	}

	userProfile struct {
		Base[model.UserProfileModel]
	}
)

func NewUserProfile(conn *gorm.DB) UserProfile {
	model := model.UserProfileModel{}
	base := NewBase(conn, model, model.TableName())
	return &userProfile{
		base,
	}
}

func (m *userProfile) FindByUserID(ctx context.Context, userID string) (*model.UserProfileModel, error) {
	query := m.getConn(ctx).Model(model.UserProfileModel{})
	result := new(model.UserProfileModel)
	err := query.Where("user_id", userID).First(result).Error
	if err != nil {
		return nil, m.maskError(err)
	}
	return result, nil
}

func (m *userProfile) UpdateByUserID(ctx context.Context, userID string, data model.UserProfileModel) (model.UserProfileModel, error) {
	err := m.getConn(ctx).Model(&data).Where("user_id = ?", userID).Updates(&data).Error
	return data, err
}

func (m *userProfile) DeleteByUserID(ctx context.Context, userID string) error {
	return m.getConn(ctx).Where("user_id = ?", userID).Delete(&model.UserProfileModel{}).Error
}
