package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
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
