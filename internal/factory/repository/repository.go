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
}

func Init(cfg *config.Configuration, db *gorm.DB) Factory {
	f := Factory{}

	f.Db = db

	f.Role = dbRepository.NewRole(f.Db)
	f.User = dbRepository.NewUser(f.Db)
	f.UserProfile = dbRepository.NewUserProfile(f.Db)

	return f
}
