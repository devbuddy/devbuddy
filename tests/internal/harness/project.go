package harness

import (
	"fmt"
	"math/rand/v2"
	"path/filepath"
	"testing"
)

// ProjectShell is the minimal context surface needed to create a test project:
// where to put it, how to write files, and how to enter it.
type ProjectShell interface {
	WriteLines(t *testing.T, path string, lines ...string)
	Cd(t *testing.T, path string) []string
	ProjectsDir() string
}

// NewProject creates a project directory with the given dev.yml content,
// cd's into it, and returns its absolute path.
func NewProject(t *testing.T, c ProjectShell, devYmlLines ...string) string {
	t.Helper()
	name := fmt.Sprintf("project-%x", rand.Int32())
	projectPath := filepath.Join(c.ProjectsDir(), "orgname", name)
	c.WriteLines(t, filepath.Join(projectPath, "dev.yml"), devYmlLines...)
	c.Cd(t, projectPath)
	return projectPath
}
