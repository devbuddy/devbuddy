package tasks

import (
	"fmt"

	color "github.com/logrusorgru/aurora"

	"github.com/pior/dad/pkg/executor"
)

type Custom struct {
	condition string
	command   string
}

func (c *Custom) Load(definition map[interface{}]interface{}) (bool, error) {
	if payload, ok := definition["custom"]; ok {
		properties := payload.(map[interface{}]interface{})

		command, _ := properties["meet"]
		condition, _ := properties["met?"]
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
