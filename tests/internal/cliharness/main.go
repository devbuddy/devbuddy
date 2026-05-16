package cliharness

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/internal/budbuild"
)

var binaryPath string

func TestMain(m *testing.M) {
	start := time.Now()

	path, err := budbuild.TempBinaryPath("bud")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	binaryPath = path

	if err := budbuild.Build(budbuild.Options{
		OutputPath: binaryPath,
		LdFlags:    "-X main.Version=devel",
	}); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Built host binary in %s\n", time.Since(start))
	os.Exit(m.Run())
}
