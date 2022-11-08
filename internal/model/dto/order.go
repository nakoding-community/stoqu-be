package dto

import (
	"github.com/google/uuid"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
	"gitlab.com/stoqu/stoqu-be/pkg/util/str"
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
		Total     int     `json:"total" validate:"required,min=1"`
		Price     float64 `json:"price" validate:"required"`
		Status    string  `json:"status" validate:"required"`

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

// mapper
func (dto *UpsertOrderRequest) ToOrderTrx(mapStockLookups map[string]model.StockLookupModel) model.OrderTrxModel {
	orderTrxID := uuid.New().String()
	if dto.ID != "" {
		orderTrxID = dto.ID
	}
	orderTrx := model.OrderTrxModel{
		Entity: model.Entity{
			ID: orderTrxID,
		},
		OrderTrxEntity: model.OrderTrxEntity{
			TrxType:        dto.TrxType,
			Code:           str.GenCode(constant.CODE_ORDER_TRX_PREFIX),
			ShipmentType:   dto.ShipmentType,
			ShipmentNumber: dto.ShipmentNumber,
			ShipmentPrice:  dto.ShipmentPrice,
			Price:          dto.Price,
			FinalPrice:     dto.FinalPrice,
			Notes:          dto.Notes,
			PaymentStatus:  dto.PaymentStatus,
			StockStatus:    dto.StockStatus,
			Status:         dto.Status,
			IsRead:         dto.IsRead,
			PicID:          dto.PicID,
			SupplierID:     dto.SupplierID,
			CustomerID:     dto.CustomerID,
		},
	}

	// items
	for _, item := range dto.Items {
		orderTrxItemID := uuid.New().String()
		if item.ID != "" {
			orderTrxItemID = item.ID
		}
		orderTrxItem := model.OrderTrxItemModel{
			Entity: model.Entity{
				ID: orderTrxItemID,
			},
			OrderTrxItemEntity: model.OrderTrxItemEntity{
				Total:      item.Total,
				Price:      item.Price,
				Status:     item.Status,
				ProductID:  item.ProductID,
				OrderTrxID: orderTrxID,
			},
			Action: item.Action,
		}

		// lookups
		for _, lookup := range item.StockLookups {
			if stockLookup, ok := mapStockLookups[lookup.ID]; ok {
				orderTrxItemLookup := model.OrderTrxItemLookupModel{
					Entity: model.Entity{
						ID: lookup.ID,
					},
					OrderTrxItemLookupEntity: model.OrderTrxItemLookupEntity{
						StockLookupEntity: stockLookup.StockLookupEntity,
						OrderTrxItemID:    orderTrxItemID,
					},
					Action: lookup.Action,
				}
				orderTrxItem.OrderTrxItemLookups = append(orderTrxItem.OrderTrxItemLookups, orderTrxItemLookup)
			}
		}

		orderTrx.OrderTrxItems = append(orderTrx.OrderTrxItems, orderTrxItem)
	}

	// receipts
	for _, receipt := range orderTrx.OrderTrxReceipts {
		orderTrxReceiptID := uuid.New().String()
		if receipt.ID != "" {
			orderTrxReceiptID = receipt.ID
		}
		orderTrxReceipt := model.OrderTrxReceiptModel{
			Entity: model.Entity{
				ID: orderTrxReceiptID,
			},
			OrderTrxReceiptEntity: model.OrderTrxReceiptEntity{
				ReceiptUrl: receipt.ReceiptUrl,
				OrderTrxID: orderTrxID,
			},
			Action: receipt.Action,
		}
		orderTrx.OrderTrxReceipts = append(orderTrx.OrderTrxReceipts, orderTrxReceipt)
	}

	return orderTrx
}
