package integration

import (
	"testing"

	"github.com/devbuddy/devbuddy/tests/context"
	"github.com/devbuddy/devbuddy/tests/internal/harness"
)

type Project = harness.Project

func TestMain(m *testing.M) {
	harness.TestMain(m)
}

func CreatePTYContext(t *testing.T) *context.TestContext {
	return harness.CreatePTYContext(t)
}

func CreatePTYContextAndInit(t *testing.T) *context.TestContext {
	return harness.CreatePTYContextAndInit(t)
}

func CreatePTYContextAndProject(t *testing.T, devYmlLines ...string) (*context.TestContext, Project) {
	return harness.CreatePTYContextAndProject(t, devYmlLines...)
}

func CreateProject(t *testing.T, c *context.TestContext, devYmlLines ...string) Project {
	return harness.CreateProject(t, c, devYmlLines...)
}

func OutputContains(t *testing.T, lines []string, subStrings ...string) {
	harness.OutputContains(t, lines, subStrings...)
}

func OutputNotContains(t *testing.T, lines []string, subStrings ...string) {
	harness.OutputNotContains(t, lines, subStrings...)
}

func OutputEqual(t *testing.T, lines []string, expectedLines ...string) {
	harness.OutputEqual(t, lines, expectedLines...)
}
