package entity

type OrderTrxItemLookupEntity struct {
	StockLookupEntity

	// fk
	OrderTrxItemID string `json:"order_trx_item_id" gorm:"not null"`
}

type OrderTrxItemLookupModel struct {
	Entity
	OrderTrxItemLookupEntity

	// helper
	Action string `json:"-" gorm:"-"`
}

func (OrderTrxItemLookupModel) TableName() string {
	return "order_trx_item_lookups"
}
