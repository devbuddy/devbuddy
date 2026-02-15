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

// CreateGitRepo initializes a simple Git repo with an SSH remote
func (p *projectWriter) CreateGitRepo(t *testing.T) {
	t.Helper()
	p.CreateGitRepoWithRemote(t, "git@github.com:org1/repo1.git")
}

// CreateGitRepoWithRemote initializes a simple Git repo with the given remote URL
func (p *projectWriter) CreateGitRepoWithRemote(t *testing.T, remoteURL string) {
	t.Helper()

	p.runGit(t, "init")
	p.runGit(t, "commit", "-m", "Commit1", "--allow-empty")
	p.runGit(t, "remote", "add", "origin", remoteURL)
}

func (p *projectWriter) runGit(t *testing.T, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = p.path
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"GIT_COMMITTER_NAME=John",
		"GIT_AUTHOR_NAME=John",
		"GIT_COMMITTER_EMAIL=john@doo.com",
		"GIT_AUTHOR_EMAIL=john@doo.com",

		"GIT_CONFIG_COUNT=1",
		"GIT_CONFIG_KEY_0=init.defaultBranch",
		"GIT_CONFIG_VALUE_0=main",
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

func (m *manifestWriter) Empty() {
	WriteFile(m.path, []byte(""))
}

func (m *manifestWriter) WriteString(value string) {
	WriteFile(m.path, []byte(value))
}
