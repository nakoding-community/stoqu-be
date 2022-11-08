package entity

type OrderTrxStatusEntity struct {
	Status   string `json:"status" gorm:"size:20;not null"`
	IsActive bool   `json:"is_active" gorm:"uniqueIndex:idx_unique_order_trx_status"`

	// fk
	OrderTrxID string `json:"order_trx_id" gorm:"uniqueIndex:idx_unique_order_trx_status;not null"`
}

type OrderTrxStatusModel struct {
	Entity
	OrderTrxStatusEntity

	// helper
	Action string `json:"-" gorm:"-"`
}

func (OrderTrxStatusModel) TableName() string {
	return "order_trx_statuses"
}
