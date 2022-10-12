package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateRoleRequest struct {
		Name string `json:"name" validate:"required"`
	}

	UpdateRoleRequest struct {
		ID   string `param:"id" validate:"required"`
		Name string `json:"name"`
	}
)

// response
type (
	RoleResponse struct {
		model.RoleModel
	}
	RoleResponseDoc struct {
		Body struct {
			Meta res.Meta     `json:"meta"`
			Data RoleResponse `json:"data"`
		} `json:"body"`
	}
)
