package entity

import "github.com/google/uuid"

type OrderTrxItemLookupEntity struct {
	StockLookupEntity

	// fk
	OrderTrxItemID uuid.UUID `json:"order_trx_item_id" gorm:"type:uuid;not null"`
}

type OrderTrxItemLookupModel struct {
	Entity
	OrderTrxItemLookupEntity
}

func (OrderTrxItemLookupModel) TableName() string {
	return "order_trx_item_lookups"
}
