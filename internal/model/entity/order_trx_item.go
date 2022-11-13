package entity

type OrderTrxItemEntity struct {
	Total       int     `json:"total" gorm:"not null"`
	TotalPacked int     `json:"total_packed" gorm:"not null"`
	Price       float64 `json:"price" gorm:"not null"`
	Status      string  `json:"status" gorm:"type:varchar(20);not null"`

	// fk
	OrderTrxID string `json:"order_trx_id" gorm:"not null"`
	ProductID  string `json:"product_id" gorm:"not null"`
	RackID     string `json:"rack_id" gorm:"not null"`
}

type OrderTrxItemModel struct {
	Entity
	OrderTrxItemEntity

	// relations
	OrderTrxItemLookups []OrderTrxItemLookupModel `json:"order_trx_item_lookups" gorm:"foreignKey:OrderTrxItemID;constraint:OnDelete:CASCADE;"`
	Product             *ProductModel             `json:"product" gorm:"foreignKey:ProductID;"`
	OrderTrx            *OrderTrxModel            `json:"order_trx" gorm:"foreignKey:OrderTrxID;"`
	Rack                *RackModel                `json:"rack" gorm:"foreignKey:RackID;"`

	// helper
	Action string `json:"-" gorm:"-"`
}

func (OrderTrxItemModel) TableName() string {
	return "order_trx_items"
}
