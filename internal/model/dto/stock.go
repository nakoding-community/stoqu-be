package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateStockRequest struct {
		Name      string `json:"name"`
		ProductID string `json:"product_id" validate:"required"`
	}
)

// response
type (
	StockResponse struct {
		model.StockModel
	}
	StockResponseDoc struct {
		Body struct {
			Meta res.Meta      `json:"meta"`
			Data StockResponse `json:"data"`
		} `json:"body"`
	}

	StockViewResponse struct {
		model.StockView
	}
	StockViewResponseDoc struct {
		Body struct {
			Meta res.Meta          `json:"meta"`
			Data StockViewResponse `json:"data"`
		} `json:"body"`
	}
)
