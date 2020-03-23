package client

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	accessToken string = os.Getenv("SFX_API_TOKEN")
	realm       string = os.Getenv("SFX_REALM")
)

func TestWebsockConnection(t *testing.T) {
	if accessToken == "" {
		t.Skip("Skipping test as no token is provided")
	}
	sfx := NewSignalFxClient(realm, accessToken, NewConfiguredClient())

	text := `# Using org values so this can be run against any account without issue
# also ensuring that it can handle comments within the text
A = data('sf.org.numDatapointsReceived').mean(over='1m').scale(60).sum().publish(label='DPM Received')`
	now := time.Now().UTC()
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(4*time.Second))
	_, datapoints, err := sfx.readStreamData(ctx, text, map[string]interface{}{
		"start":     toUnixMilliseconds(now.Add(-1 * 10 * time.Minute)),
		"immediate": true,
		"timezone":  "UTC",
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(datapoints), 1)
	t.Log(datapoints[0])
}

func TestMultipleDatastreams(t *testing.T) {
	if accessToken == "" {
		t.Skip("Skipping test as no token is provided")
	}
	sfx := NewSignalFxClient(realm, accessToken, NewConfiguredClient())

	text := `# Using org values so this can be run against any account without issue
# also ensuring that it can handle comments within the text
A = data('sf.org.numDatapointsReceived').mean(over='1m').scale(60).sum().publish(label='DPM Received')
B = data('sf.org.subscription.datapointsPerMinute').publish(label='DPM Limit')`
	now := time.Now().UTC()
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(4*time.Second))
	_, datapoints, err := sfx.readStreamData(ctx, text, map[string]interface{}{
		"start":     toUnixMilliseconds(now.Add(-1 * 10 * time.Minute)),
		"immediate": true,
		"timezone":  "UTC",
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(datapoints), 1)
	t.Log(datapoints[0])
}
