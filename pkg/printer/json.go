package printer

import (
	"encoding/json"
	"fmt"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

func jsonPrinter(detectorID string, results []*types.Result) error {
	data := map[string]interface{}{
		"DetectorID": detectorID,
		"Results":    results,
	}
	buff, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Println(string(buff))
	return nil
}
