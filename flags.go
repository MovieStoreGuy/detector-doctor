package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/MovieStoreGuy/detector-doctor/pkg/printer"
)

const helpMessage = `%s help message

Usage:
$ %s [--flags, ...] detectorIDs, ... 

Description:
This application inspects the state of configured detector to establish if a detector is misbehaving.
It will read the list of detector IDs from the arguments passed to the command line.

Flags:
`

const (
	defaultVersion      = false
	defaultVerbose      = false
	defaultDisableHTTP2 = false
	defaultToken        = ""
	defaultRealm        = "us0"
	defaultOutputFormat = "text"
)

var (
	paramVersion      bool
	paramVerbose      bool
	paramDisableHTTP2 bool
	paramToken        string
	paramRealm        string
	paramPrinter      string
	paramFilters      = printer.NewFlagFilter()
)

func init() {
	flag.BoolVar(&paramVersion, "version", defaultVersion, "shows the current version of the application")
	flag.BoolVar(&paramVerbose, "verbose", defaultVerbose, "enables verbose outout throughout the application")
	flag.BoolVar(&paramDisableHTTP2, "disable-http2", defaultDisableHTTP2, "disables using HTTP 2.0 when making requests")
	flag.StringVar(&paramToken, "token", defaultToken, "sets the token to access the API")
	flag.StringVar(&paramRealm, "realm", defaultRealm, "sets the realm to use when accessing the API")
	flag.StringVar(&paramPrinter, "output-format", defaultOutputFormat, "sets the output format of the results")
	flag.Var(paramFilters, "filter", "used to ")
	flag.Usage = func() {
		app := path.Base(os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), helpMessage, app, app)
		flag.PrintDefaults()
	}
}
