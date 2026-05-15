package harness

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/context"
	"github.com/devbuddy/devbuddy/tests/internal/budbuild"
)

var config context.Config

func TestMain(m *testing.M) {
	var err error
	config, err = context.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err)
		os.Exit(1)
	}

	config.BinaryPath, err = budbuild.TempBinaryPath("bud")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting with config: %+v\n", config)

	fmt.Printf("Building linux binary\n")
	start := time.Now()
	if err := budbuild.Build(budbuild.Options{
		OutputPath: config.BinaryPath,
		Env:        []string{"GOOS=linux", "CGO_ENABLED=0"},
	}); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Built in %s\n", time.Since(start))

	os.Exit(m.Run())
}
