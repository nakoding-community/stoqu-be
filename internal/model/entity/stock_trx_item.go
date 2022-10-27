package entity

type StockTrxItemEntity struct {
	TotalSeal    int    `json:"total_seal" gorm:"not null"`
	TotalNotSeal int    `json:"total_not_seal" gorm:"not null"`
	ConvertType  string `json:"convert_type"`

	// fk
	StockTrxID string `json:"stock_trx_id"`
	StockID    string `json:"stock_id"`
	ProductID  string `json:"product_id"`
}

type StockTrxItemModel struct {
	Entity
	StockTrxItemEntity

	// relations
	StockTrxItemLookups []StockTrxItemLookupModel `json:"stock_trx_item_lookups" gorm:"foreignKey:StockTrxItemID;constraint:OnDelete:CASCADE;"`
}

func (StockTrxItemModel) TableName() string {
	return "stock_trx_items"
}
