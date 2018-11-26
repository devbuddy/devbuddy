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
		// task.addAction(&brewInstall{formula: f})
		task.
			addActionWithBuilder(fmt.Sprintf("installing %s", formula), func(ctx *context) error {
				result := command(ctx, "brew", "install", formula).Run()
				if result.Error != nil {
					return fmt.Errorf("failed to run brew install: %s", result.Error)
				}

				return nil
			}).
			addConditionFunc(func(ctx *context) *actionResult {
				brew := helpers.NewHomebrew()

				if brew.IsInstalled(formula) {
					return actionNotNeeded()
				}
				return actionNeeded("package %s is not installed", formula)
			})
	}

	return nil
}

// type brewInstall struct {
// 	formula string
// }

// func (b *brewInstall) description() string {
// 	return fmt.Sprintf("installing %s", b.formula)
// }

// func (b *brewInstall) needed(ctx *context) *actionResult {
// 	brew := helpers.NewHomebrew()

// 	if brew.IsInstalled(b.formula) {
// 		return actionNotNeeded()
// 	}
// 	return actionNeeded("package %s is not installed", b.formula)
// }

// func (b *brewInstall) run(ctx *context) error {
// 	result := command(ctx, "brew", "install", b.formula).Run()
// 	if result.Error != nil {
// 		return fmt.Errorf("failed to run brew install: %s", result.Error)
// 	}

// 	return nil
// }
