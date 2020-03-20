package checks

import (
	"context"

	"github.com/MovieStoreGuy/detector-doctor/pkg/client"
	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

// Check is the typed function to make it easier to update the function definition
type Check func(ctx context.Context, detectorID string, sfx *client.SignalFx) ([]*types.Result, error)
