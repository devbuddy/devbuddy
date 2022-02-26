package taskapi

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/context"
)

type TaskAction interface {
	Description() string
	Needed(*context.Context) *ActionResult
	Run(*context.Context) error
	Feature() *autoenv.FeatureInfo
}

type ActionResult struct {
	Needed bool
	Reason string
	Error  error
}

func Failed(errorMessage string, args ...interface{}) *ActionResult {
	return &ActionResult{Error: fmt.Errorf(errorMessage, args...)}
}

func Needed(message string, args ...interface{}) *ActionResult {
	return &ActionResult{Needed: true, Reason: fmt.Sprintf(message, args...)}
}

func NotNeeded() *ActionResult {
	return &ActionResult{Needed: false}
}
