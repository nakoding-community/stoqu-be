package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	UpdateReminderStockHistoryRequest struct {
		ID     string `param:"id" validate:"required"`
		IsRead bool   `json:"is_read"`
	}

	UpdateReminderStockHistoryBulkReadRequest struct {
		IDs []string `json:"ids" validate:"required"`
	}
)

// response
type (
	ReminderStockHistoryResponse struct {
		model.ReminderStockHistoryModel
	}
	ReminderStockHistoryResponseDoc struct {
		Body struct {
			Meta res.Meta                     `json:"meta"`
			Data ReminderStockHistoryResponse `json:"data"`
		} `json:"body"`
	}

	ReminderStockHistoryBulkReadResponse struct {
		Status string `json:"status"`
	}
	ReminderStockHistoryBulkReadResponseDoc struct {
		Body struct {
			Meta res.Meta                             `json:"meta"`
			Data ReminderStockHistoryBulkReadResponse `json:"data"`
		} `json:"body"`
	}

	ReminderStockHistoryCountUnreadResponse struct {
		Count int64 `json:"count"`
	}
	ReminderStockHistoryCountUnreadResponseDoc struct {
		Body struct {
			Meta res.Meta                                `json:"meta"`
			Data ReminderStockHistoryCountUnreadResponse `json:"data"`
		} `json:"body"`
	}
)
