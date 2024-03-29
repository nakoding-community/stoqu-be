package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateStockLookupRequest struct {
		IsSeal               bool    `json:"is_seal"`
		Value                float64 `json:"value"`
		RemainingValue       float64 `json:"remaining_value"`
		RemainingValueBefore float64 `json:"remaining_value_before"`
		StockRackID          string  `json:"stock_rack_id" validate:"required"`
	}

	UpdateStockLookupRequest struct {
		ID                   string  `param:"id" validate:"required"`
		IsSeal               bool    `json:"is_seal"`
		Value                float64 `json:"value"`
		RemainingValue       float64 `json:"remaining_value"`
		RemainingValueBefore float64 `json:"remaining_value_before"`
		StockRackID          string  `json:"stock_rack_id" validate:"required"`
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

	StockLookupViewResponse struct {
		model.StockLookupView
	}
	StockLookupViewResponseDoc struct {
		Body struct {
			Meta res.Meta                `json:"meta"`
			Data StockLookupViewResponse `json:"data"`
		} `json:"body"`
	}
)
