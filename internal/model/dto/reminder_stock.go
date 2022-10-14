package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreateReminderStockRequest struct {
		Name     string `json:"name"`
		MinStock int64  `json:"min_stock"`
	}

	UpdateReminderStockRequest struct {
		ID       string `param:"id" validate:"required"`
		Name     string `json:"name"`
		MinStock int64  `json:"min_stock"`
	}
)

// response
type (
	ReminderStockResponse struct {
		model.ReminderStockModel
	}
	ReminderStockResponseDoc struct {
		Body struct {
			Meta res.Meta              `json:"meta"`
			Data ReminderStockResponse `json:"data"`
		} `json:"body"`
	}
)
