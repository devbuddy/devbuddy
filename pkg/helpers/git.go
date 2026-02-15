package helpers

import (
	"fmt"
	"regexp"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
)

// GitRepo represents a local git repository
type GitRepo struct {
	ctx  *context.Context
	path string
}

// NewGitRepo returns a GitRepo
func NewGitRepo(ctx *context.Context, path string) *GitRepo {
	return &GitRepo{ctx: ctx, path: path}
}

// GetRemoteURL returns the URL of the origin remote
func (r *GitRepo) GetRemoteURL() (string, error) {
	cmd := executor.New("git", "config", "--get", "remote.origin.url")
	cmd.Cwd = r.path
	result := r.ctx.Executor.CaptureAndTrim(cmd)
	if result.Error != nil {
		return "", fmt.Errorf("failed to get the origin remote url: %w", result.Error)
	}
	return result.Output, nil
}

// GetCurrentBranch returns the name of the branch or "HEAD" for special cases
func (r *GitRepo) GetCurrentBranch() (string, error) {
	cmd := executor.New("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Cwd = r.path
	result := r.ctx.Executor.CaptureAndTrim(cmd)
	if result.Error != nil {
		return "", fmt.Errorf("failed to get the current branch: %w", result.Error)
	}
	return result.Output, nil
}

var githubRemotePatterns = []*regexp.Regexp{
	regexp.MustCompile(`^git@github\.com:([^/]+)/([^/.]+?)(?:\.git)?$`),
	regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/.]+?)(?:\.git)?$`),
}

func (r *GitRepo) buildGithubURL() (string, error) {
	remoteURL, err := r.GetRemoteURL()
	if err != nil {
		return "", err
	}
	for _, re := range githubRemotePatterns {
		matches := re.FindStringSubmatch(remoteURL)
		if matches != nil {
			return fmt.Sprintf("https://github.com/%s/%s", matches[1], matches[2]), nil
		}
	}
	return "", fmt.Errorf("unrecognized git remote url: %s", remoteURL)
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
