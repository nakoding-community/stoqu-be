package entity

type PacketEntity struct {
	Code  string `json:"code" gorm:"not null;unique;size:50"`
	Name  string `json:"name" gorm:"not null;unique;size:50"`
	Value int64  `json:"value"`

	// fk
	UnitID string `json:"unit_id" gorm:"not null"`
}

type PacketModel struct {
	Entity
	PacketEntity

	// relations
	Unit *UnitModel `json:"unit" gorm:"foreignKey:UnitID;"`

	// helper
	UnitName string `json:"unit_name" gorm:"-"`
}

func (PacketModel) TableName() string {
	return "packets"
}
