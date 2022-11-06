package entity

type StockLookupEntity struct {
	Code                 string  `json:"code" gorm:"size:50;not null;unique"`
	IsSeal               bool    `json:"is_seal" gorm:"not null"`
	Value                float64 `json:"value" gorm:"not null"`
	RemainingValue       float64 `json:"remaining_value" gorm:"not null"`
	RemainingValueBefore float64 `json:"remaining_value_before" gorm:"not null"`

	// fk
	StockRackID string `json:"stock_rack_id" gorm:"not null"`
}

type StockLookupModel struct {
	Entity
	StockLookupEntity
}

func (StockLookupEntity) TableName() string {
	return "stock_lookups"
}

type StockLookupView struct {
	Entity
	StockLookupEntity

	// join
	ProductID string `json:"product_id" filter:"column:stocks.product_id"`
	RackID    string `json:"rack_id" filter:"column:stock_racks.rack_id"`
}
