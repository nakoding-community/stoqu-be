package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateProductLookupRequest struct {
		IsSeal                   bool    `json:"is_seal"`
		TypeValue                float64 `json:"type_value"`
		RemainingTypeValue       float64 `json:"remaining_type_value"`
		RemainingTypeValueBefore float64 `json:"remaining_type_value_before"`
		ProductID                string  `json:"product_id" validate:"required"`
	}

	UpdateProductLookupRequest struct {
		ID                       string  `param:"id" validate:"required"`
		IsSeal                   bool    `json:"is_seal"`
		TypeValue                float64 `json:"type_value"`
		RemainingTypeValue       float64 `json:"remaining_type_value"`
		RemainingTypeValueBefore float64 `json:"remaining_type_value_before"`
	}
)

// response
type (
	ProductLookupResponse struct {
		model.ProductLookupModel
	}
	ProductLookupResponseDoc struct {
		Body struct {
			Meta res.Meta              `json:"meta"`
			Data ProductLookupResponse `json:"data"`
		} `json:"body"`
	}
)
