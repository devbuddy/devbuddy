// Package harness provides shared helpers for integration test packages.
//
// Three concerns live here:
//   - building the `bud` binary (Linux for Docker-based tests, host for fast tests)
//   - Docker-based test contexts wrapping tests/context.TestContext
//   - host subprocess test contexts (no Docker, no PTY)
//
// Each test package owns its own TestMain and calls the appropriate Build*
// function during setup.
package harness

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// BuildLinuxBinary cross-compiles ./cmd/bud for linux/amd64 and returns the
// output path. The binary is intended to be mounted into a Docker container.
func BuildLinuxBinary() (string, error) {
	path, err := tempBinaryPath("bud-linux")
	if err != nil {
		return "", err
	}
	if err := buildBud(path, []string{"GOOS=linux", "CGO_ENABLED=0"}, ""); err != nil {
		return "", err
	}
	return path, nil
}

// BuildHostBinary compiles ./cmd/bud for the host platform with
// `-X main.Version=devel` and returns the output path.
func BuildHostBinary() (string, error) {
	path, err := tempBinaryPath("bud-host")
	if err != nil {
		return "", err
	}
	if err := buildBud(path, nil, "-X main.Version=devel"); err != nil {
		return "", err
	}
	return path, nil
}

func buildBud(outputPath string, env []string, ldflags string) error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return fmt.Errorf("finding repo root: %w", err)
	}

	args := []string{"build"}
	if ldflags != "" {
		args = append(args, "-ldflags", ldflags)
	}
	args = append(args, "-o", outputPath, "./cmd/bud")

	cmd := exec.Command("go", args...)
	cmd.Dir = repoRoot
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("building %s: %w\n%s", outputPath, err, output)
	}
	return nil
}

func findRepoRoot() (string, error) {
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

func tempBinaryPath(name string) (string, error) {
	path := filepath.Join(os.TempDir(), fmt.Sprintf("bud-test-%d", os.Getpid()), name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("creating binary directory: %w", err)
	}
	return path, nil
}
