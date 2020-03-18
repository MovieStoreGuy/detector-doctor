package types

import "errors"

var (
	// ErrNoDetectorFound is return when no valid detector is given
	ErrNoDetectorFound = errors.New("no detector passed")

	ErrAPIIssue = errors.New("unable to process request")

	ErrNotImplemented = errors.New("currently not implemented")
)