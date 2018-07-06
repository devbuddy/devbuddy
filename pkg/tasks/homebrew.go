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

	for _, f := range formulas {
		task.addAction(&brewInstall{formula: f})
	}

	return nil
}

type brewInstall struct {
	formula string
	success bool
}

func (b *brewInstall) description() string {
	return fmt.Sprintf("installing %s", b.formula)
}

func (b *brewInstall) needed(ctx *context) (bool, error) {
	brew := helpers.NewHomebrew()

	installed := brew.IsInstalled(b.formula)

	return !installed, nil
}

func (b *brewInstall) run(ctx *context) error {
	err := command(ctx, "brew", "install", b.formula).Run()

	if err != nil {
		return fmt.Errorf("Homebrew failed: %s", err)
	}

	b.success = true
	return nil
}
