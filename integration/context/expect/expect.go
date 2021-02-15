// Copyright 2016 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Updated for the usage of the DevBuddy project.

// Package expect implements a small expect-style interface
package expect

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
)

var (
	ErrCommandExited  = errors.New("command exited")
	ErrAlreadyClosed  = errors.New("already closed")
	ErrProcessStopped = errors.New("process stopped")
)

type ExpectProcess struct {
	closed atomic.Value

	cmd *exec.Cmd
	pty *os.File

	readerLines chan string
	readerError chan error

	// StopSignal is the signal Stop sends to the process; defaults to SIGKILL.
	StopSignal os.Signal

	// Debug enables some debug logs
	Debug bool
}

// NewExpect creates a new process for expect testing.
func NewExpect(name string, arg ...string) *ExpectProcess {
	// if env[] is nil, use current system env
	return NewExpectWithEnv(name, arg, nil)
}

// NewExpectWithEnv creates a new process with user defined env variables for expect testing.
func NewExpectWithEnv(name string, args []string, env []string) *ExpectProcess {
	cmd := exec.Command(name, args...)
	cmd.Env = env
	cmd.Stderr = cmd.Stdout
	cmd.Stdin = nil

	ep := &ExpectProcess{
		cmd:         cmd,
		StopSignal:  syscall.SIGKILL,
		readerLines: make(chan string, 100),
		readerError: make(chan error, 1),
	}
	ep.closed.Store(false)

	return ep
}

func (ep *ExpectProcess) Start() (err error) {
	if ep.pty, err = pty.Start(ep.cmd); err != nil { // start the process with a pty
		return err
	}

	// Set process in raw mode
	_, err = term.MakeRaw(int(ep.pty.Fd()))
	if err != nil {
		return err
	}

	go ep.reader()
	return nil
}

func (ep *ExpectProcess) reader() {
	r := bufio.NewReader(ep.pty)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			// by calling cmd.Wait() here, we can return interesting errors from expect functions.
			if cerr := ep.cmd.Wait(); cerr != nil {
				err = fmt.Errorf("process stopped: %w", cerr)
			} else {
				err = ErrProcessStopped
			}
			ep.readerError <- err
			close(ep.readerError)
			break
		}

		if line != "" {
			ep.debugLine(fmt.Sprintf("received %q", line))
			ep.readerLines <- line
		}
	}

	ep.debugLine("reader stopped")
}

func (ep *ExpectProcess) debugLine(line string) {
	if ep.Debug {
		fmt.Printf("%s[%d]: %v\n", ep.cmd.Path, ep.cmd.Process.Pid, line)
	}
}

// ExpectFunc returns the first line satisfying the function f.
func (ep *ExpectProcess) ExpectFunc(f func(string) bool) (string, error) {
	if ep.closed.Load().(bool) {
		return "", ErrAlreadyClosed
	}

	for {
		select {
		case err := <-ep.readerError:
			ep.debugLine(fmt.Sprintf("readerError=%v", err))
			return "", err
		case line := <-ep.readerLines:
			if f(line) {
				return line, nil
			}
		}
	}
}

// Expect returns the first line containing the given string.
func (ep *ExpectProcess) Expect(s string) (string, error) {
	return ep.ExpectFunc(func(txt string) bool { return strings.Contains(txt, s) })
}

// Line returns one line.
func (ep *ExpectProcess) Line() (string, error) {
	return ep.ExpectFunc(func(txt string) bool { return true })
}

// Stop kills the expect process and waits for it to exit.
func (ep *ExpectProcess) Stop() error {
	return ep.close(true)
}

// Signal sends a signal to the expect process
func (ep *ExpectProcess) Signal(sig os.Signal) error {
	return ep.cmd.Process.Signal(sig)
}

// Close waits for the expect process to exit.
func (ep *ExpectProcess) Close() error {
	return ep.close(false)
}

func (ep *ExpectProcess) close(kill bool) error {
	ep.debugLine("close()")

	if ep.closed.Load().(bool) {
		return ErrAlreadyClosed
	}

	defer func() {
		ep.closed.Store(true)
	}()

	if kill {
		ep.debugLine("kill()")
		_ = ep.Signal(ep.StopSignal)
	}

	err := <-ep.readerError // waiting for the reader to close after calling cmd.Wait()
	if err != nil {
		ep.debugLine(fmt.Sprintf("readerError=%v", err))
		if err == ErrProcessStopped {
			err = nil
		} else if !kill && strings.Contains(err.Error(), "exit status") { // non-zero exit code
			err = nil
		} else if kill && strings.Contains(err.Error(), "signal:") {
			err = nil
		}
	}

	ep.debugLine("pty.Close()")
	ep.pty.Close()

	return err
}

func (ep *ExpectProcess) Send(command string) error {
	ep.debugLine(fmt.Sprintf("sending %q", command))
	_, err := io.WriteString(ep.pty, command)
	return err
}
