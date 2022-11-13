package factory

import (
	"context"

	"gitlab.com/stoqu/stoqu-be/internal/config"
	"gitlab.com/stoqu/stoqu-be/internal/driver/db"
	"gitlab.com/stoqu/stoqu-be/internal/driver/db/automigration"
	"gitlab.com/stoqu/stoqu-be/internal/driver/db/seeder"
	"gitlab.com/stoqu/stoqu-be/internal/driver/firebase"
	"gitlab.com/stoqu/stoqu-be/internal/driver/ws"
	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/factory/usecase"
	"gitlab.com/stoqu/stoqu-be/pkg/util/gracefull"
)

type Factory struct {
	Repository repository.Factory
	Usecase    usecase.Factory
	WsHub      *ws.Hub
}

func Init(cfg *config.Configuration) (Factory, gracefull.ProcessStopper, error) {
	var stoppers []gracefull.ProcessStopper
	stopper := func(ctx context.Context) error {
		for _, st := range stoppers {
			err := st(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}

	f := Factory{}

	// db
	stopperDb, err := db.Init(cfg)
	if err != nil {
		panic(err)
	}
	stoppers = append(stoppers, stopperDb)
	conn, err := db.GetConnection(cfg.App.Connection)
	if err != nil {
		return f, stopper, err
	}

	// migration
	automigration.Init(cfg)

	// seeder
	seeder.Init(cfg)

	// ws
	f.WsHub = ws.NewHub()

	// firestore
	stopperFs, err := firebase.FirestoreInit()
	stoppers = append(stoppers, stopperFs)
	fsClient := firebase.GetFirestoreClient()
	if err != nil {
		return f, stopper, err
	}

	// repository
	f.Repository = repository.Init(cfg, conn, fsClient)

	// usecase
	f.Usecase = usecase.Init(cfg, f.Repository)

	return f, stopper, nil
}
