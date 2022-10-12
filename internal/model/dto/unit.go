package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateUnitRequest struct {
		Name string `json:"name" validate:"required"`
	}

	UpdateUnitRequest struct {
		ID   string `param:"id" validate:"required"`
		Name string `json:"name"`
	}
)

// response
type (
	UnitResponse struct {
		model.UnitModel
	}
	UnitResponseDoc struct {
		Body struct {
			Meta res.Meta     `json:"meta"`
			Data UnitResponse `json:"data"`
		} `json:"body"`
	}
)
