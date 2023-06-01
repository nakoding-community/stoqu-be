package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.com/stoqu/stoqu-be/internal/config"
	"gitlab.com/stoqu/stoqu-be/pkg/util/gracefull"
	"gorm.io/gorm"
)

type (
	Db interface {
		Init() (*gorm.DB, gracefull.ProcessStopper, error)
	}

	db struct {
		Host string
		User string
		Pass string
		Port string
		Name string
	}
)

var (
	dbConnections map[string]*gorm.DB
)

func Init(cfg *config.Configuration) (gracefull.ProcessStopper, error) {
	dbConnections = make(map[string]*gorm.DB)
	var (
		stoppers         []gracefull.ProcessStopper
		dbConfigurations = map[string]Db{}
	)

	for _, v := range cfg.Databases {
		dbConfigurations[strings.ToUpper(v.DBName)] = &dbPostgreSQL{
			db: db{
				Host: v.DBHost,
				Name: v.DBName,
				Port: v.DBPort,
				Pass: v.DBPass,
				User: v.DBUser,
			},
			SslMode:  v.DBSSL,
			Tz:       v.DBTZ,
			LogLevel: v.DBLogLevel,
		}
	}

	for k, v := range dbConfigurations {
		db, stopper, err := v.Init()
		if err != nil {
			logrus.Info(err)
			return nil, fmt.Errorf("failed to connect to database %s", k)
		}
		dbConnections[k] = db
		stoppers = append(stoppers, stopper)
		logrus.Info(fmt.Sprintf("successfully connected to database %s", k))
	}

	stopper := func(ctx context.Context) (err error) {
		for _, stopper := range stoppers {
			err = stopper(ctx)
			if err != nil {
				return
			}
		}
		return
	}

	return stopper, nil
}

func GetConnection(name string) (*gorm.DB, error) {
	name = strings.ToUpper(name)
	if dbConnections[name] == nil {
		return nil, errors.New("connection is undefined")
	}
	return dbConnections[name], nil
}
