package types

const (
	defaultNoIssueExpression     = "No required action to take"
	defaultUserIssueExpresison   = "There is a issue with this detector due to user error"
	defaultSystemIssueExpression = "There is an unknown issue with this detector that can not be manually resolved"
)

// Result used when evaluation a known condition for a detector to be in an errored state
type Result struct {
	IssueType Issue  `json:"issueType"`
	Tested    string `json:"tested"`
	Msg       string `json:"message"`
}

// WithMessage updates the Result's message and returns a deepcopy instead of
// modifying the orginal result
func (r *Result) WithMessage(msg string) *Result {
	return &Result{
		IssueType: r.IssueType,
		Tested:    r.Tested,
		Msg:       msg,
	}
}

// CheckUserIssue is convience function to return an UserIssue iff ok is false
func CheckUserIssue(ok bool, tested string) *Result {
	if ok {
		return NoIssue(tested)
	}
	return UserIssue(tested)
}

// CheckSystemIssue is a convience function to return a SystemIssue iff ok is false
func CheckSystemIssue(ok bool, tested string) *Result {
	if ok {
		return NoIssue(tested)
	}
	return SystemIssue(tested)
}

// NoIssue returns a Result with its IssueType set as OK and default message
func NoIssue(tested string) *Result {
	return &Result{
		IssueType: OK,
		Tested:    tested,
		Msg:       defaultNoIssueExpression,
	}
}

// UserIssue returns a Result with its IssueType set as User and a default message
func UserIssue(tested string) *Result {
	return &Result{
		IssueType: User,
		Tested:    tested,
		Msg:       defaultUserIssueExpresison,
	}
}

// SystemIssue returns a Result with its IssueType set as System and a default message
func SystemIssue(tested string) *Result {
	return &Result{
		IssueType: System,
		Tested:    tested,
		Msg:       defaultSystemIssueExpression,
	}
}
