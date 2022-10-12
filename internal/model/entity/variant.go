package entity

type VariantEntity struct {
	Code       string `json:"code" gorm:"size:50;unique;not null"`
	Name       string `json:"name" gorm:"size:50;uniqueIndex:idx_unique_variant;not null"`
	ITL        string `json:"itl" gorm:"size:100;uniqueIndex:idx_unique_variant;not null"`
	UniqueCode string `json:"unique_code"`

	// fk
	BrandID string `json:"brand_id" gorm:"uniqueIndex:idx_unique_variant;not null"`
}

type VariantModel struct {
	Entity
	VariantEntity

	// relations
	Brand *BrandModel `json:"brand" gorm:"foreignKey:BrandID;constraint:OnDelete:CASCADE;"`
}

func (VariantModel) TableName() string {
	return "variants"
}
