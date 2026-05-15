package worktree

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
)

// Worktree is one record from `git worktree list --porcelain`.
type Worktree struct {
	Path     string
	Head     string
	Branch   string
	Detached bool
	Bare     bool
}

func ParseListPorcelain(output string) ([]Worktree, error) {
	var worktrees []Worktree
	var current *Worktree

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			if current != nil {
				worktrees = append(worktrees, *current)
				current = nil
			}
			continue
		}

		key, value, _ := strings.Cut(line, " ")
		switch key {
		case "worktree":
			if current != nil {
				worktrees = append(worktrees, *current)
			}
			if value == "" {
				return nil, fmt.Errorf("invalid git worktree record: missing path")
			}
			current = &Worktree{Path: value}
		case "HEAD":
			if current == nil {
				return nil, fmt.Errorf("invalid git worktree record: HEAD before worktree")
			}
			current.Head = value
		case "branch":
			if current == nil {
				return nil, fmt.Errorf("invalid git worktree record: branch before worktree")
			}
			current.Branch = strings.TrimPrefix(value, "refs/heads/")
		case "detached":
			if current == nil {
				return nil, fmt.Errorf("invalid git worktree record: detached before worktree")
			}
			current.Detached = true
		case "bare":
			if current == nil {
				return nil, fmt.Errorf("invalid git worktree record: bare before worktree")
			}
			current.Bare = true
		}
	}

	if current != nil {
		worktrees = append(worktrees, *current)
	}

	return worktrees, nil
}

func ManagedPath(repoPath, name string) (string, error) {
	slug := Slug(name)
	if slug == "" {
		return "", fmt.Errorf("worktree name must contain letters or numbers")
	}

	return filepath.Join(filepath.Dir(repoPath), filepath.Base(repoPath)+"--"+slug), nil
}

func Slug(name string) string {
	var b strings.Builder
	lastDash := false

	for _, r := range strings.TrimSpace(name) {
		keep := unicode.IsLetter(r) || unicode.IsNumber(r) || r == '.' || r == '_'
		if keep {
			b.WriteRune(r)
			lastDash = false
			continue
		}

		if !lastDash {
			b.WriteRune('-')
			lastDash = true
		}
	}

	return strings.Trim(b.String(), "-")
}

func CheckedOutBranch(worktrees []Worktree, branch string) *Worktree {
	branch = strings.TrimPrefix(branch, "refs/heads/")
	for i := range worktrees {
		if worktrees[i].Branch == branch {
			return &worktrees[i]
		}
	}
	return nil
}
