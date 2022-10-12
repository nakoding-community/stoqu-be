package entity

type ProductEntity struct {
	Code     string  `json:"code" gorm:"size:50;not null"`
	PriceUSD float64 `json:"price_usd" gorm:"not null"`
	PriceIDR float64 `json:"price_idr" gorm:"not null"`

	// fk
	BrandID   string `json:"brand_id" gorm:"not null"`
	VariantID string `json:"variant_id" gorm:"not null"`
	PacketID  string `json:"packet_id" gorm:"not null"`
}

type ProductModel struct {
	Entity
	ProductEntity

	// relations
	Brand          *BrandModel          `json:"brand" gorm:"foreignKey:BrandID;"`
	Variant        *VariantModel        `json:"variant" gorm:"foreignKey:VariantID;"`
	Packet         *PacketModel         `json:"type" gorm:"foreignKey:PacketID;"`
	ProductLookups []ProductLookupModel `json:"product_lookups" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

func (ProductModel) TableName() string {
	return "products"
}
