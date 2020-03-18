package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/MovieStoreGuy/detector-doctor/pkg/client"
	"github.com/MovieStoreGuy/detector-doctor/pkg/printer"
	"github.com/MovieStoreGuy/detector-doctor/pkg/processor"
)

func main() {
	flag.Parse()

	if paramVersion {
		fmt.Println(GetRuntimeVersions())
	}

	writer, err := printer.GetConfiguredPrinter(paramPrinter)
	if err != nil {
		fmt.Println("Unknown output format supplied", paramPrinter)
		fmt.Println("Supported output formats are", printer.ConfiguredPrinters())
		os.Exit(1)
	}

	proc := processor.NewDefaultService(client.NewSignalFxClient(
		paramRealm,
		paramToken,
		client.NewConfiguredClient(),
	))

	for _, detectorID := range flag.Args() {
		results, err := proc.Run(context.Background(), detectorID)
		if err != nil {
			// something bad has happened, figure out what to do later
			os.Exit(1)
		}
		// Results Printer that shows the results for a given detector ID
		writer(detectorID, results)
	}
}
