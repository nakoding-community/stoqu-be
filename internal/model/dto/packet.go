package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	CreatePacketRequest struct {
		UnitID string `json:"unit_id" validate:"required"`
		Value  int    `json:"value" validate:"required"`
	}

	UpdatePacketRequest struct {
		ID     string `param:"id" validate:"required"`
		UnitID string `json:"unit_id" validate:"required"`
		Value  int    `json:"value"`
	}
)

// response
type (
	PacketResponse struct {
		model.PacketModel
	}
	PacketResponseDoc struct {
		Body struct {
			Meta res.Meta       `json:"meta"`
			Data PacketResponse `json:"data"`
		} `json:"body"`
	}
)
