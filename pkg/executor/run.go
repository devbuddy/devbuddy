package executor

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
)

func runWithPTY(cmd *exec.Cmd, outputWriter io.Writer) error {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	defer func() {
		_ = ptmx.Close()
	}()

	// Handle pty size
	chResize := make(chan os.Signal, 1)
	signal.Notify(chResize, syscall.SIGWINCH)
	go func() {
		for range chResize {
			_ = pty.InheritSize(os.Stdin, ptmx)
		}
	}()
	chResize <- syscall.SIGWINCH // Initial resize.
	defer func() {
		signal.Stop(chResize)
		close(chResize)
	}()

	// Set Stdin in raw mode
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer func() {
		_ = terminal.Restore(int(os.Stdin.Fd()), oldState)
	}()

	go func() {
		_, _ = io.Copy(ptmx, os.Stdin)
	}()
	_, _ = io.Copy(outputWriter, ptmx)

	return cmd.Wait()
}
