package helpers

import (
	"github.com/pior/dad/pkg/executor"
)

// GitRepo represents a local git repository
type GitRepo struct {
	path string
}

// NewGitRepo returns a  GitRepo
func NewGitRepo(path string) *GitRepo {
	return &GitRepo{path: path}
}

// GetRemoteURL returns the URL of the origin remote
func (r *GitRepo) GetRemoteURL() (url string, err error) {
	url, err = executor.New("git", "config", "--get", "remote.origin.url").SetCwd(r.path).CaptureAndTrim()
	if err != nil {
		return "", err
	}
	return
}

// GetCurrentBranch returns the name of the branch or "HEAD" for special cases
func (r *GitRepo) GetCurrentBranch() (url string, err error) {
	url, err = executor.New("git", "rev-parse", "--abbrev-ref", "HEAD").SetCwd(r.path).CaptureAndTrim()
	if err != nil {
		return "", err
	}
	return
}
