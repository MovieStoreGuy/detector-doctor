package processor

import (
	"context"
	"net/http"
	"sync"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

type runner struct {
	client  *http.Client
	workers []worker
}

// NewDefaultService creates a processor that will inspect a
func NewDefaultService(cli *http.Client) Service {
	r := &runner{
		client:  cli,
		workers: make([]worker, 0),
	}
	// Add default workers functions into the service
	return r
}

func (r *runner) Run(ctx context.Context, detectorID string) ([]*types.Result, error) {
	// To spend things up, fetch the currently detector state from SignalFx
	// to avoid over querying the API and risk potentially being rate limited

	// Start all internal workers to fetch the required information for the detector
	finished := &sync.WaitGroup{}
	finished.Add(len(r.workers))
	for _, w := range r.workers {
		w.work(detectorID, finished)
	}
	finished.Wait()
	results := make([]*types.Result, 0)
	for _, w := range r.workers {
		result, err := w.getResults()
		switch {
		case err == nil:
			// Do nothing we expect it to be nil when we are good
		case err == ErrJobNotComplete:
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
