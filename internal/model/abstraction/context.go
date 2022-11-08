package abstraction

import (
	"gorm.io/gorm"
)

type AuthContext struct {
	ID    string
	Name  string
	Email string
}

type TrxContext struct {
	Db *gorm.DB
}
