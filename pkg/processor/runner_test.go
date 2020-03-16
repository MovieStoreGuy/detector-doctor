package processor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyService(t *testing.T) {
	r := &runner{}
	results, err := r.Run(context.TODO(), "")
	require.NoError(t, err)
	require.NotNil(t, results)
}
