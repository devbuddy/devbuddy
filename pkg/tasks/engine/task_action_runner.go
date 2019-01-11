package engine

import (
	"context"
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/tasks"
)

func runAction(ctx *context.Context, action tasks.TaskAction) error {
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

		err := action.run(ctx)
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
