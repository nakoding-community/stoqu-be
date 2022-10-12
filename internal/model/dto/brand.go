package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateBrandRequest struct {
		Name       string `json:"name" validate:"required"`
		SupplierID string `json:"supplier_id"`
	}

	UpdateBrandRequest struct {
		ID         string `param:"id" validate:"required"`
		Name       string `json:"name"`
		SupplierID string `json:"supplier_id"`
	}
)

// response
type (
	BrandResponse struct {
		model.BrandModel
	}
	BrandResponseDoc struct {
		Body struct {
			Meta res.Meta      `json:"meta"`
			Data BrandResponse `json:"data"`
		} `json:"body"`
	}
)
