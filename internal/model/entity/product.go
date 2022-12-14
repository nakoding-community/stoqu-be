package entity

type ProductEntity struct {
	Code       string  `json:"code" gorm:"size:50;not null"`
	Name       string  `json:"name" gorm:"size:255;not null"`
	PriceUSD   float64 `json:"price_usd" gorm:"not null"`
	PriceFinal float64 `json:"price_final" gorm:"not null"`

	// fk
	BrandID   string `json:"brand_id" gorm:"not null"`
	VariantID string `json:"variant_id" gorm:"not null"`
	PacketID  string `json:"packet_id" gorm:"not null"`
}

type ProductModel struct {
	Entity
	ProductEntity

	// relations
	Brand   *BrandModel   `json:"brand" gorm:"foreignKey:BrandID;"`
	Variant *VariantModel `json:"variant" gorm:"foreignKey:VariantID;"`
	Packet  *PacketModel  `json:"packet" gorm:"foreignKey:PacketID;"`
	Stock   *StockModel   `json:"stock" gorm:"foreignKey:ProductID;"`
}

func (ProductModel) TableName() string {
	return "products"
}

type ProductView struct {
	Entity
	ProductEntity

	// join
	BrandName         string `json:"brand_name" filter:"column:brands.name"`
	SupplierName      string `json:"supplier_name" filter:"column:users.name"`
	VariantName       string `json:"variant_name" filter:"column:variants.name"`
	VariantUniqueCode string `json:"variant_unique_code" filter:"column:variants.unique_code"`
	PacketValue       int64  `json:"packet_value" filter:"column:packets.value"`
	UnitName          string `json:"unit_name" filter:"column:units.name"`
}
