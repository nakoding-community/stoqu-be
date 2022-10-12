package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
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
