package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateProductRequest struct {
		Name       string  `json:"name"`
		PriceUSD   float64 `json:"price_usd"`
		PriceFinal float64 `json:"price_final"`
		BrandID    string  `json:"brand_id" validate:"required"`
		VariantID  string  `json:"variant_id" validate:"required"`
		PacketID   string  `json:"packet_id" validate:"required"`
	}

	UpdateProductRequest struct {
		ID         string  `param:"id" validate:"required"`
		Name       string  `json:"name"`
		PriceUSD   float64 `json:"price_usd"`
		PriceFinal float64 `json:"price_final"`
	}
)

// response
type (
	ProductResponse struct {
		model.ProductModel
	}
	ProductResponseDoc struct {
		Body struct {
			Meta res.Meta        `json:"meta"`
			Data ProductResponse `json:"data"`
		} `json:"body"`
	}

	ProductViewResponse struct {
		model.ProductView
	}
	ProductViewResponseDoc struct {
		Body struct {
			Meta res.Meta            `json:"meta"`
			Data ProductViewResponse `json:"data"`
		} `json:"body"`
	}
)
