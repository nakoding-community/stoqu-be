package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateProductRequest struct {
		PriceUSD  float64 `json:"price_usd"`
		PriceIDR  float64 `json:"price_idr"`
		BrandID   string  `json:"brand_id" validate:"required"`
		VariantID string  `json:"variant_id" validate:"required"`
		PacketID  string  `json:"packet_id" validate:"required"`
	}

	UpdateProductRequest struct {
		ID       string  `param:"id" validate:"required"`
		PriceUSD float64 `json:"price_usd"`
		PriceIDR float64 `json:"price_idr"`
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
)
