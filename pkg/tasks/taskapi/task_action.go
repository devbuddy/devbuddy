package taskapi

import "fmt"

type taskAction interface {
	Description() string
	Needed(*Context) *ActionResult
	Run(*Context) error
}

type ActionResult struct {
	Needed bool
	Reason string
	Error  error
}

func ActionFailed(errorMessage string, args ...interface{}) *ActionResult {
	return &ActionResult{Error: fmt.Errorf(errorMessage, args...)}
}

func ActionNeeded(message string, args ...interface{}) *ActionResult {
	return &ActionResult{Needed: true, Reason: fmt.Sprintf(message, args...)}
}

func ActionNotNeeded() *ActionResult {
	return &ActionResult{Needed: false}
}
