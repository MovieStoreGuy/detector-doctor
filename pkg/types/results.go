package types

type Issue int8

const (
	OK Issue = iota
	System
	User
)

const (
	DefaultNoIssueExpression     = "No required action to take"
	DefaultUserIssueExpresison   = "There is a issue with this detector due to user error"
	DefaultSystemIssueExpression = "There is an unknown issue with this detector that can not be manually resolved"
)

type Result struct {
	IssueType Issue
	Tested    string
	Msg       string
}

func (r *Result) WithMessage(msg string) *Result {
	return &Result{
		IssueType: r.IssueType,
		Tested:    r.Tested,
		Msg:       msg,
	}
}

func NoIssue(tested string) *Result {
	return &Result{
		IssueType: OK,
		Tested:    tested,
		Msg:       DefaultNoIssueExpression,
	}
}

func UserIssue(tested string) *Result {
	return &Result{
		IssueType: User,
		Tested:    tested,
		Msg:       DefaultUserIssueExpresison,
	}
}

func SystemIssue(tested string) *Result {
	return &Result{
		IssueType: System,
		Tested:    tested,
		Msg:       DefaultSystemIssueExpression,
	}
}
