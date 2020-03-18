package printer

import (
	"fmt"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

func textPrinter(detectorID string, results []*types.Result) error {
	fmt.Printf("Detector ID: %s", detectorID)
	for _, result := range results {
		fmt.Printf("Error:%s, Tested:%s, Message:%s", result.IssueType.String(), result.Tested, result.Msg)
	}
	return nil
}
