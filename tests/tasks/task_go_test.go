package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/internal/harness"
)

const testingGoMain = `
package main

import (
	"fmt"
	tiny "github.com/devbuddy/tiny-test-go-package"
)

func main() {
	fmt.Println("Is it working:", tiny.Working())
}
`

func Test_Task_Go(t *testing.T) {
	// Test with Go 1.23.6
	// Use sub-tests to avoid downloading the go runtime multiple times.
	c := harness.NewDockerPTYInit(t)

	t.Run("installs_and_runs_go_modules", func(t *testing.T) {
		harness.NewProject(t, c,
			`up:`,
			`- go:`,
			`    version: "1.23.6"`,
		)

		lines := c.Run(t, "bud up", harness.Timeout(2*time.Minute))
		harness.OutputContains(t, lines, "Golang (1.23.6)", "install golang distribution")
		harness.OutputContains(t, lines, "activated: golang[1.23.6]")

		lines = c.Run(t, "go version")
		harness.OutputContains(t, lines, "go version go1.23.6")

		c.Write(t, "main.go", testingGoMain)
		c.Run(t, "go mod init github.com/orgname/project", harness.Timeout(15*time.Second))
		c.Run(t, "go mod tidy", harness.Timeout(15*time.Second))
		lines = c.Run(t, "go run main.go", harness.Timeout(time.Minute))
		harness.OutputContains(t, lines, "Is it working: true")
	})

	t.Run("modules_false_is_rejected", func(t *testing.T) {
		harness.NewProject(t, c,
			`up:`,
			`- go:`,
			`    version: "1.23.6"`,
			`    modules: false`,
		)

		lines := c.Run(t, "bud up", harness.ExitCode(1), harness.Timeout(2*time.Minute))
		harness.OutputContains(t, lines, `task "go": "modules: false" is no longer supported for task "go"`)
	})
}
