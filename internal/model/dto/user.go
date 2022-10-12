package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateUserRequest struct {
		Name        string  `json:"name" validate:"required" example:"user"`
		Email       *string `json:"email,omitempty" validate:"omitempty" example:"user@gmail.com"`
		Password    string  `json:"password"`
		PhoneNumber string  `json:"phone_number"`
		RoleID      *string `json:"role_id,omitempty" validate:"required_without=RoleName,omitempty" example:"1"`
		RoleName    *string `json:"role_name,omitempty" validate:"required_without=RoleID,omitempty"`
	}

	UpdateUserRequest struct {
		ID          string `param:"id" validate:"required"`
		Name        string `json:"name"`
		Email       string `json:"email" validate:"omitempty,email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
	}
)

// response
type (
	UserResponse struct {
		model.UserModel
		Role        *model.RoleModel        `json:"role"`
		UserProfile *model.UserProfileModel `json:"user_profile"`
	}
	UserResponseDoc struct {
		Body struct {
			Meta res.Meta     `json:"meta"`
			Data UserResponse `json:"data"`
		} `json:"body"`
	}
)
