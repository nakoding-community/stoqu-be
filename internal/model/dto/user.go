package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateUserRequest struct {
		Name        string  `json:"name" validate:"required"`
		Email       *string `json:"email,omitempty" validate:"omitempty"`
		Password    string  `json:"password"`
		PhoneNumber string  `json:"phone_number"`
		RoleID      *string `json:"role_id,omitempty" validate:"required_without=RoleName,omitempty"`
		RoleName    *string `json:"role_name,omitempty" validate:"required_without=RoleID,omitempty"`
	}
)

type (
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
