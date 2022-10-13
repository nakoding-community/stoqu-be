package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gitlab.com/stoqu/stoqu-be/internal/config"

	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
	"gitlab.com/stoqu/stoqu-be/pkg/util/str"
	"gitlab.com/stoqu/stoqu-be/pkg/util/trxmanager"
)

type User interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.UserResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.UserResponse, error)
	Create(ctx context.Context, payload dto.CreateUserRequest) (dto.UserResponse, error)
	Update(ctx context.Context, payload dto.UpdateUserRequest) (dto.UserResponse, error)
	Delete(ctx context.Context, payload dto.ByIDRequest) (dto.UserResponse, error)
}

type user struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewUser(cfg *config.Configuration, f repository.Factory) User {
	return &user{f, cfg}
}

func (u *user) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.UserResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := "lower(name) LIKE ? OR lower(email) Like ?"
		searchVal := "%" + strings.ToLower(filterParam.Search) + "%"
		search = &abstraction.Search{
			Query: searchQuery,
			Args:  []interface{}{searchVal, searchVal},
		}
	}

	users, info, err := u.Repo.User.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, user := range users {
		role, err := u.Repo.Role.FindByID(ctx, user.RoleID)
		if err != nil {
			return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
		}

		userProfile, err := u.Repo.UserProfile.FindByUserID(ctx, user.ID)
		if err != nil {
			return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
		}

		result = append(result, dto.UserResponse{
			UserModel:   user,
			Role:        role,
			UserProfile: userProfile,
		})
	}

	return result, pagination, nil
}

func (u *user) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.UserResponse, error) {
	var result dto.UserResponse

	user, err := u.Repo.User.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	role, err := u.Repo.Role.FindByID(ctx, user.RoleID)
	if err != nil {
		return result, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	userProfile, err := u.Repo.UserProfile.FindByUserID(ctx, user.ID)
	if err != nil {
		return result, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	result = dto.UserResponse{
		UserModel:   *user,
		Role:        role,
		UserProfile: userProfile,
	}

	return result, nil
}

func (u *user) Create(ctx context.Context, payload dto.CreateUserRequest) (result dto.UserResponse, err error) {
	var email string
	if payload.Email != nil {
		email = *payload.Email
	} else {
		email = str.GenRandStr(10) + "@gmail.com"
	}

	var (
		userID = uuid.New().String()
		user   = model.UserModel{
			Entity: model.Entity{
				ID: userID,
			},
			UserEntity: model.UserEntity{
				Name:     payload.Name,
				Email:    email,
				Password: payload.Password,
			},
		}
		userProfile = &model.UserProfileModel{
			UserProfileEntity: model.UserProfileEntity{
				Phone:  payload.PhoneNumber,
				UserID: userID,
			},
		}
	)

	var role *model.RoleModel
	if payload.RoleName != nil {
		role, err = u.Repo.Role.FindByName(ctx, *payload.RoleName)
		if err != nil {
			return result, res.ErrorBuilder(res.Constant.Error.BadRequest, err, "role name invalid")
		}
	} else {
		role, err = u.Repo.Role.FindByID(ctx, *payload.RoleID)
		if err != nil {
			return result, res.ErrorBuilder(res.Constant.Error.BadRequest, err, "role id invalid")
		}
	}
	user.RoleID = role.ID

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.User.Create(ctx, user)
		if err != nil {
			return err
		}

		_, err = u.Repo.UserProfile.Create(ctx, *userProfile)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.UserResponse{
		UserModel:   user,
		Role:        role,
		UserProfile: userProfile,
	}

	return result, nil
}

func (u *user) Update(ctx context.Context, payload dto.UpdateUserRequest) (result dto.UserResponse, err error) {
	var (
		userData = &model.UserModel{
			UserEntity: model.UserEntity{
				Name:     payload.Name,
				Email:    payload.Email,
				Password: payload.Password,
			},
		}
		userProfile = &model.UserProfileModel{
			UserProfileEntity: model.UserProfileEntity{
				Phone: payload.PhoneNumber,
			},
		}
	)

	if payload.Password != "" {
		userData.HashPassword()
		userData.Password = ""
	}

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.User.UpdateByID(ctx, payload.ID, *userData)
		if err != nil {
			return err
		}

		_, err = u.Repo.UserProfile.UpdateByUserID(ctx, payload.ID, *userProfile)
		if err != nil {
			return err
		}

		userData, err = u.Repo.User.FindByID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}
		userProfile, err = u.Repo.UserProfile.FindByUserID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.UserResponse{
		UserModel:   *userData,
		UserProfile: userProfile,
	}

	return result, nil
}

func (u *user) Delete(ctx context.Context, payload dto.ByIDRequest) (result dto.UserResponse, err error) {
	var data *model.UserModel

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		data, err = u.Repo.User.FindByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.UserProfile.DeleteByUserID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.User.DeleteByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.UserResponse{
		UserModel: *data,
	}

	return result, nil
}
