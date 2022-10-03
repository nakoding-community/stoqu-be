package account_lock

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"sync/atomic"
	"testing"
	"time"
)

type AccountLockSuite struct {
	suite.Suite
	rc *redis.Client
	mr *miniredis.Miniredis
}

func (s *AccountLockSuite) SetupSuite() {
	mr, err := miniredis.Run()
	s.Require().NoError(err)
	s.rc = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	s.mr = mr
}

func (s *AccountLockSuite) TearDownSuite() {
	s.rc.Close()
	s.mr.Close()
}

func (s *AccountLockSuite) TestProducerWithOptions() {
	fixtures := []struct {
		maxRetry                 int
		min, max                 time.Duration
		wantMaxRetry             int
		wantRetryDelayFuncNotNil bool
	}{
		{0, 0, 0, 0, false},
		{1, 0, 0, 1, false},
		{0, 10 * time.Second, 20 * time.Second, 0, true},
	}
	for _, f := range fixtures {
		f := f
		s.T().Run(fmt.Sprintf("%v", f), func(t *testing.T) {
			t.Parallel()
			max := WithMaxTries(f.maxRetry)
			var retryDelay Option
			if f.min > 0 && f.max > 0 {
				retryDelay = WithRetryDelayWithJitter(f.min, f.max)
			}
			opts := []Option{max}
			if retryDelay != nil {
				opts = append(opts, retryDelay)
			}
			lp := NewProducer(s.rc, opts...)
			assert.NotNil(t, lp)
			assert.Equal(t, f.wantMaxRetry, lp.cfg.lock.maxTries)
			assert.Equal(t, f.wantRetryDelayFuncNotNil, lp.cfg.lock.retryDelay != nil)
		})

	}
}

func (s *AccountLockSuite) TestAutoExtendLock() {
	milli := time.Millisecond
	lp := NewProducer(s.rc, WithMaxTries(3)) // just use max retry 3 to prevent the test taking long time to complete
	fixtures := []struct {
		lockName   string
		lockExpiry time.Duration
		wantExpiry time.Duration
	}{
		{"test", 500 * milli, 500 * milli},
		{"test_zero", 0, defaultExpiry},
		{"test_1sec", 1000 * milli, 1000 * milli},
	}
	for _, f := range fixtures {
		f := f
		s.T().Run(f.lockName, func(t *testing.T) {
			t.Parallel()
			l := lp.New(context.Background(), f.lockName, f.lockExpiry).(*lock)
			assert.Equal(t, l.expiredAfter, f.wantExpiry)
			assert.False(t, l.Valid()) // the lock must not be valid before aquired
			assert.NoError(t, l.Acquire(context.Background()))
			time.Sleep(l.expiredAfter + (l.expiredAfter / 5))  // 1.2 of the lock expire
			assert.True(t, l.Valid())                          // the lock must still be valid
			assert.NoError(t, l.Acquire(context.Background())) // the lock must be auto extended, thus aquire must not be failed but return immediately
			// test release can be called multiple time without error
			for i := 0; i < 2; i++ {
				bool, err := l.Release(context.Background())
				if i == 0 {
					assert.True(t, bool)
				} else {
					assert.False(t, bool)
				}
				assert.NoError(t, err)
			}
			// make sure we can aquire the same lock after it released
			assert.NoError(t, l.Acquire(context.Background()))
			_, err := l.Release(context.Background())
			assert.NoError(t, err)
		})
	}
}

func (s *AccountLockSuite) TestConcurrentAccess() {
	lp := NewProducer(s.rc, WithMaxTries(3), WithRetryDelayWithJitter(10*time.Nanosecond, time.Second)) // just use max retry 3 to prevent the test taking long time to complete
	fixtures := []struct {
		operation      string
		numbGoRoutines int
		lockExpiry     time.Duration
	}{
		{"acquire", 4, 10 * time.Second},
		{"release", 4, 500 * time.Millisecond},
		{"acquire", 100, 10 * time.Second},
		{"acquire", 25, 10 * time.Second},
		{"release", 100, 500 * time.Millisecond},
		{"release", 25, 500 * time.Millisecond},
	}
	for idx, f := range fixtures {
		idx := idx
		f := f
		s.T().Run(fmt.Sprintf("%s#%d", f.operation, idx), func(t *testing.T) {
			t.Parallel()
			l := lp.New(context.Background(), fmt.Sprintf("test#%d", idx), f.lockExpiry).(*lock)
			var endTimes []<-chan time.Time
			var opCalled uint32
			// TODO(bistokdl) is there better way to determine
			// whether the test only call the operation once?
			switch f.operation {
			case "release":
				l.releasedCallBack = func() {
					atomic.AddUint32(&opCalled, 1)
				}
			default:
				l.acquiredCallBack = func() {
					atomic.AddUint32(&opCalled, 1)
				}
			}

			//pre test
			if f.operation == "release" {
				assert.NoError(t, l.Acquire(context.Background()))
			}
			for i := 0; i < f.numbGoRoutines; i++ {
				endTime := make(chan time.Time, 1)
				endTimes = append(endTimes, endTime)
				go func() {
					switch f.operation {
					case "acquire":
						err := l.Acquire(context.Background())
						endTime <- time.Now()
						assert.NoError(t, err)
						assert.True(t, l.Valid())
					default:
						assert.True(t, l.Valid())
						_, err := l.Release(context.Background())
						endTime <- time.Now()
						assert.NoError(t, err)
					}
				}()
			}
			// post test
			if f.operation == "acquire" {
				_, _ = l.Release(context.Background())
			}
			prev := <-endTimes[0]
			for i := 1; i < len(endTimes); i++ {
				curr := <-endTimes[i]
				if !assert.True(t, opCalled == 1 && curr.Unix() == prev.Unix(), "idx: %d, prev: %v, curr: %v", i, prev, curr) {
					t.Logf("operation called count: %d, prev (ms: %d, s: %d), curr (ms: %d, s: %d)", opCalled, prev.UnixMilli(), prev.Unix(), curr.UnixMilli(), curr.Unix())
				}
				prev = curr
			}
		})
	}
}

func TestAccountLockSuite(t *testing.T) {
	suite.Run(t, new(AccountLockSuite))
}
