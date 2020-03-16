package processor

import (
	"context"
	"sync"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

type runner struct {
	client  interface{}
	workers []worker
}

// NewDefaultService creates a processor that will inspect a
func NewDefaultService(cli interface{}) Service {
	r := &runner{
		client:  cli,
		workers: make([]workers, 0),
	}
	// Add default workers functions in
	return r
}

func (r *runner) Run(ctx context.Context, detectorID string) ([]*types.Result, error) {
	// To spend things up, fetch the currently detector state from SignalFx
	// to avoid over querying the API and risk potentially being rate limited

	// Start all internal workers to fetch the required information for the detector
	finished := &sync.WaitGroup{}
	finished.Add(len(r.workers))
	for _, w := range r.workers {
		w.Work(detectorID, finished)
	}
	finished.Wait()
	results := make([]*types.Result, 0)
	for _, w := range r.workers {
		result, err := w.GetResults()
		switch err {
		case nil:
			// Do nothing we expect it to be nil when we are good
		case ErrJobNotComplete:
			// Should not be able to reach this state since we explictly wait for each job to complete
			panic("reached an unreachable state")
		default:
			// One of the jobs had experienced an issue and we are bubble that error up
			return nil, err
		}
		results = append(results, result...)
	}
	return results, nil
}
