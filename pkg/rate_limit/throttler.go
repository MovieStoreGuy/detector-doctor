package ratelimit

import (
	"context"
	"errors"
	"math/rand"
	"sync/atomic"
	"time"

	"golang.org/x/sync/semaphore"
)

func init() {
	// Ensuring rand is random based on our start time
	rand.Seed(time.Now().UnixNano())
}

// Throttler allows for up to fix number of actions to happen within a
// fixed interval with accounting for jitters
type Throttler struct {
	limit             int64
	consumed, waiting int64
	interval, jitter  time.Duration

	bucket  *semaphore.Weighted
	running bool
}

// NewThrottler will configure a new rater limiting bucket with max of limit / (interval + [0,jitter)) ==> rate
// Once the limit has been breached, it will pause computation for that thread until the interval + jitter amount
// has been reached
func NewThrottler(limit int64, interval, jitter time.Duration) *Throttler {
	if limit < 0 {
		limit = 1
	}
	if interval < 0 {
		interval = 100 * time.Millisecond
	}
	if jitter < 0 {
		jitter = 0
	}
	return &Throttler{
		limit:    limit,
		interval: interval,
		jitter:   jitter,
		bucket:   semaphore.NewWeighted(limit),
	}
}

func withJitter(ticker *time.Timer, jiter time.Duration) <-chan time.Time {
	if jiter < 0 {
		return ticker.C
	}
	delay := rand.Int63n(int64(jiter))
	time.Sleep(time.Duration(delay))
	return ticker.C
}

// Start runs a background go routine that was reset the sync limit
// based on the interval and jitter time
func (lb *Throttler) Start() {
	if lb.running {
		return
	}
	lb.running = true
	lb.reset()
	go func() {
		timer := time.NewTimer(lb.interval)
		defer timer.Stop()
		for {
			if !lb.running {
				return
			}
			select {
			case <-withJitter(timer, lb.jitter):
				lb.reset()
			}
		}
	}()
}

// Stop will stop the background task to reset the bucket from running
// and will release anyone waiting for the
func (lb *Throttler) Stop() {
	if !lb.running {
		return
	}
	lb.running = false
	// Ensure all waiting proccesses are being let through
	for lb.waiting > 0 {
		lb.bucket.Release(lb.consumed)
	}
}

// Consume will take one place in the internal semiphore and will
// block once we have reached out limit and delay the processing of blocked
// this does mean that blocked will be delayed instead of dropped.
func (lb *Throttler) Consume(ctx context.Context) error {
	if !lb.running {
		return errors.New("not started")
	}
	atomic.AddInt64(&lb.waiting, 1)
	if err := lb.bucket.Acquire(ctx, 1); err != nil {
		return err
	}
	atomic.AddInt64(&lb.waiting, -1)
	atomic.AddInt64(&lb.consumed, 1)
	return nil
}

func (lb *Throttler) reset() {
	lb.bucket.Release(lb.consumed)
	atomic.StoreInt64(&lb.consumed, 0)
}
