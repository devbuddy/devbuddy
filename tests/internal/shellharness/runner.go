package shellharness

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Runner struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	lines  <-chan string
	done   chan error
	mu     sync.Mutex
	nextID int
}

type Result struct {
	Lines    []string
	ExitCode int
}

func Start(shellPath string, args ...string) (*Runner, error) {
	cmd := exec.Command(shellPath, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("opening stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("opening stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("opening stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting shell: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	return &Runner{
		cmd:   cmd,
		stdin: stdin,
		lines: scanLines(stdout, stderr),
		done:  done,
	}, nil
}

func (r *Runner) Run(ctx context.Context, command string) (Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Commands run in a long-lived shell, so we cannot rely on process exit to
	// know when one command is done. Append a unique sentinel after the user's
	// command and read until that sentinel appears. The sentinel also carries
	// "$?", which gives us the exit code without prompt parsing.
	r.nextID++
	sentinel := fmt.Sprintf("__BUD_SENTINEL_%d__", r.nextID)
	wrapped := fmt.Sprintf("%s\nprintf '%%s %%s\\n' %s \"$?\"\n", command, strconv.Quote(sentinel))
	if _, err := io.WriteString(r.stdin, wrapped); err != nil {
		return Result{}, fmt.Errorf("writing command: %w", err)
	}

	var output []string
	for {
		select {
		case err := <-r.done:
			return Result{}, fmt.Errorf("shell exited before sentinel: %w", err)
		case <-ctx.Done():
			return Result{}, ctx.Err()
		case line, ok := <-r.lines:
			if !ok {
				return Result{}, io.ErrUnexpectedEOF
			}
			exitCode, ok := parseSentinel(line, sentinel)
			if ok {
				return Result{Lines: output, ExitCode: exitCode}, nil
			}
			output = append(output, line)
		}
	}
}

func (r *Runner) Close() error {
	// Closing stdin asks the shell to exit. Once the shell exits, cmd.Wait
	// returns, stdout/stderr pipes close, and scanLines can drain and stop.
	_ = r.stdin.Close()
	select {
	case err := <-r.done:
		return err
	case <-time.After(time.Second):
		return r.cmd.Process.Kill()
	}
}

func scanLines(readers ...io.Reader) <-chan string {
	lines := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(readers))
	for _, reader := range readers {
		// Each pipe gets its own scanner so stdout and stderr cannot block each
		// other. The scanner loop stops when its pipe is closed, normally after
		// Runner.Close closes stdin and the shell process exits. The waiter below
		// closes the merged output channel only after every pipe scanner stops.
		go func() {
			defer wg.Done()
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				lines <- strings.TrimSuffix(scanner.Text(), "\r")
			}
		}()
	}
	go func() {
		wg.Wait()
		close(lines)
	}()
	return lines
}

func parseSentinel(line, sentinel string) (int, bool) {
	value, ok := strings.CutPrefix(line, sentinel+" ")
	if !ok {
		return 0, false
	}
	exitCode, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}
	return exitCode, true
}
