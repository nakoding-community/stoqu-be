package entity

type OrderTrxReceiptEntity struct {
	ReceiptUrl string `json:"receipt_url" gorm:"type:text;not null"`

	// fk
	OrderTrxID string `json:"order_trx_id" gorm:"not null"`
}

type OrderTrxReceiptModel struct {
	Entity
	OrderTrxReceiptEntity

	// helper
	Action string `json:"-" gorm:"-"`
}

func (OrderTrxReceiptModel) TableName() string {
	return "order_trx_receipts"
}
