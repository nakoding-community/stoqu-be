package entity

type UnitEntity struct {
	Code string `json:"code" gorm:"not null;unique;size:50"`
	Name string `json:"name" gorm:"not null;unique;size:50"`
}

type UnitModel struct {
	Entity
	UnitEntity

	// relations
	Packets []PacketModel `json:"-" gorm:"foreignKey:UnitID;"`
}

func (UnitModel) TableName() string {
	return "units"
}
