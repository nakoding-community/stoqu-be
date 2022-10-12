package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	LoginAuthRequest struct {
		Email    string `json:"email" validate:"required,email" example:"admin@gmail.com"`
		Password string `json:"password" validate:"required" example:"admin"`
	}
	RegisterAuthRequest struct {
		model.UserEntity
	}
)

// response
type (
	AuthLoginResponse struct {
		Token string `json:"token"`
		Role  string `json:"role"`
		model.UserModel
	}
	AuthLoginResponseDoc struct {
		Body struct {
			Meta res.Meta          `json:"meta"`
			Data AuthLoginResponse `json:"data"`
		} `json:"body"`
	}

	AuthRegisterResponse struct {
		model.UserModel
	}
	AuthRegisterResponseDoc struct {
		Body struct {
			Meta res.Meta             `json:"meta"`
			Data AuthRegisterResponse `json:"data"`
		} `json:"body"`
	}
)
