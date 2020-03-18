package checks

import (
	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

// CheckDetector inspects the settings of the detector to see if there is any user
// settings that have been set that could of caused issues with detector.
func CheckDetector(det *types.Detector) ([]*types.Result, error) {
	if det == nil {
		return nil, types.ErrNoDetectorFound
	}
	results := []*types.Result{
		types.CheckUserIssue(det.OverMTSLimit, "Over MTS limit").
			WithMessage("Reduce the number of time series by applying further filtering"),
		// Would like to follow up with their support to understand when a detector can be locked
		types.CheckSystemIssue(!det.Locked, "Locked").
			WithMessage("detector is locked by an unknown reason, detector can not be updated"),
		types.CheckUserIssue(len(det.ProgramText) > 0, "Has program text").
			WithMessage("Ensure that there is Signalflow set for this detector"),
		types.CheckUserIssue(len(det.Rules) > 0, "Has alert rules").
			WithMessage("Ensure that there is alerting rules set for this detector"),
	}
	for _, r := range det.Rules {
		results = append(results,
			types.CheckUserIssue(r.Disabled, "disabled rule").
				WithMessage("Ensuring rules are not disabled"),
			types.CheckUserIssue(len(r.Notifications) > 0, "Notification rules set").
				WithMessage("Ensure that there is a notification path"),
		)
	}
	return results, nil
}
