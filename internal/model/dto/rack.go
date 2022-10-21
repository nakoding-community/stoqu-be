package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateRackRequest struct {
		Name string `json:"name" validate:"required"`
	}

	UpdateRackRequest struct {
		ID   string `param:"id" validate:"required"`
		Name string `json:"name"`
	}
)

// response
type (
	RackResponse struct {
		model.RackModel
	}
	RackResponseDoc struct {
		Body struct {
			Meta res.Meta     `json:"meta"`
			Data RackResponse `json:"data"`
		} `json:"body"`
	}
)
