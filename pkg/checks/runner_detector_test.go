package checks

import (
	"context"
	"testing"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
	"github.com/stretchr/testify/require"
)

var (
	notifications = &types.Detector{
		Locked:       false,
		OverMTSLimit: false,
		Rules: []types.Rule{
			types.Rule{},
		},
	}
	userError = &types.Detector{
		OverMTSLimit: true,
	}
	systemIssue = &types.Detector{
		Locked: true,
	}
)

func init() {
	// configure global cache
	getGlobalCache().setDetector("notifications", notifications)
	getGlobalCache().setDetector("user", userError)
	getGlobalCache().setDetector("system", systemIssue)
}

func TestMissingDetector(t *testing.T) {
	t.Parallel()
	_, err := CheckDetector(context.Background(), "", nil)
	require.Error(t, err)
	require.Equal(t, err, types.ErrMissingClient)
}

func TestUserIssueDetector(t *testing.T) {
	t.Parallel()
	results, err := CheckDetector(context.Background(), "user", nil)
	require.NoError(t, err)
	require.NotNil(t, results)
	for _, res := range results {
		require.NotEqual(t, res.IssueType, types.System)
	}
}

func TestIncompleteRules(t *testing.T) {
	t.Parallel()
	results, err := CheckDetector(context.Background(), "notifications", nil)
	require.NoError(t, err)
	require.NotNil(t, results)
	for _, res := range results {
		require.NotEqual(t, res.IssueType, types.System)
	}
}
