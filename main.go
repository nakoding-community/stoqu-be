package main

import (
	"sync"
	"time"

	"gitlab.com/stoqu/stoqu-be/internal/config"
	"gitlab.com/stoqu/stoqu-be/internal/driver/cron"
	"gitlab.com/stoqu/stoqu-be/internal/driver/http"
	"gitlab.com/stoqu/stoqu-be/internal/factory"
	"gitlab.com/stoqu/stoqu-be/pkg/util/gracefull"
)

// @title stoqu-be
// @version 0.0.1
// @description This is a doc for stoqu-be.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @host localhost:3000
// @BasePath /
func main() {
	cfg, err := config.Load("")
	if err != nil {
		panic(err)
	}

	// factory
	wg := new(sync.WaitGroup)
	f, stopperFactory, err := factory.Init(cfg)
	if err != nil {
		panic(err)
	}

	starterApi, stopperApi := http.Init(cfg, f)
	cron.Init(cfg, f)

	wg.Add(1)
	go func() {
		gracefull.StartProcessAtBackground(starterApi)
		gracefull.StopProcessAtBackground(time.Second*10, stopperApi, stopperFactory)
		wg.Done()
	}()

	wg.Wait()
}
