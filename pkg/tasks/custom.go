package tasks

import (
	"fmt"

	color "github.com/logrusorgru/aurora"

	"github.com/pior/dad/pkg/executor"
)

func init() {
	allTasks["custom"] = NewCustom
}

type Custom struct {
	condition string
	command   string
}

func NewCustom() Task {
	return &Custom{}
}

func (c *Custom) Load(definition interface{}) (bool, error) {
	def, ok := definition.(map[interface{}]interface{})
	if !ok {
		return false, nil
	}

	if payload, ok := def["custom"]; ok {
		properties := payload.(map[interface{}]interface{})

		command, ok := properties["meet"]
		if !ok {
			return false, nil
		}
		condition, ok := properties["met?"]
		if !ok {
			return false, nil
		}
		c.command = command.(string)
		c.condition = condition.(string)
		return true, nil
	}
	return false, nil
}

func (c *Custom) Perform() error {
	fmt.Printf("%s Custom: %s\n", color.Brown("★"), color.Cyan(c.command))

	code, err := executor.RunShellSilent(c.condition)
	if err != nil {
		fmt.Printf("Failed to run the condition command: %s", err)
		return taskFailed
	}
	if code == 0 {
		fmt.Println(color.Green("  Already good!"))
		return nil
	}

	// The condition command was run and returned a non-zero exit code.
	// It means we should run this custom task

	// fmt.Println(color.Brown("  Running"))
	code, err = executor.RunShellSilent(c.command)
	if err != nil {
		fmt.Printf("Command failed: %s", err)
		return taskFailed
	}
	if code != 0 {
		fmt.Println(color.Sprintf(color.Red("  Command exited with code %d"), code))
		return nil
	}
	fmt.Println(color.Green("  Done!"))
	return nil
}
