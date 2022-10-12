package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateVariantRequest struct {
		Name       string `json:"name" validate:"required"`
		ITL        string `json:"itl"`
		UniqueCode string `json:"unique_code,omitempty"`
		BrandID    string `json:"brand_id" validate:"required"`
	}

	UpdateVariantRequest struct {
		ID         string `param:"id" validate:"required"`
		Name       string `json:"name"`
		ITL        string `json:"itl"`
		UniqueCode string `json:"unique_code,omitempty"`
	}
)

// response
type (
	VariantResponse struct {
		model.VariantModel
	}
	VariantResponseDoc struct {
		Body struct {
			Meta res.Meta        `json:"meta"`
			Data VariantResponse `json:"data"`
		} `json:"body"`
	}
)
