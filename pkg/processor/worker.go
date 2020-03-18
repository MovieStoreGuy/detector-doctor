package processor

import (
	"errors"
	"sync"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

var (
	// ErrJobNotComplete returns when done is not set within a worker
	ErrJobNotComplete = errors.New("job has not completed")
)

type worker struct {
	done    bool
	issue   error
	results []*types.Result

	client interface{}

	// Runner is the unique method to query facts regarding the detector
	runner func(detectorId string, client interface{}) ([]*types.Result, error)
}

func newWorker(f func(string, interface{}) ([]*types.Result, error)) *worker {
	if f == nil {
		panic("function parameter required")
	}
	return &worker{
		runner: f,
	}
}

// Work abstracts the knowledge the underlying running is running within its own go routine
func (w *worker) work(detectorID string, finished *sync.WaitGroup) {
	if w.runner == nil {
		return
	}
	async := func() {
		defer finished.Done()
		w.results, w.issue = w.runner(detectorID, nil)
		w.done = true
	}

	go async()
}

// GetResults will not return any data until internal job is completed
func (w *worker) getResults() ([]*types.Result, error) {
	if !w.done {
		return nil, ErrJobNotComplete
	}
	return w.results, w.issue
}
