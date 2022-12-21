package entity

type BrandEntity struct {
	Code string `json:"code" gorm:"size:50;unique;not null"`
	Name string `json:"name" gorm:"size:255;not null"`

	// fk
	SupplierID string `json:"supplier_id" gorm:"not null"`
}

type BrandModel struct {
	Entity
	BrandEntity

	// relations
	Variants []VariantModel `json:"variants" table:"variants" gorm:"foreignKey:BrandID"`
	Supplier *UserModel     `json:"supplier" table:"users" gorm:"foreignKey:SupplierID"`

	// helper
	SupplierName string `json:"supplier_name" gorm:"-"`
}

func (BrandModel) TableName() string {
	return "brands"
}
