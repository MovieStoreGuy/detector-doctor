package types

// Issue is used as an enum value
type Issue int8

func (i Issue) String() string {
	s, exist := issuesStrings[i]
	if !exist {
		return "Unknown Issue"
	}
	return s
}

var issuesStrings = map[Issue]string{
	OK:     "OK",
	System: "System",
	User:   "User",
}

const (
	// OK is used when no issue is detected
	OK Issue = iota + 1
	// System is used when the issue relates to an issue with SignalFx
	System
	// User is used when the issue relates to an issue that is result of user action
	User
	// Informational is used when as guide to highlight potential issue
	Informational
)
