package harness

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/context"
)

var config context.Config

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags

	var err error
	config, err = context.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Starting with config: %+v\n", config)

	fmt.Printf("Building linux binary\n")
	start := time.Now()
	repoRoot, err := findRepoRoot()
	if err != nil {
		fmt.Printf("Error finding repo root: %s\n", err)
		os.Exit(1)
	}
	cmd := exec.Command("go", "build", "-o", config.BinaryPath, "./cmd/bud")
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "GOOS=linux", "CGO_ENABLED=0")

	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error when building binary: %s\n%s\n", err.Error(), string(cmdOutput))
		os.Exit(1)
	}
	fmt.Printf("Built in %s\n", time.Since(start))

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
