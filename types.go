package castor

import (
	"fmt"
	"time"
)

// ExitError implements the cli.ExitError interface
type ExitError struct {
	error
	code int
}

// ExitCode returns the error code
func (e ExitError) ExitCode() int {
	return e.code
}

// ExitErr returns an ExitError based on an error code an another error value.
func ExitErr(code int, err error) ExitError {
	switch v := err.(type) {
	// TODO: should we ignore `code` in this case?
	case ExitError:
		return v
	default:
		return ExitError{err, code}
	}
}

// ExitErrorF returns an ExitError based on an error code and a format specifier.
func ExitErrorF(code int, format string, a ...interface{}) ExitError {
	return ExitError{fmt.Errorf(format, a...), code}
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Ref struct {
	Ref string `json:"ref"`
}

type WithLogin struct {
	Login string `json:"login"`
}

type PRsSearch struct {
	IssueCount int        `json:"issueCount"`
	Nodes      []SearchPR `json:"nodes"`
}

type SearchPR struct {
	URL            string         `json:"url"`
	Number         int            `json:"number"`
	Title          string         `json:"title"`
	Author         WithLogin      `json:"author"`
	HeadRefName    string         `json:"headRefName"`
	Labels         Labels         `json:"Labels"`
	ReviewRequests ReviewRequests `json:"reviewRequests"`
}

type LoginAndName struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

type Review struct {
	State       string    `json:"state"`
	SubmittedAt time.Time `json:"submittedAt"`
	URL         string    `json:"url"`
	Author      WithLogin `json:"Author"`
}

type RequestedReviewer struct {
	RequestedReviewer LoginAndName `json:"requestedReviewer"`
}

type ReviewRequests struct {
	TotalCount int                 `json:"totalCount"`
	Nodes      []RequestedReviewer `json:"Nodes"`
}

type Labels struct {
	TotalCount int     `json:"totalCount"`
	Nodes      []Label `json:"nodes"`
}
