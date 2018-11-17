package tasks

import "fmt"

type taskAction interface {
	description() string
	needed(*context) *actionResult
	run(*context) error
}

type actionResult struct {
	Needed bool
	Reason string
	Error  error
}

func actionFailed(errorMessage string, args ...interface{}) *actionResult {
	return &actionResult{Error: fmt.Errorf(errorMessage, args...)}
}

func actionNeeded(message string, args ...interface{}) *actionResult {
	return &actionResult{Needed: true, Reason: fmt.Sprintf(message, args...)}
}

func actionNotNeeded() *actionResult {
	return &actionResult{Needed: false}
}

func runAction(ctx *context, action taskAction) error {
	desc := action.description()

	result := action.needed(ctx)
	if result.Error != nil {
		return fmt.Errorf("The task action (%s) failed to detect whether it need to run: %s", desc, result.Error)
	}

	if result.Needed {
		if desc != "" {
			ctx.ui.TaskActionHeader(desc)
		}
		ctx.ui.Debug("Reason: \"%s\"", result.Reason)

		err := action.run(ctx)
		if err != nil {
			return fmt.Errorf("The task action failed to run: %s", err)
		}
	}

	result = action.needed(ctx)
	if result.Error != nil {
		return fmt.Errorf("The task action failed to detect if it is resolved: %s", result.Error)
	}

	if result.Needed {
		return fmt.Errorf("The task action did not produce the expected result: %s", result.Reason)
	}

	return nil
}
