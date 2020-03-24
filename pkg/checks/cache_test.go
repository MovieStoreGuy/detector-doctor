package checks

import (
	"testing"
	"time"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestExplictWait(t *testing.T) {
	t.Parallel()
	require.Nil(t, getGlobalCache().explictWaitGetDetector("never-exist", time.Second))

	// Slow write to the cache while on contention
	write := func() {
		tick := time.NewTicker(120 * time.Millisecond)
		<-tick.C
		getGlobalCache().setDetector("delayed-detector", &types.Detector{})
	}
	go write()
	require.NotNil(t, getGlobalCache().explictWaitGetDetector("delayed-detector", 1*time.Second))
	// Should already be in the cache so there should be no wait time
	require.NotNil(t, getGlobalCache().explictWaitGetDetector("delayed-detector", 1*time.Second))
}
