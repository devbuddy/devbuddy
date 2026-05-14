package cliharness

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

var binaryPath string

func TestMain(m *testing.M) {
	start := time.Now()
	repoRoot, err := findRepoRoot()
	if err != nil {
		fmt.Printf("Error finding repo root: %s\n", err)
		os.Exit(1)
	}

	binaryPath = filepath.Join(os.TempDir(), fmt.Sprintf("bud-test-%d", os.Getpid()), "bud")
	if err := os.MkdirAll(filepath.Dir(binaryPath), 0755); err != nil {
		fmt.Printf("Error creating binary directory: %s\n", err)
		os.Exit(1)
	}
	cmd := exec.Command("go", "build", "-ldflags", "-X main.Version=devel", "-o", binaryPath, "./cmd/bud")
	cmd.Dir = repoRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error when building binary: %s\n%s\n", err.Error(), string(output))
		os.Exit(1)
	}

	fmt.Printf("Built host binary in %s\n", time.Since(start))
	os.Exit(m.Run())
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
