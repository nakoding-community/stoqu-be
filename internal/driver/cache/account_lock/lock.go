package account_lock

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"sync"
	"sync/atomic"
	"time"
)

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type stopAutoExtendLock func()

type (
	Lock interface {
		// Acquire the lock
		// once a lock is acquired, the lock will be automatically extended
		// until the lock is released, or the context used on Acquire is canceled.
		// the time before extension is 3/4 of the lock TLL (expiry) duration.
		// the lock must be released by the caller, or it will block other process
		// to acquire the same lock until the lock is expired.
		Acquire(context.Context) error
		// Release the lock
		Release(context.Context) (bool, error)
		// Valid to check if the acquired lock is still valid on the moment
		// of this call
		Valid() bool
	}

	lock struct {
		noCopy
		acquiredCallBack func() // for unit test purpose only
		releasedCallBack func() // for unit test purpose only

		expiredAfter   time.Duration
		stopAutoExtend stopAutoExtendLock
		rm             *redsync.Mutex
		mu             sync.RWMutex
		acqOrRel       int32 // 0: released, 1: acquired
	}
)

func newLock(expiry time.Duration, m *redsync.Mutex) *lock {
	l := &lock{rm: m}
	l.expiredAfter = expiry
	return l
}

func autoExtendLock(ctx context.Context, mu *sync.RWMutex, m *redsync.Mutex, dur time.Duration) stopAutoExtendLock {
	stop := make(chan struct{}, 1)
	extendTicker := time.NewTicker(dur)
	once := sync.Once{}
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-ctx.Done():
				extendTicker.Stop()
				return
			case <-extendTicker.C:
				mu.Lock()
				_, err := m.ExtendContext(ctx)
				mu.Unlock()
				if err != nil {
					// TODO log error here
				}
			}
		}
	}()
	return func() {
		once.Do(func() {
			extendTicker.Stop()
			stop <- struct{}{}
		})
	}
}

// Acquire will acquire the lock and run a watchdog process
// that will repeatedly extend the lock expiry time until one of following condition is met:
// - the lock's Release method is called
// - the context used as argument is canceled or timeout
// the watchdog process will extend the lock every ~3/4 of the lock's expire period.
// if the lock is not in acquired state, concurrent call to Acquired (while the process is still being run) will return,
// only after the first call has finished. in other words it mimics the behavior of sync.Once
func (l *lock) Acquire(ctx context.Context) error {
	if atomic.LoadInt32(&l.acqOrRel) == 0 {
		return l.acquireSlow(ctx)
	}
	return nil
}

func (l *lock) acquireSlow(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.acqOrRel == 1 {
		return nil
	}
	var acq int32
	defer func() {
		atomic.StoreInt32(&l.acqOrRel, acq)
		if l.acquiredCallBack != nil {
			l.acquiredCallBack()
		}
	}()
	err := l.rm.LockContext(ctx)
	if err != nil {
		return err
	}
	acq = 1
	d := 3 * l.expiredAfter / 4 // ~ 3/4 of the expiry time
	l.stopAutoExtend = autoExtendLock(ctx, &l.mu, l.rm, d)
	return nil
}

// Release will release the lock and stop the lock's extend watchdog process
// Release will return false if the lock is already released
// if the lock is not in acquired release, concurrent call to Release (while the process is still being run) will return,
// only after the first call has finished. in other words it mimics the behavior of sync.Once
func (l *lock) Release(ctx context.Context) (bool, error) {
	if atomic.LoadInt32(&l.acqOrRel) == 1 {
		return l.releaseSlow(ctx)
	}
	return false, nil
}

func (l *lock) Valid() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.rm.Until().After(time.Now())
}

func (l *lock) releaseSlow(ctx context.Context) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.acqOrRel == 0 {
		return false, nil
	}
	var rel int32 = 1
	defer func() {
		atomic.StoreInt32(&l.acqOrRel, rel)
		if l.releasedCallBack != nil {
			l.releasedCallBack()
		}

	}()
	l.stopAutoExtend()
	res, err := l.rm.UnlockContext(ctx)
	if err != nil {
		return false, err
	}
	rel = 0
	return res, nil

}
