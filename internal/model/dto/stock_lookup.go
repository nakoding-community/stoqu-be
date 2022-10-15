package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateStockLookupRequest struct {
		IsSeal                   bool    `json:"is_seal"`
		TypeValue                float64 `json:"type_value"`
		RemainingTypeValue       float64 `json:"remaining_type_value"`
		RemainingTypeValueBefore float64 `json:"remaining_type_value_before"`
		StockID                  string  `json:"stock_id" validate:"required"`
	}

	UpdateStockLookupRequest struct {
		ID                       string  `param:"id" validate:"required"`
		IsSeal                   bool    `json:"is_seal"`
		TypeValue                float64 `json:"type_value"`
		RemainingTypeValue       float64 `json:"remaining_type_value"`
		RemainingTypeValueBefore float64 `json:"remaining_type_value_before"`
	}
)

// response
type (
	StockLookupResponse struct {
		model.StockLookupModel
	}
	StockLookupResponseDoc struct {
		Body struct {
			Meta res.Meta            `json:"meta"`
			Data StockLookupResponse `json:"data"`
		} `json:"body"`
	}
)
