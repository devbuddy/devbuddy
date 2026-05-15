package project

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sahilm/fuzzy"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/utils"
	wt "github.com/devbuddy/devbuddy/pkg/worktree"
)

func FindBestMatch(expr string, cfg *config.Config) (found *Project, err error) {
	// 1. check if the expression is an ID
	p, err := NewFromID(expr, cfg)
	if err == nil && p.Exists() {
		return p, nil
	}

	// 2. fuzzy search on all projects
	projects, err := getAllProjects(cfg.SourceDir)
	if err != nil {
		return
	}

	if len(projects) == 0 {
		err = fmt.Errorf("no projects found at all! Try cloning one first")
		return
	}

	found = projectMatch(expr, projects)
	if found == nil {
		found = worktreeMatch(expr, projects)
	}
	if found == nil {
		err = fmt.Errorf("no project found for %s", expr)
	}
	return
}

func FindBestLinkMatch(expr string, index []string) string {
	matches := fuzzy.Find(expr, index)
	if matches.Len() >= 1 {
		return matches[0].Str
	}

	return ""
}

func projectMatch(expr string, projects []*Project) *Project {
	// First, try to match on project name only
	names := []string{}
	for _, p := range projects {
		names = append(names, p.Name())
	}
	matches := fuzzy.Find(expr, names)
	if matches.Len() >= 1 {
		return projects[matches[0].Index]
	}

	// Then, extend match to the organisation name as well
	names = []string{}
	for _, p := range projects {
		names = append(names, p.FullName())
	}
	matches = fuzzy.Find(expr, names)
	if matches.Len() >= 1 {
		return projects[matches[0].Index]
	}

	return nil
}

func worktreeMatch(expr string, projects []*Project) *Project {
	type candidate struct {
		project *Project
		branch  string
	}

	candidatesByPath := map[string]candidate{}
	for _, p := range projects {
		output, err := exec.Command("git", "-C", p.Path, "worktree", "list", "--porcelain").Output()
		if err != nil {
			continue
		}

		worktrees, err := wt.ParseListPorcelain(string(output))
		if err != nil {
			continue
		}

		for _, worktree := range worktrees {
			if worktree.Branch == "" {
				continue
			}

			c := candidate{
				project: projectFromWorktree(p, worktree.Path, worktree.Branch),
				branch:  worktree.Branch,
			}
			if worktree.Branch == expr {
				return c.project
			}
			candidatesByPath[worktree.Path] = c
		}
	}

	candidates := make([]candidate, 0, len(candidatesByPath))
	index := make([]string, 0, len(candidatesByPath))
	for _, c := range candidatesByPath {
		candidates = append(candidates, c)
		index = append(index, c.branch)
	}

	matches := fuzzy.Find(expr, index)
	if matches.Len() >= 1 {
		return candidates[matches[0].Index].project
	}

	return nil
}

func projectFromWorktree(base *Project, path string, branch string) *Project {
	hosting := *base.hosting
	hosting.repository = filepath.Base(path)
	return &Project{
		hosting:        &hosting,
		Path:           path,
		WorktreeBranch: branch,
	}
}

func getAllProjects(sourceDir string) ([]*Project, error) {
	var projects []*Project

	for _, platform := range getPlatformNames() {
		platformPath := filepath.Join(sourceDir, platform)
		if !utils.PathExists(platformPath) {
			continue
		}

		var orgPath string
		var projPath string

		orgs, err := listChildDir(platformPath)
		if err != nil {
			return nil, err
		}

		for _, org := range orgs {
			orgPath = filepath.Join(platformPath, org)

			repos, err := listChildDir(orgPath)
			if err != nil {
				return nil, err
			}

			for _, repo := range repos {
				projPath = filepath.Join(orgPath, repo)

				projects = append(projects, &Project{
					hosting: &hostingInfo{
						platform:     platform,
						organisation: org,
						repository:   repo,
					},
					Path: projPath,
				})
			}
		}
	}

	return projects, nil
}

func listChildDir(path string) (paths []string, err error) {
	files, err := os.ReadDir(path)
	if err != nil {
		err = fmt.Errorf("error listing files in %s: %w", path, err)
		return
	}

	for _, f := range files {
		if f.IsDir() {
			paths = append(paths, f.Name())
		}
	}
	return
}
