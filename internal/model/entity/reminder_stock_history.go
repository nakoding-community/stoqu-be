package entity

type ReminderStockHistoryEntity struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	IsRead bool   `json:"is_read" gorm:"not null;default:0"`
}

type ReminderStockHistoryModel struct {
	Entity
	ReminderStockHistoryEntity
}

func (ReminderStockHistoryModel) TableName() string {
	return "reminder_stock_histories"
}
