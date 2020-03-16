package processor

import (
	"context"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

// Service defines the dispatcher that will issue jobs to the runnning workers and
// return the results of all issues found
type Service interface {

	// Run will spawn all the workers and get all the results returned.
	Run(ctx context.Context, detectorId string) ([]*types.Result, error)
}
