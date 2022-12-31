package entity

type StockRackEntity struct {
	TotalSeal    int64 `json:"total_seal" gorm:"not null"`
	TotalNotSeal int64 `json:"total_not_seal" gorm:"not null"`
	Total        int64 `json:"total" gorm:"not null"`

	// fk
	StockID string `json:"stock_id" gorm:"not null"`
	RackID  string `json:"rack_id" gorm:"not null"`
}

type StockRackModel struct {
	Entity
	StockRackEntity

	// relations
	Rack         *RackModel         `json:"rack" gorm:"foreignKey:RackID;"`
	StockLookups []StockLookupModel `json:"stock_lookups" gorm:"foreignKey:StockRackID;constraint:OnDelete:CASCADE;"`
}

func (StockRackModel) TableName() string {
	return "stock_racks"
}
