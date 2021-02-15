package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/devbuddy/devbuddy/tests/context"
)

func Test_Task_Go_Module(t *testing.T) {
	c := CreateContextAndInit(t)

	CreateProject(t, c, "project",
		`up:`,
		`- go:`,
		`    version: "1.15"`,
		`    modules: true`,
	)

	c.Run("export GOPATH=~")

	lines := c.Run("bud up", context.Timeout(2*time.Minute))
	OutputContains(t, lines, "Golang (1.15)", "install golang distribution")
	OutputContains(t, lines, "golang activated. (1.15+mod)")

	lines = c.Run("go version")
	OutputContains(t, lines, "go version go1.15")

	// Compile source with a dependency with a module

	c.Write("main.go", `
package main

import (
	"fmt"
	tiny "github.com/devbuddy/tiny-test-go-package"
)

func main() {
	fmt.Println("Is it working:", tiny.Working())
}
	`)
	c.Run("go mod init")
	lines = c.Run("go run main.go", context.Timeout(time.Minute))
	OutputContains(t, lines, "Is it working: true")
}

func Test_Task_Go_Pre_Module(t *testing.T) {
	c := CreateContextAndInit(t)

	_ = CreateProject(t, c, "project",
		`up:`,
		`- go: "1.15"`,
	)

	c.Run("export GOPATH=~")

	lines := c.Run("bud up", context.Timeout(2*time.Minute))
	OutputContains(t, lines, "Golang (1.15)", "install golang distribution")
	OutputContains(t, lines, "golang activated. (1.15)")

	lines = c.Run("go version")
	require.Len(t, lines, 1)
	require.Contains(t, lines[0], "go version go1.15")
}

func Test_Task_Go_Check_GOPATH(t *testing.T) {
	c := CreateContextAndInit(t)

	CreateProject(t, c, "project",
		`up:`,
		`- go: "1.15"`,
	)

	lines := c.Run("bud up", context.ExitCode(1))
	OutputContains(t, lines, "The GOPATH environment variable should be set to ~")
}
