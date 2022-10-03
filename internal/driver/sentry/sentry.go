package sentry

import (
	"fmt"
	"log"
	"time"

	sentry "github.com/getsentry/sentry-go"
	"gitlab.com/stoqu/stoqu-be/internal/config"
)

func Init() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.Config.Driver.Sentry.Dsn,
		TracesSampleRate: 1.0,
		ServerName:       fmt.Sprintf("%s v%s", config.Config.App.Name, config.Config.App.Version),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sentry.Flush(2 * time.Second)
}
