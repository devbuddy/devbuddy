package integration

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/context"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags

	config = context.LoadConfig()
	fmt.Printf("Starting with config: %+v\n", config)

	fmt.Printf("Building linux binary\n")
	start := time.Now()
	cmd := exec.Command("go", "build", "-o", config.BinaryPath, "../cmd/bud")
	cmd.Env = append(os.Environ(), "GOOS=linux", "CGO_ENABLED=0")

	cmdOutput, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error when building binary: %s\n%s\n", err.Error(), string(cmdOutput))
		os.Exit(1)
	}
	fmt.Printf("Built in %s\n", time.Since(start))

	os.Exit(m.Run())
}
