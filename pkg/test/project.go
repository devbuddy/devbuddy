package test

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Flaque/filet"
	"github.com/stretchr/testify/require"
)

type projectWriter struct {
	path string
}

// Project returns a test project writer
func Project(path string) *projectWriter {
	return &projectWriter{path}
}

// GitInit initializes a simple Git repo
func (p *projectWriter) CreateGitRepo(t *testing.T) {
	cmd := exec.Command("git", "init")
	cmd.Dir = p.path
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "commit", "-m", "Commit1", "--allow-empty")
	cmd.Dir = p.path
	cmd.Env = []string{
		"GIT_COMMITTER_NAME=John",
		"GIT_AUTHOR_NAME=John",
		"GIT_COMMITTER_EMAIL=john@doo.com",
		"GIT_AUTHOR_EMAIL=john@doo.com",
	}
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "remote", "add", "origin", "git@github.com:org1/repo1.git")
	cmd.Dir = p.path
	require.NoError(t, cmd.Run())
}

type manifestWriter struct {
	path string
}

// GitInit initializes a simple Git repo
func (p *projectWriter) Manifest() *manifestWriter {
	return &manifestWriter{path: filepath.Join(p.path, "dev.yml")}
}

func (m *manifestWriter) Empty(t *testing.T) {
	filet.File(t, m.path, "")
}

func (m *manifestWriter) WriteString(t *testing.T, value string) {
	filet.File(t, m.path, value)
}
