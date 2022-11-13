package cron

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"gitlab.com/stoqu/stoqu-be/internal/config"
	"gitlab.com/stoqu/stoqu-be/internal/factory"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	"gitlab.com/stoqu/stoqu-be/pkg/util/gracefull"
)

var stopper = gracefull.NilStopper

// !TODO: add stopper
func Init(cfg *config.Configuration, f factory.Factory) gracefull.ProcessStopper {
	if !cfg.Driver.Cron.Enabled {
		return stopper
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	s := gocron.NewScheduler(loc)

	// second_5
	jobSec, _ := s.Every(5).Second().Do(func() {
		f.Usecase.ReminderStockHistory.GenerateRecurring(context.Background(), constant.REMINDER_STOCK_SECOND)
	})
	jobSec.SingletonMode()

	// daily
	jobDaily, _ := s.Every(1).Day().Do(func() {
		f.Usecase.ReminderStockHistory.GenerateRecurring(context.Background(), constant.REMINDER_STOCK_DAILY)
	})
	jobDaily.SingletonMode()

	s.StartAsync()

	return stopper
}
