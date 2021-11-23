package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

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
	p.runGit(t, "init")
	p.runGit(t, "commit", "-m", "Commit1", "--allow-empty")
	p.runGit(t, "remote", "add", "origin", "git@github.com:org1/repo1.git")
}

func (p *projectWriter) runGit(t *testing.T, args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Dir = p.path
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"GIT_COMMITTER_NAME=John",
		"GIT_AUTHOR_NAME=John",
		"GIT_COMMITTER_EMAIL=john@doo.com",
		"GIT_AUTHOR_EMAIL=john@doo.com",
	}
	output, err := cmd.CombinedOutput()
	fmt.Printf("git output: %s\n", output)
	require.NoError(t, err)
}

type manifestWriter struct {
	path string
}

// GitInit initializes a simple Git repo
func (p *projectWriter) Manifest() *manifestWriter {
	return &manifestWriter{path: filepath.Join(p.path, "dev.yml")}
}

func (m *manifestWriter) Empty(t *testing.T) {
	WriteFile(m.path, []byte(""))
}

func (m *manifestWriter) WriteString(t *testing.T, value string) {
	WriteFile(m.path, []byte(value))
}
