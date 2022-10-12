package usecase

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/stoqu/stoqu-be/internal/config"
	resConstant "gitlab.com/stoqu/stoqu-be/pkg/util/response/constant"

	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
	"gitlab.com/stoqu/stoqu-be/pkg/util/trxmanager"

	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	Login(ctx context.Context, payload dto.LoginAuthRequest) (dto.AuthLoginResponse, error)
	Register(ctx context.Context, payload dto.RegisterAuthRequest) (dto.AuthRegisterResponse, error)
}

type auth struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewAuth(cfg *config.Configuration, f repository.Factory) Auth {
	return &auth{f, cfg}
}

func (u *auth) Login(ctx context.Context, payload dto.LoginAuthRequest) (dto.AuthLoginResponse, error) {
	var result dto.AuthLoginResponse

	data, err := u.Repo.User.FindByEmail(ctx, payload.Email)
	if data == nil {
		return result, res.CustomErrorBuilder(http.StatusUnauthorized, resConstant.E_INVALID_CREDENTIALS, err, "")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(data.PasswordHash), []byte(payload.Password)); err != nil {
		return result, res.CustomErrorBuilder(http.StatusUnauthorized, resConstant.E_INVALID_CREDENTIALS, err, "")
	}

	token, err := data.GenerateToken()
	if err != nil {
		return result, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}

	result = dto.AuthLoginResponse{
		Token:     token,
		UserModel: *data,
	}

	return result, nil
}

func (u *auth) Register(ctx context.Context, payload dto.RegisterAuthRequest) (dto.AuthRegisterResponse, error) {
	var result dto.AuthRegisterResponse
	var user model.UserModel
	var err error

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {

		userID := uuid.New().String()
		user = model.UserModel{
			Entity: model.Entity{
				ID: userID,
			},
			UserEntity: payload.UserEntity,
		}
		userProfile := model.UserProfileModel{
			UserProfileEntity: model.UserProfileEntity{
				Address: "xxx",
				Phone:   "021",
				UserID:  userID,
			},
		}

		role, err := u.Repo.Role.FindByName(ctx, "customer")
		if err != nil {
			return err
		}
		user.RoleID = role.ID

		_, err = u.Repo.User.Create(ctx, user)
		if err != nil {
			return err
		}

		_, err = u.Repo.UserProfile.Create(ctx, userProfile)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.AuthRegisterResponse{
		UserModel: user,
	}

	return result, nil
}
