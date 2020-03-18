package checks_test

import (
	"testing"

	"github.com/MovieStoreGuy/detector-doctor/pkg/checks"
	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestMissingDetector(t *testing.T) {
	t.Parallel()
	_, err := checks.CheckDetector(nil)
	require.Error(t, err)
}

func TestUserIssueDetector(t *testing.T) {
	t.Parallel()
	det := &types.Detector{
		Locked:       false,
		OverMTSLimit: true,
	}
	results, err := checks.CheckDetector(det)
	require.NoError(t, err)
	require.NotNil(t, results)
	for _, res := range results {
		require.NotEqual(t, res.IssueType, types.System)
	}
}

func TestIncompleteRules(t *testing.T) {
	t.Parallel()
	det := &types.Detector{
		Rules: []types.Rule{
			types.Rule{},
			types.Rule{},
		},
	}
	results, err := checks.CheckDetector(det)
	require.NoError(t, err)
	require.NotNil(t, results)
	for _, res := range results {
		require.NotEqual(t, res.IssueType, types.System)
	}
}
