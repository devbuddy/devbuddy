package tasks

import (
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	t := registerTaskDefinition("homebrew")
	t.name = "Homebrew"
	t.parser = parserHomebrew
}

func parserHomebrew(config *taskConfig, task *Task) error {
	var formulas []string

	for _, value := range config.payload.([]interface{}) {
		if v, ok := value.(string); ok {
			formulas = append(formulas, v)
		} else {
			return fmt.Errorf("invalid homebrew formulas")
		}
	}

	if len(formulas) == 0 {
		return fmt.Errorf("no homebrew formulas specified")
	}

	task.header = strings.Join(formulas, ", ")

	for _, formula := range formulas {
		action := newAction(fmt.Sprintf("installing %s", formula), func(ctx *context) error {
			result := command(ctx, "brew", "install", formula).Run()
			if result.Error != nil {
				return fmt.Errorf("failed to run brew install: %s", result.Error)
			}

			return nil
		})
		action.onFunc(func(ctx *context) *actionResult {
			brew := helpers.NewHomebrew()

			if brew.IsInstalled(formula) {
				return actionNotNeeded()
			}
			return actionNeeded("package %s is not installed", formula)
		})
		task.addAction(action)
	}

	return nil
}
