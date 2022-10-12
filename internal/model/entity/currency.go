package entity

type CurrencyEntity struct {
	Code   string  `json:"code" gorm:"not null;unique;size:50"`
	Name   string  `json:"name" gorm:"not null;unique;size:50"`
	IsAuto bool    `json:"is_auto"`
	Value  float64 `json:"value"`
}

type CurrencyModel struct {
	Entity
	CurrencyEntity
}

func (CurrencyModel) TableName() string {
	return "currencies"
}
