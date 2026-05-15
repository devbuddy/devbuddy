package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/internal/cliharness"
)

type Context = cliharness.Context
type Project = cliharness.Project
type RunOption = cliharness.RunOption

func TestMain(m *testing.M) {
	cliharness.TestMain(m)
}

func CreateContext(t *testing.T) *cliharness.Context {
	return cliharness.CreateContext(t)
}

func CreateContextAndInit(t *testing.T) *cliharness.Context {
	return cliharness.CreateContextAndInit(t)
}

func CreateContextAndProject(t *testing.T, devYmlLines ...string) (*cliharness.Context, Project) {
	return cliharness.CreateContextAndProject(t, devYmlLines...)
}

func CreateProject(t *testing.T, c *cliharness.Context, devYmlLines ...string) Project {
	return cliharness.CreateProject(t, c, devYmlLines...)
}

func ExitCode(code int) RunOption {
	return cliharness.ExitCode(code)
}

func Timeout(timeout time.Duration) RunOption {
	return cliharness.Timeout(timeout)
}

func OutputContains(t *testing.T, lines []string, subStrings ...string) {
	cliharness.OutputContains(t, lines, subStrings...)
}

func OutputNotContains(t *testing.T, lines []string, subStrings ...string) {
	cliharness.OutputNotContains(t, lines, subStrings...)
}

func OutputEqual(t *testing.T, lines []string, expectedLines ...string) {
	cliharness.OutputEqual(t, lines, expectedLines...)
}
