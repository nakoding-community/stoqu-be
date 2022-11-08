package dto

import (
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
)

// request
type (
	UpsertOrderRequest struct {
		ID             string `json:"id"`
		TrxType        string `json:"trx_type" validate:"required,oneof=in out"`
		ShipmentType   string `json:"shipment_type"`
		ShipmentNumber string `json:"shipment_number"`

		ShipmentPrice int     `json:"shipment_price"`
		Price         float64 `json:"price" validate:"required"`
		FinalPrice    float64 `json:"final_price" validate:"required"`

		Notes string `json:"notes"`

		PaymentStatus string `json:"payment_status" validate:"required,oneof=paid unpaid dp"`
		StockStatus   string `json:"stock_status" validate:"required"`
		Status        string `json:"status" validate:"required,oneof=pending progress completed canceled"`

		IsRead bool `json:"is_read"`

		PicID      string `json:"pic_id" validate:"required"`
		SupplierID string `json:"supplier_id"`
		CustomerID string `json:"customer_id"`

		Items    []UpsertOrderItemRequest    `json:"items" validate:"required,dive,min=1"`
		Receipts []UpsertOrderReceiptRequest `json:"receipts" validate:"dive"`
	}

	UpsertOrderItemRequest struct {
		ID        string  `json:"id"`
		ProductID string  `json:"product_id" validate:"required"`
		Quantity  int     `json:"quantity" validate:"required,min=1"`
		Price     float64 `json:"price" validate:"required"`

		StockRackID  string                         `json:"stock_rack_id"`
		StockLookups []UpsertOrderItemLookupRequest `json:"stock_lookups"`

		Action string `json:"action" validate:"required,oneof=insert update delete"`
	}

	UpsertOrderItemLookupRequest struct {
		ID     string `json:"id" validate:"required"`
		Action string `json:"action" validate:"required,oneof=insert update delete"`
	}

	UpsertOrderReceiptRequest struct {
		ID         string `json:"id"`
		ReceiptUrl string `json:"receipt_url" validate:"required"`
		Action     string `json:"action" validate:"required,oneof=insert update delete"`
	}
)

// response
type (
	OrderViewResponse struct {
		model.OrderView
	}
	OrderViewResponseDoc struct {
		Body struct {
			Meta res.Meta        `json:"meta"`
			Data model.OrderView `json:"data"`
		} `json:"body"`
	}

	OrderUpsertResponse struct {
		Status string `json:"status"`
	}
	OrderUpsertResponseDoc struct {
		Body struct {
			Meta res.Meta            `json:"meta"`
			Data OrderUpsertResponse `json:"data"`
		} `json:"body"`
	}
)
