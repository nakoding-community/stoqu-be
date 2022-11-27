package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	OrderReportQuery struct {
		StartDate string `json:"start_date" filter:"column:order_trxs.created_at;operator:>="`
		EndDate   string `json:"end_date" filter:"column:order_trxs.created_at;operator:<="`
		Status    string `json:"status" filter:"column:order_trxs.status"`
	}

	ProductReportQuery struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
)

// response
type (
	OrderReportResponse struct {
		model.OrderView
	}
	OrderReportResponseDoc struct {
		Body struct {
			Meta res.Meta        `json:"meta"`
			Data model.OrderView `json:"data"`
		} `json:"body"`
	}

	OrderProductReportResponse struct {
		model.OrderViewProduct
	}
	OrderProductReportResponseDoc struct {
		Body struct {
			Meta res.Meta               `json:"meta"`
			Data model.OrderViewProduct `json:"data"`
		} `json:"body"`
	}
)
