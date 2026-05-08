package executor

import (
	"os"
	"os/exec"

	"golang.org/x/term"
)

func runPassthrough(cmd *exec.Cmd) error {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		cmd.Stdin = os.Stdin
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
