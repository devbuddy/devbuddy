package tasks

import "fmt"

type taskAction interface {
	description() string
	Needed(*Context) *ActionResult
	Run(*Context) error
}

type ActionResult struct {
	Needed bool
	Reason string
	Error  error
}

func actionFailed(errorMessage string, args ...interface{}) *ActionResult {
	return &ActionResult{Error: fmt.Errorf(errorMessage, args...)}
}

func actionNeeded(message string, args ...interface{}) *ActionResult {
	return &ActionResult{Needed: true, Reason: fmt.Sprintf(message, args...)}
}

func actionNotNeeded() *ActionResult {
	return &ActionResult{Needed: false}
}

func runAction(ctx *Context, action taskAction) error {
	desc := action.description()

	result := action.Needed(ctx)
	if result.Error != nil {
		return fmt.Errorf("The task action (%s) failed to detect whether it need to run: %s", desc, result.Error)
	}

	if result.Needed {
		if desc != "" {
			ctx.ui.TaskActionHeader(desc)
		}
		ctx.ui.Debug("Reason: \"%s\"", result.Reason)

		err := action.Run(ctx)
		if err != nil {
			return fmt.Errorf("The task action failed to run: %s", err)
		}

		result = action.Needed(ctx)
		if result.Error != nil {
			return fmt.Errorf("The task action failed to detect if it is resolved: %s", result.Error)
		}

		if result.Needed {
			return fmt.Errorf("The task action did not produce the expected result: %s", result.Reason)
		}
	}

	return nil
}
