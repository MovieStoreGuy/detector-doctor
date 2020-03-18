package types

const (
	DefaultNoIssueExpression     = "No required action to take"
	DefaultUserIssueExpresison   = "There is a issue with this detector due to user error"
	DefaultSystemIssueExpression = "There is an unknown issue with this detector that can not be manually resolved"
)

type Result struct {
	IssueType Issue  `json:"issueType"`
	Tested    string `json:"tested"`
	Msg       string `json:"message"`
}

func (r *Result) WithMessage(msg string) *Result {
	return &Result{
		IssueType: r.IssueType,
		Tested:    r.Tested,
		Msg:       msg,
	}
}

func CheckUserIssue(ok bool, tested string) *Result {
	if ok {
		return NoIssue(tested)
	}
	return UserIssue(tested)
}

func CheckSystemIssue(ok bool, tested string) *Result {
	if ok {
		return NoIssue(tested)
	}
	return SystemIssue(tested)
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
