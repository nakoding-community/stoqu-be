package entity

type StockEntity struct {
	TotalSeal    int64 `json:"price_usd" gorm:"not null"`
	TotalNotSeal int64 `json:"price_final" gorm:"not null"`
	Total        int64 `json:"toital" gorm:"not null"`

	// fk
	ProductID string `json:"product_id" gorm:"not null"`
	BrandID   string `json:"brand_id" gorm:"not null"`
	VariantID string `json:"variant_id" gorm:"not null"`
	PacketID  string `json:"packet_id" gorm:"not null"`
}

type StockModel struct {
	Entity
	StockEntity

	// relations
	Product *ProductModel `json:"product" gorm:"foreignKey:ProductID;"`
	Brand   *BrandModel   `json:"brand" gorm:"foreignKey:BrandID;"`
	Variant *VariantModel `json:"variant" gorm:"foreignKey:VariantID;"`
	Packet  *PacketModel  `json:"type" gorm:"foreignKey:PacketID;"`
}

func (StockModel) TableName() string {
	return "stocks"
}

type StockView struct {
	Entity
	StockEntity

	// join
	ProductCode  string `json:"product_code" filter:"column:products.code"`
	ProductName  string `json:"product_name" filter:"column:products.name"`
	BrandName    string `json:"brand_name" filter:"column:brands.name"`
	SupplierName string `json:"supplier_name" filter:"column:users.name"`
	VariantName  string `json:"variant_name" filter:"column:variants.name"`
	PacketValue  int64  `json:"packet_value" filter:"column:packets.value"`
	UnitName     string `json:"unit_name" filter:"column:units.name"`
}
