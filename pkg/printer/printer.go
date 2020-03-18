package printer

import (
	"fmt"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

var (
	printers = map[string]Printer{
		"text": textPrinter,
		"json": jsonPrinter,
	}
)

// Printer defines the expected printer and will write results
type Printer func(detectorID string, results []*types.Result) error

// ConfiguredPrinters returns a list of all currently configured printers that can be used
func ConfiguredPrinters() []string {
	results := make([]string, 0)
	for p := range printers {
		results = append(results, p)
	}
	return results
}

// GetConfiguredPrinter will return the printer requested
func GetConfiguredPrinter(printer string) (Printer, error) {
	p, exist := printers[printer]
	if !exist {
		return nil, fmt.Errorf("no configured printer %s", printer)
	}
	return p, nil
}
