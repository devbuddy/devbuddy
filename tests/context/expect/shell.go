package expect

import (
	"fmt"
	"strings"
)

type ShellExpect struct {
	process      *ExpectProcess
	promptString string

	Debug bool
}

func NewShellExpect(expect *ExpectProcess, prompt string) *ShellExpect {
	return &ShellExpect{
		process:      expect,
		promptString: prompt,
	}
}

// Init synchronize the prompt detection by expecting the initial prompt.
func (c *ShellExpect) Init() error {
	line, err := c.process.Line()
	if err != nil {
		return err
	}

	// First line should be the prompt, but can be prepended with terminal initialization chars.

	if strings.HasSuffix(norm(line), norm(c.promptString)) {
		c.debugLine("Found initial prompt\n")
		return nil
	}

	return fmt.Errorf("expected initial prompt, got %q", line)
}

func (c *ShellExpect) Run(command string) ([]string, error) {
	command = strings.TrimSuffix(command, "\n") + "\n"

	err := c.process.Send(command)
	if err != nil {
		return nil, fmt.Errorf("sending command %q: %w", command, err)
	}

	output, err := c.waitPrompt()
	if err != nil {
		return nil, fmt.Errorf("waiting prompt after command %q: %w", command, err)
	}

	return output, nil
}

func (c *ShellExpect) waitPrompt() ([]string, error) {
	var output []string

	c.debugLine("Waiting for the prompt")
	for {
		line, err := c.process.Line()
		if err != nil {
			return nil, fmt.Errorf("expecting output: %w", err)
		}

		if norm(line) == norm(c.promptString) {
			c.debugLine("Received prompt")
			return output, nil
		}

		c.debugLine(fmt.Sprintf("Received output: %q", line))
		output = append(output, trim(norm(line)))
	}
}

func (c *ShellExpect) debugLine(line string) {
	if c.Debug {
		fmt.Println(strings.TrimSuffix(line, "\n"))
	}
}

func norm(s string) string {
	return strings.Replace(s, "\r", "", -1)
}

func trim(s string) string {
	return strings.TrimSuffix(s, "\n")
}
