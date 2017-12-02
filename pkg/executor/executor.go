package executor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	. "github.com/logrusorgru/aurora"
)

func Run(program string, args ...string) error {
	fmt.Println("Running: ", Bold(Cyan(program)), Cyan(strings.Join(args, " ")))

	cmd := exec.Command(program, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
