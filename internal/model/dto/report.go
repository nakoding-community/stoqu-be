package dto

import (
	"time"

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
		StartDate     string `json:"start_date"`
		EndDate       string `json:"end_date"`
		StartDateTime time.Time
		EndDateTime   time.Time
		Group         string `json:"group"`
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

	ProductReportResponse struct {
		model.OrderView
	}
	ProductReportResponseDoc struct {
		Body struct {
			Meta res.Meta          `json:"meta"`
			Data model.ProductView `json:"data"`
		} `json:"body"`
	}
)
