package dto

import (
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// response
type (
	DashboardResponse struct {
		TotalBrand   int64                         `json:"total_brand"`
		TotalProduct int64                         `json:"total_product"`
		TotalStock   int64                         `json:"total_stock"`
		TotalOrder   int64                         `json:"total_order"`
		OrderDaily   []DashboardOrderDailyResponse `json:"order_daily"`
	}

	DashboardOrderDailyResponse struct {
		Day   string `json:"day"`
		Total int    `json:"total"`
	}

	DashboardResponseDoc struct {
		Body struct {
			Meta res.Meta          `json:"meta"`
			Data DashboardResponse `json:"data"`
		} `json:"body"`
	}
)
