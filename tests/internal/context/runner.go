package context

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

// shellRunnerImpl drives a long-lived non-PTY shell over stdin/stdout pipes.
// Commands cannot rely on process exit to delimit output, so each command is
// followed by a unique sentinel that also carries the exit code.
type shellRunnerImpl struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	lines  <-chan string
	done   chan error
	mu     sync.Mutex
	nextID int
}

type runnerResult struct {
	Lines    []string
	ExitCode int
}

func startShellRunner(shellPath string, args ...string) (*shellRunnerImpl, error) {
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

	return &shellRunnerImpl{
		cmd:   cmd,
		stdin: stdin,
		lines: scanLines(stdout, stderr),
		done:  done,
	}, nil
}

func (r *shellRunnerImpl) RunWithExitCode(command string, timeout time.Duration) ([]string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := r.run(ctx, command)
	if err != nil {
		return nil, 0, err
	}
	return result.Lines, result.ExitCode, nil
}

func (r *shellRunnerImpl) run(ctx context.Context, command string) (runnerResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	sentinel := fmt.Sprintf("__BUD_SENTINEL_%d__", r.nextID)
	wrapped := fmt.Sprintf("%s\nprintf '%%s %%s\\n' %s \"$?\"\n", command, strconv.Quote(sentinel))
	if _, err := io.WriteString(r.stdin, wrapped); err != nil {
		return runnerResult{}, fmt.Errorf("writing command: %w", err)
	}

	var output []string
	for {
		select {
		case err := <-r.done:
			return runnerResult{}, fmt.Errorf("shell exited before sentinel: %w", err)
		case <-ctx.Done():
			return runnerResult{}, ctx.Err()
		case line, ok := <-r.lines:
			if !ok {
				return runnerResult{}, io.ErrUnexpectedEOF
			}
			exitCode, ok := parseSentinel(line, sentinel)
			if ok {
				return runnerResult{Lines: output, ExitCode: exitCode}, nil
			}
			output = append(output, line)
		}
	}
}

func (r *shellRunnerImpl) Close() error {
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
		go func() {
			defer wg.Done()
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				lines <- strings.TrimSuffix(scanner.Text(), "\r")
			}
			_ = scanner.Err()
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
