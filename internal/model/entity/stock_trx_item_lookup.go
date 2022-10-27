package entity

type StockTrxItemLookupEntity struct {
	StockLookupEntity
	Code string `json:"code" gorm:"size:50;not null"`

	// fk
	StockTrxItemID string `json:"stock_trx_item_id"`
}

type StockTrxItemLookupModel struct {
	Entity
	StockTrxItemLookupEntity
}

func (StockTrxItemLookupModel) TableName() string {
	return "stock_trx_item_lookups"
}
