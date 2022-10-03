package trxmanager

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/stoqu/stoqu-be/pkg/util/ctxval"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type trxManager struct {
	db *gorm.DB
}

type trxFn func(ctx context.Context) error

func New(db *gorm.DB) *trxManager {
	return &trxManager{db}
}

func (g *trxManager) WithTrx(ctx context.Context, fn trxFn) (err error) {
	trx := g.db.Begin()
	ctx = ctxval.SetTrxValue(ctx, trx)

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			trx.Rollback()
			logrus.Error(p)
			err = errors.New("panic happened because: " + fmt.Sprintf("%v", p))
		} else if err != nil {
			// error occurred, rollback
			trx.Rollback()
		} else {
			// all good, commit
			err = trx.Commit().Error
		}
	}()

	err = fn(ctx)
	return err
}
