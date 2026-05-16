// Package budbuild compiles the `bud` binary for integration test harnesses.
//
// Both the docker harness (cross-compiled Linux binary mounted into a
// container) and the fast cliharness (host binary executed directly) need to
// build `cmd/bud` before their tests start. The variations are small enough
// to share: cross-compile env, optional ldflags, and output path.
package budbuild

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Options describe a single build invocation. OutputPath must be an absolute
// path to the binary file (its parent directory must exist).
type Options struct {
	OutputPath string
	Env        []string // appended to os.Environ()
	LdFlags    string   // -ldflags value (omitted if empty)
}

// Build compiles ./cmd/bud from the repo root with the given options.
func Build(opts Options) error {
	repoRoot, err := FindRepoRoot()
	if err != nil {
		return fmt.Errorf("finding repo root: %w", err)
	}

	args := []string{"build"}
	if opts.LdFlags != "" {
		args = append(args, "-ldflags", opts.LdFlags)
	}
	args = append(args, "-o", opts.OutputPath, "./cmd/bud")

	cmd := exec.Command("go", args...)
	cmd.Dir = repoRoot
	if len(opts.Env) > 0 {
		cmd.Env = append(os.Environ(), opts.Env...)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("building %s: %w\n%s", opts.OutputPath, err, output)
	}
	return nil
}

// FindRepoRoot walks up from the working directory until it finds a go.mod.
func FindRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", dir)
		}
		dir = parent
	}
}

// TempBinaryPath returns a unique absolute path under os.TempDir for a test
// binary, and ensures the parent directory exists. The path is process-scoped
// to avoid collisions with parallel test runs.
func TempBinaryPath(name string) (string, error) {
	path := filepath.Join(os.TempDir(), fmt.Sprintf("bud-test-%d", os.Getpid()), name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("creating binary directory: %w", err)
	}
	return path, nil
}
