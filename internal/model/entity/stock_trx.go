package entity

type StockTrxEntity struct {
	Code       string `json:"code" gorm:"size:50;not null"`
	TrxType    string `json:"trx_type" gorm:"size:20;not null"`
	OrderTrxID string `json:"order_trx_id"`
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

type StockTrxView struct {
	Entity
	StockTrxEntity

	// join
	OrderCode string `json:"order_code" filter:"column:order_trxs.code"`

	// relations
	StockTrxItems []StockTrxItemModel `json:"stock_lookups" gorm:"foreignKey:StockTrxID;constraint:OnDelete:CASCADE;"`
}
