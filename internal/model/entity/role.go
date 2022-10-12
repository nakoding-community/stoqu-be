package entity

type RoleEntity struct {
	Code string `json:"code" gorm:"not null;unique;size:50"`
	Name string `json:"name" gorm:"not null;unique;size:50"`
}

type RoleModel struct {
	Entity
	RoleEntity

	// relations
	Users []UserModel `json:"-" gorm:"foreignKey:RoleID"`
}

func (RoleModel) TableName() string {
	return "roles"
}
