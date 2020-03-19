package printer

import (
	"fmt"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

func textPrinter(detectorID string, results []*types.Result) error {
	fmt.Println()
	fmt.Println("Detector ID: ", detectorID)
	for _, result := range results {
		fmt.Printf("Error:%s, Tested:%s, Message:%s\n", result.IssueType.String(), result.Tested, result.Msg)
	}
	return nil
}
