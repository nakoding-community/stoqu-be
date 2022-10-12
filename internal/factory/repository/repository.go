package repository

import (
	"cloud.google.com/go/firestore"
	"firebase.google.com/go/messaging"
	el "github.com/olivere/elastic/v7"
	"gitlab.com/stoqu/stoqu-be/internal/config"
	dbRepository "gitlab.com/stoqu/stoqu-be/internal/repository/db"
	"gorm.io/gorm"
)

type Factory struct {
	Db        *gorm.DB
	Es        *el.Client
	Fcm       *messaging.Client
	Firestore *firestore.Client

	Role        dbRepository.Role
	User        dbRepository.User
	UserProfile dbRepository.UserProfile

	Unit           dbRepository.Unit
	Packet         dbRepository.Packet
	ConvertionUnit dbRepository.ConvertionUnit
	Currency       dbRepository.Currency
	ReminderStock  dbRepository.ReminderStock

	Brand   dbRepository.Brand
	Variant dbRepository.Variant
}

func Init(cfg *config.Configuration, db *gorm.DB) Factory {
	f := Factory{}

	f.Db = db
	f.Role = dbRepository.NewRole(f.Db)
	f.User = dbRepository.NewUser(f.Db)
	f.UserProfile = dbRepository.NewUserProfile(f.Db)

	f.Unit = dbRepository.NewUnit(f.Db)
	f.Packet = dbRepository.NewPacket(f.Db)
	f.ConvertionUnit = dbRepository.NewConvertionUnit(f.Db)
	f.Currency = dbRepository.NewCurrency(f.Db)
	f.ReminderStock = dbRepository.NewReminderStock(f.Db)

	f.Brand = dbRepository.NewBrand(f.Db)
	f.Variant = dbRepository.NewVariant(f.Db)

	return f
}
