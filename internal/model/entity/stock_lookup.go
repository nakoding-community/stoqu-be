package entity

type StockLookupEntity struct {
	Code                 string  `json:"code" gorm:"size:50;not null;unique"`
	IsSeal               bool    `json:"is_seal" gorm:"not null"`
	Value                float64 `json:"value" gorm:"not null"`
	RemainingValue       float64 `json:"remaining_value" gorm:"not null"`
	RemainingValueBefore float64 `json:"remaining_value_before" gorm:"not null"`

	// fk
	StockID string `json:"stock_id" gorm:"not null"`
}

type StockLookupModel struct {
	Entity
	StockLookupEntity
}

func (StockLookupEntity) TableName() string {
	return "stock_lookups"
}
