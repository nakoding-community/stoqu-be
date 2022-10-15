package entity

type StockLookupEntity struct {
	Code                     string  `json:"code" gorm:"size:50;not null;unique"`
	IsSeal                   bool    `json:"is_seal" gorm:"not null"`
	TypeValue                float64 `json:"type_value" gorm:"not null"`
	RemainingTypeValue       float64 `json:"remaining_type_value" gorm:"not null"`
	RemainingTypeValueBefore float64 `json:"remaining_type_value_before" gorm:"not null"`

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
