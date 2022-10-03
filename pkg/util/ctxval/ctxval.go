package ctxval

import (
	"context"

	"gorm.io/gorm"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
)

type key string

var (
	keyAuth key = "x-auth"
	keyTrx  key = "x-trx"
)

func SetAuthValue(ctx context.Context, payload *abstraction.AuthContext) context.Context {
	return context.WithValue(ctx, keyAuth, payload)
}

func GetAuthValue(ctx context.Context) *abstraction.AuthContext {
	val, ok := ctx.Value(keyAuth).(*abstraction.AuthContext)
	if ok {
		return val
	}
	return nil
}

func SetTrxValue(ctx context.Context, trx *gorm.DB) context.Context {
	return context.WithValue(ctx, keyTrx, trx)
}

func GetTrxValue(ctx context.Context) *gorm.DB {
	val, ok := ctx.Value(keyTrx).(*gorm.DB)
	if ok {
		return val
	}

	return nil
}
