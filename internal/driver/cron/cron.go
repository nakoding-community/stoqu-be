package cron

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
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

	// daily
	jobDaily, _ := s.Every(1).Day().At("06:00").Do(func() {
		logrus.Info("Cron reminder stock daily, running ...")
		err := f.Usecase.ReminderStockHistory.GenerateRecurring(context.Background(), constant.REMINDER_STOCK_DAILY)
		if err != nil {
			logrus.Error("Cron reminder stock error", err.Error())
		}
	})
	jobDaily.SingletonMode()

	// monthly
	jobMonthly, _ := s.Every(1).Month().Do(func() {
		logrus.Info("Cron reminder stock monthly, running ...")
		err := f.Usecase.ReminderStockHistory.GenerateRecurring(context.Background(), constant.REMINDER_STOCK_DAILY)
		if err != nil {
			logrus.Error("Cron reminder stock error", err.Error())
		}
	})
	jobMonthly.SingletonMode()

	s.StartAsync()

	return stopper
}
