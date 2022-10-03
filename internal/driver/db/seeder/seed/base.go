package seed

import "gorm.io/gorm"

type Seed interface {
	Run(conn *gorm.DB) error
	GetTag() string
}
