package printer

import (
	"fmt"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

func textPrinter(detectorID string, res []*types.Result, filters ...Filter) error {
	results := res
	for _, f := range filters {
		results = f(results)
	}
	fmt.Println()
	fmt.Println("Detector ID: ", detectorID)
	for _, result := range results {
		fmt.Printf("Error:%s, Tested:%s, Message:%s\n", result.IssueType.String(), result.Tested, result.Msg)
	}
	return nil
}
