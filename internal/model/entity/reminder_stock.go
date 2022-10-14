package entity

type ReminderStockEntity struct {
	Code     string `json:"code" gorm:"size:50;unique;not null"`
	Name     string `json:"name" gorm:"size:255;not null"`
	MinStock int64  `json:"min_stock" gorm:"default:0"`
}

type ReminderStockModel struct {
	Entity
	ReminderStockEntity
}

func (ReminderStockModel) TableName() string {
	return "reminder_stocks"
}
