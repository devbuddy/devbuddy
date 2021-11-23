package helpers

import (
	"fmt"
	"regexp"

	"github.com/devbuddy/devbuddy/pkg/executor"
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
func (r *GitRepo) GetRemoteURL() (string, error) {
	result := executor.New("git", "config", "--get", "remote.origin.url").SetCwd(r.path).CaptureAndTrim()
	if result.Error != nil {
		return "", fmt.Errorf("failed to get the origin remote url: %w", result.Error)
	}
	return result.Output, nil
}

// GetCurrentBranch returns the name of the branch or "HEAD" for special cases
func (r *GitRepo) GetCurrentBranch() (string, error) {
	result := executor.New("git", "rev-parse", "--abbrev-ref", "HEAD").SetCwd(r.path).CaptureAndTrim()
	if result.Error != nil {
		return "", fmt.Errorf("failed to get the current branch: %w", result.Error)
	}
	return result.Output, nil
}

func (r *GitRepo) buildGithubURL() (string, error) {
	remoteURL, err := r.GetRemoteURL()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile("git@github.com:([^/]+)/([^.]+).git")
	matches := re.FindStringSubmatch(remoteURL)
	if matches == nil {
		return "", fmt.Errorf("unrecognized git remote url: %s", remoteURL)
	}
	url := fmt.Sprintf("https://github.com/%s/%s", matches[1], matches[2])
	return url, nil
}

// BuildGithubProjectURL builds the Github page url from the git remote url for a specific branch
func (r *GitRepo) BuildGithubProjectURL() (string, error) {
	baseURL, err := r.buildGithubURL()
	if err != nil {
		return "", err
	}
	branch, err := r.GetCurrentBranch()
	if err != nil {
		return "", err
	}
	return baseURL + "/tree/" + branch, nil
}

// BuildGithubPullrequestURL builds the Github pullrequest  url from the git remote url for a specific branch
func (r *GitRepo) BuildGithubPullrequestURL() (string, error) {
	baseURL, err := r.buildGithubURL()
	if err != nil {
		return "", err
	}
	branch, err := r.GetCurrentBranch()
	if err != nil {
		return "", err
	}
	return baseURL + "/pull/" + branch + "?expand=1", nil
}
