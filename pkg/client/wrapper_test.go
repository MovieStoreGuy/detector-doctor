package client

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
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
	messages, datapoints, err := sfx.GetMetricTimeSeries(ctx, text, map[string]interface{}{
		"start":     toUnixMilliseconds(now.Add(-1 * 10 * time.Minute)),
		"immediate": true,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(datapoints), 1)
	require.GreaterOrEqual(t, len(messages), 1)
	for _, msg := range messages {
		meta, cast := msg.(*types.MessageMetadata)
		if !cast {
			continue
		}
		t.Logf("%+v", meta)
	}
}

func TestMultipleDatastreams(t *testing.T) {
	if accessToken == "" {
		t.Skip("Skipping test as no token is provided")
	}
	sfx := NewSignalFxClient(realm, accessToken, NewConfiguredClient())

	text := `# Using org values so this can be run against any account without issue
# also ensuring that it can handle comments within the text
A = data('sf.org.numDatapointsReceived').mean(over='1m').scale(60).sum().publish(label='DPM Received')
B = data('sf.org.subscription.datapointsPerMinute').publish(label='DPM Limit')
detect(when(A.sum() > B.sum(), '2m')).publish('Over DPM limit')`
	now := time.Now().UTC()
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(4*time.Second))
	messages, datapoints, err := sfx.GetMetricTimeSeries(ctx, text, map[string]interface{}{
		"start":     now.Add(-1 * 10 * time.Minute),
		"immediate": true,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(datapoints), 1)
	require.GreaterOrEqual(t, len(messages), 1)
	for _, msg := range messages {
		meta, cast := msg.(*types.MessageMetadata)
		if !cast {
			continue
		}
		t.Logf("%+v", meta)
	}
}
