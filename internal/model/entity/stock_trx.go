package entity

type StockTrxEntity struct {
	TrxType string `json:"trx_type" gorm:"size:20;not null"`
	Code    string `json:"code" gorm:"size:50;not null"`
}

type StockTrxModel struct {
	Entity
	StockTrxEntity

	// relations
	StockTrxItems []StockTrxItemModel `json:"stock_lookups" gorm:"foreignKey:StockTrxID;constraint:OnDelete:CASCADE;"`
}

func (StockTrxModel) TableName() string {
	return "stock_trxs"
}
