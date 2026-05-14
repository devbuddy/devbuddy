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

func CreateContext(t *testing.T) *context.TestContext {
	return harness.CreateContext(t)
}

func CreateContextAndInit(t *testing.T) *context.TestContext {
	return harness.CreateContextAndInit(t)
}

func CreateContextAndProject(t *testing.T, devYmlLines ...string) (*context.TestContext, Project) {
	return harness.CreateContextAndProject(t, devYmlLines...)
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
