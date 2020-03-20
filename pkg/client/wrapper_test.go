package client

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	accessToken = os.Getenv("SFX_API_TOKEN")
	realm = os.Getenv("SFX_REALM")
)

func TestWebsockConnection(t *testing.T) {
	if accessToken == "" {
		t.Skip("Skipping test as no token is provided")
	}
	t.Parallel()
	sfx := NewSignalFxClient(realm, accessToken, NewConfiguredClient())

	text := `
# Using org values so this can be run against any account without issue
# also ensuring that it can handle comments within the text
A = data('sf.org.numDatapointsReceived').mean(over='1m').scale(60).sum().publish(label='DPM Received')
	`
	now := time.Now()
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(4 * time.Second))
	err := sfx.readStreamData(ctx, text, map[string]interface{}{
		"start": now.Add(-1 * 1 * time.Hour).Unix(),
		"stop" : now.Unix(),
		"timezone": "utc",
	})
	require.NoError(t, err)
} 
