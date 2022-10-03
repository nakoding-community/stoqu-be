package entity

type UserProfileEntity struct {
	Address string `json:"address"`
	Phone   string `json:"phone"`
	UserID  string `json:"user_id" gorm:"not null"`
}

type UserProfileModel struct {
	Entity
	UserProfileEntity
}

func (UserProfileModel) TableName() string {
	return "user_profiles"
}
