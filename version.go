package main

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	// Version is the git tag stored against the project
	Version = "unknown"

	// GitHash is the git hash stored with this version
	GitHash = "unkown"
)

// GetRuntimeVersions will show the current version of the application
func GetRuntimeVersions() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Version:%s\n", Version))
	sb.WriteString(fmt.Sprintf("GitHash:%s\n", GitHash))
	sb.WriteString(fmt.Sprintf("Go version:%v", runtime.Version()))
	return sb.String()
}
