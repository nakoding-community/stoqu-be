package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateCurrencyRequest struct {
		Name   string  `json:"name" validate:"required"`
		IsAuto bool    `json:"is_auto"`
		Value  float64 `json:"value"`
	}

	UpdateCurrencyRequest struct {
		ID     string  `param:"id" validate:"required"`
		Name   string  `json:"name"`
		IsAuto bool    `json:"is_auto"`
		Value  float64 `json:"value"`
	}

	ConvertCurrencyRequest struct {
		USD int `json:"usd"`
	}
)

// response
type (
	CurrencyResponse struct {
		model.CurrencyModel
	}
	CurrencyResponseDoc struct {
		Body struct {
			Meta res.Meta         `json:"meta"`
			Data CurrencyResponse `json:"data"`
		} `json:"body"`
	}

	ConvertCurrencyResponse struct {
		IDR int `json:"idr"`
	}
	ConvertCurrencyResponseDoc struct {
		Body struct {
			Meta res.Meta                `json:"meta"`
			Data ConvertCurrencyResponse `json:"data"`
		} `json:"body"`
	}

	CurrencyNowAPI struct {
		Date string  `json:"date"`
		Idr  float64 `json:"idr"`
	}
)
