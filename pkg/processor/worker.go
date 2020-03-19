package processor

import (
	"context"
	"errors"
	"sync"

	"github.com/MovieStoreGuy/detector-doctor/pkg/checks"
	"github.com/MovieStoreGuy/detector-doctor/pkg/client"
	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

var (
	// ErrJobNotComplete returns when done is not set within a worker
	ErrJobNotComplete = errors.New("job has not completed")
)

type worker struct {
	rwlock  sync.RWMutex
	done    bool
	issue   error
	results []*types.Result

	client *client.SignalFx

	// Runner is the unique method to query facts regarding the detector
	runner checks.Check
}

func newWorker(sfx *client.SignalFx, f checks.Check) *worker {
	if f == nil {
		panic("function parameter required")
	}
	return &worker{
		client: sfx,
		runner: f,
	}
}

// Work abstracts the knowledge the underlying running is running within its own go routine
func (w *worker) work(ctx context.Context, detectorID string, finished *sync.WaitGroup) {
	if w.runner == nil {
		return
	}

	async := func() {
		defer finished.Done()
		defer w.rwlock.Unlock()
		w.rwlock.Lock()
		w.done = false
		w.results, w.issue = w.runner(ctx, detectorID, w.client)
		w.done = true
	}

	go async()
}

// GetResults will not return any data until internal job is completed
func (w *worker) getResults() ([]*types.Result, error) {
	w.rwlock.RLock()
	defer w.rwlock.RUnlock()
	if !w.done {
		return nil, ErrJobNotComplete
	}
	return w.results, w.issue
}
