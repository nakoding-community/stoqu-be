package account_lock

import (
	"context"
	"crypto/rand"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"math/big"
	"time"
)

const (
	defaultExpiry = time.Second * 8
)

type (
	Producer interface {
		// New will create an account mutex with given TTL (or expiry duration)
		// the TTL is optional, if the value is zero or less it will use default
		// value: 8 seconds.
		New(ctx context.Context, accountID string, TTL time.Duration) Lock
	}
	producer struct {
		*redsync.Redsync
		cfg struct {
			lock struct {
				maxTries   int
				retryDelay redsync.DelayFunc
			}
		}
	}
	Option func(producer *producer)
)

func NewProducer(rc *redis.Client, opts ...Option) *producer {
	r := &producer{}
	r.cfg.lock.maxTries = 32
	for _, opt := range opts {
		opt(r)
	}
	r.Redsync = redsync.New(goredis.NewPool(rc))
	return r
}

func (p *producer) New(ctx context.Context, accountID string, TTL time.Duration) Lock {
	opts := []redsync.Option{
		redsync.WithTries(p.cfg.lock.maxTries),
	}
	if TTL <= 0 {
		TTL = defaultExpiry
	}
	opts = append(opts, redsync.WithExpiry(TTL))
	if p.cfg.lock.retryDelay != nil {
		opts = append(opts, redsync.WithRetryDelayFunc(p.cfg.lock.retryDelay))
	}
	m := p.NewMutex(accountID, opts...)
	return newLock(TTL, m)
}

// WithMaxTries will override the maximum number of tries the lock will use
// when attempting the locking mechanism. default value for maximum tries is 32
func WithMaxTries(maxTries int) Option {
	return func(r *producer) {
		r.cfg.lock.maxTries = maxTries
	}
}

// WithRetryDelayWithJitter will override the delay mechanism between each retry locking attempt
// it will use jitter by random number between 0 and the current retry attempt
// thus the retry duration will be random number between min and max plus the random jitter
func WithRetryDelayWithJitter(min, max time.Duration) Option {
	return func(r *producer) {
		bimax := big.NewInt(int64(max))
		r.cfg.lock.retryDelay = func(tries int) time.Duration {
			j, _ := rand.Int(rand.Reader, big.NewInt(int64(tries)))
			n, _ := rand.Int(rand.Reader, bimax)
			// the duration is between min and max + random jitter (j)
			return time.Duration(n.Int64()) + min + time.Duration(j.Int64())
		}
	}
}
