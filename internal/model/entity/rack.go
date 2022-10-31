package entity

type RackEntity struct {
	Code string `json:"code" gorm:"not null;unique;size:50"`
	Name string `json:"name" gorm:"not null;unique;size:50"`
}

type RackModel struct {
	Entity
	RackEntity

	StockRacks []StockRackModel `json:"stock_racks" gorm:"foreignKey:RackID;constraint:OnDelete:CASCADE;"`
}

func (RackModel) TableName() string {
	return "racks"
}
