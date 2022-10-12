package dto

import (
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateConvertionUnitRequest struct {
		Origin      string  `json:"origin" validate:"required"`
		Destination string  `json:"destination" validate:"required"`
		Total       float64 `json:"total" validate:"required"`
	}
)

type (
	UpdateConvertionUnitRequest struct {
		ID          string  `param:"id" validate:"required"`
		Origin      string  `json:"origin"`
		Destination string  `json:"destination"`
		Total       float64 `json:"total"`
	}
)

// response
type (
	ConvertionUnitResponse struct {
		ID          string  `json:"id"`
		Origin      string  `json:"origin"`
		Destination string  `json:"Destination"`
		Total       float64 `json:"total"`
	}
	ConvertionUnitResponseDoc struct {
		Body struct {
			Meta res.Meta               `json:"meta"`
			Data ConvertionUnitResponse `json:"data"`
		} `json:"body"`
	}
)
