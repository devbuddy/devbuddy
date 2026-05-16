package integration

import (
	"fmt"
	"os"
	"testing"

	"github.com/devbuddy/devbuddy/tests/internal/harness"
)

func TestMain(m *testing.M) {
	if err := harness.SetupDocker(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}
