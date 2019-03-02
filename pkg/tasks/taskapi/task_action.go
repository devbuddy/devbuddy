package taskapi

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
)

type TaskAction interface {
	Description() string
	Needed(*Context) *ActionResult
	Run(*Context) error
	Feature() *autoenv.FeatureInfo
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
