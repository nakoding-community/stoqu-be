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
		StartDate string `json:"start_date" filter:"column:order_trxs.created_at;operator:>="`
		EndDate   string `json:"end_date" filter:"column:order_trxs.created_at;operator:>="`
		Group     string `query:"group"`
	}

	GenerateExcelReportInput struct {
		SheetName   string                   `query:"sheet_name"`
		FileName    string                   `query:"file_name"`
		Data        []map[string]interface{} `json:"data"`
		DataMapping DataMapping              `json:"data_mapping"`
	}

	DataMapping struct {
		Headers map[string]string `json:"headers"`
		Body    map[string]string `json:"body"`
	}
)

// response
type (
	OrderReportResponse struct {
		TotalOrder  int64             `json:"total_order"`
		TotalIncome int64             `json:"total_income"`
		Orders      []model.OrderView `json:"orders"`
	}
	OrderReportResponseDoc struct {
		Body struct {
			Meta res.Meta            `json:"meta"`
			Data OrderReportResponse `json:"data"`
		} `json:"body"`
	}

	OrderProductReportResponse struct {
		Total  int64                    `json:"total"`
		Orders []model.OrderViewProduct `json:"orders"`
	}
	OrderProductReportResponseDoc struct {
		Body struct {
			Meta res.Meta                   `json:"meta"`
			Data OrderProductReportResponse `json:"data"`
		} `json:"body"`
	}
)
