package entity

type ConvertionUnitEntity struct {
	Code            string  `json:"code" gorm:"not null;unique;size:50"`
	Name            string  `json:"name" gorm:"not null;unique;size:50"`
	ValueConvertion float64 `json:"value_convertion"`

	// fk
	UnitOriginID      string `json:"unit_origin_id" gorm:"uniqueIndex:idx_unique_unit"`
	UnitDestinationID string `json:"unit_destination_id" gorm:"uniqueIndex:idx_unique_unit"`
}

type ConvertionUnitModel struct {
	Entity
	ConvertionUnitEntity
}

func (ConvertionUnitModel) TableName() string {
	return "convertion_units"
}
