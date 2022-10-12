package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	UserProfile interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.UserProfileModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.UserProfileModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.UserProfileModel, error)
		FindByCode(ctx context.Context, code string) (*model.UserProfileModel, error)
		FindByName(ctx context.Context, name string) (*model.UserProfileModel, error)
		Create(ctx context.Context, data model.UserProfileModel) (model.UserProfileModel, error)
		Creates(ctx context.Context, data []model.UserProfileModel) ([]model.UserProfileModel, error)
		UpdateByID(ctx context.Context, id string, data model.UserProfileModel) (model.UserProfileModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)

		// Custom
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
	query := m.GetConn(ctx).Model(model.UserProfileModel{})
	result := new(model.UserProfileModel)
	err := query.Where("user_id", userID).First(result).Error
	if err != nil {
		return nil, m.MaskError(err)
	}
	return result, nil
}

func (m *userProfile) UpdateByUserID(ctx context.Context, userID string, data model.UserProfileModel) (model.UserProfileModel, error) {
	err := m.GetConn(ctx).Model(&data).Where("user_id = ?", userID).Updates(&data).Error
	return data, err
}

func (m *userProfile) DeleteByUserID(ctx context.Context, userID string) error {
	return m.GetConn(ctx).Where("user_id = ?", userID).Delete(&model.UserProfileModel{}).Error
}
