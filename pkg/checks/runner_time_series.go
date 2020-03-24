package checks

import (
	"context"
	"fmt"
	"time"

	"github.com/MovieStoreGuy/detector-doctor/pkg/client"
	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

const (
	// MaxResolutionMs the largest resolution period allowed before being
	// too long for someone to effectively monitor on it
	MaxResolutionMs int = 10 * 60 * 60 * 1000 // == 10 minutes
)

// CheckTimeSeries reads the underlying time series data for a detector
func CheckTimeSeries(ctx context.Context, detectorID string, sfx *client.SignalFx) ([]*types.Result, error) {
	det := getGlobalCache().explictWaitGetDetector(detectorID, 5*time.Second)
	if det == nil {
		if sfx == nil {
			return nil, types.ErrMissingClient
		}
		d, err := sfx.GetDetectorByID(ctx, detectorID)
		if err != nil {
			return nil, err
		}
		getGlobalCache().setDetector(detectorID, d)
		det = d
	}
	// No point trying to query metric time series
	// if we know that the number of time series is too much already
	if det.OverMTSLimit {
		return []*types.Result{}, nil
	}
	results := make([]*types.Result, 0)
	messages, timeseries, err := sfx.GetMetricTimeSeries(ctx, det.ProgramText, map[string]interface{}{
		"start":     time.Now().Add(-1 * 2 * time.Hour),
		"immediate": true,
	})
	results = append(results, types.CheckSystemIssue(err == nil, "Reading program text from detector"))
	if err != nil {
		results[0] = results[0].WithMessage(err.Error())
		return results, nil
	}
	for _, msg := range messages {
		meta, cast := msg.(*types.MessageMetadata)
		if !cast {
			continue
		}
		if _, hasField := meta.Properties["sf_isDieQuickly"]; hasField {
			results = append(results,
				types.InformationalIssue("Marked as short lived time series").
					WithMessage(fmt.Sprintf("%s is marked as short lived, can cause problems", meta.Properties["sf_originatingMetric"])),
			)
		}
		if res, cast := meta.Properties["sf_resolutionMs"].(int); cast {
			results = append(results,
				// This is a user issue in the sense that the data is sent non peridoicly
				types.CheckUserIssue(res < MaxResolutionMs, "Check resolution").
					WithMessage("Ensuring the resolution of the time series isn't beyond 10m"),
			)
		}
	}
	counts := make(map[string]uint64, 0)
	totalDatapoints := 0
	for _, data := range timeseries {
		// Need to resolve the issue of how the time series are being processed
		for _, point := range data.Data {
			counts[point.TimeSeriesID] = counts[point.TimeSeriesID] + 1
			totalDatapoints++
		}
	}
	expectedLength := uint64(totalDatapoints / len(counts))
	for tsid, count := range counts {
		results = append(results,
			types.CheckInformational(count == expectedLength, "time series count").
				WithMessage(fmt.Sprintf("time series id: %s has expected number of values", tsid)),
		)
	}
	return results, nil
}
