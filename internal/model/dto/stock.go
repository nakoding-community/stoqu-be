package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	TransactionStockRequest struct {
		TrxType  string                           `json:"trx_type" validate:"oneof=in out"`
		Products []TransactionStockProductRequest `json:"products"`
	}

	TransactionStockProductRequest struct {
		ID             string   `json:"id" validate:"required"`
		Quantity       int      `json:"quantity" validate:"required,min=1"`
		RackID         string   `json:"rack_id" validate:"required"`
		StockLookupIDs []string `json:"stock_lookup_ids"`
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

	StockTransactionResponse struct {
		Status string `json:"status"`
	}
	StockTransactionResponseDoc struct {
		Body struct {
			Meta res.Meta                 `json:"meta"`
			Data StockTransactionResponse `json:"data"`
		} `json:"body"`
	}
)
