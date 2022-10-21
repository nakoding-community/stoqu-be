package entity

type RackEntity struct {
	Code string `json:"code" gorm:"not null;unique;size:50"`
	Name string `json:"name" gorm:"not null;unique;size:50"`
}

type RackModel struct {
	Entity
	RackEntity
}

func (RackModel) TableName() string {
	return "racks"
}
