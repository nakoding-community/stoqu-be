package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	TransactionStockRequest struct {
		TrxType  string                           `json:"trx_type" validate:"required,oneof=in out"`
		Products []TransactionStockProductRequest `json:"products" validate:"required,min=1,dive"`
	}

	TransactionStockProductRequest struct {
		ID             string   `json:"id" validate:"required"`
		Quantity       int      `json:"quantity" validate:"required,min=1"`
		RackID         string   `json:"rack_id" validate:"required"`
		StockLookupIDs []string `json:"stock_lookup_ids"`
	}

	ConvertionStockRequest struct {
		Origin      ConvertionStocOriginkRequest     `json:"origin" validate:"required,dive"`
		Destination ConvertionStocDestinationRequest `json:"destination" validate:"required,dive"`
	}

	ConvertionStocOriginkRequest struct {
		ProductID      string   `json:"product_id" validate:"required"`
		RackID         string   `json:"rack_id" validate:"required"`
		StockLookupIDs []string `json:"stock_lookup_ids"`
	}

	ConvertionStocDestinationRequest struct {
		PacketID string `json:"packet_id" validate:"required"`
		Total    int64  `json:"total" validate:"required,min=1"`
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

	StockConvertionResponse struct {
		Status string `json:"status"`
	}
	StockConvertionResponseDoc struct {
		Body struct {
			Meta res.Meta                `json:"meta"`
			Data StockConvertionResponse `json:"data"`
		} `json:"body"`
	}
)
