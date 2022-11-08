package entity

type OrderTrxEntity struct {
	TrxType        string  `json:"trx_type" gorm:"size:20;not null"`
	Code           string  `json:"code" gorm:"size:50;not null"`
	ShipmentType   string  `json:"shipment_type" gorm:"size:20;not null"`
	ShipmentNumber string  `json:"shipment_number" gorm:"size:150;not null"`
	ShipmentPrice  int     `json:"shipment_price" gorm:"not null;default:0"`
	Price          float64 `json:"price" gorm:"not null"`
	FinalPrice     float64 `json:"final_price" gorm:"not null"`
	Notes          string  `json:"notes"`
	PaymentStatus  string  `json:"payment_status" gorm:"size:10;not null"`
	StockStatus    string  `json:"stock_status" gorm:"size:20;not null"`
	Status         string  `json:"status" gorm:"size:20;not null"`
	IsRead         bool    `json:"is_read" gorm:"size:20;not null;default:0"`

	// fk
	PicID      string `json:"pic_id" gorm:"type:uuid;not null"`
	SupplierID string `json:"supplier_id" gorm:"type:uuid;not null"`
	CustomerID string `json:"customer_id" gorm:"type:uuid;not null"`
}

type OrderTrxModel struct {
	Entity
	OrderTrxEntity

	// relations
	OrderTrxItems    []OrderTrxItemModel    `json:"-" gorm:"foreignKey:OrderTrxID;constraint:OnDelete:CASCADE;"`
	OrderTrxStatuses []OrderTrxStatusModel  `json:"-" gorm:"foreignKey:OrderTrxID;constraint:OnDelete:CASCADE;"`
	OrderTrxReceipts []OrderTrxReceiptModel `json:"-" gorm:"foreignKey:OrderTrxID;constraint:OnDelete:CASCADE;"`
}

func (OrderTrxModel) TableName() string {
	return "order_trxs"
}

type OrderView struct {
	Entity
	OrderTrxEntity

	// join
	CustomerName string `json:"customer_name" filter:"column:customers.name"`
	SupplierName string `json:"supplier_name" filter:"column:suppliers.name"`
	PicName      string `json:"pic_name" filter:"column:pics.name"`
}
