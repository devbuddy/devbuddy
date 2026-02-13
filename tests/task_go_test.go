package integration

import (
	"testing"
	"time"

	"github.com/devbuddy/devbuddy/tests/context"
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
	c := CreateContextAndInit(t)

	t.Run("gomod_with_gopath", func(t *testing.T) {
		// Most common configuration, Modules and GOPATH set.

		c.Run(t, "export GOPATH=~")

		p := CreateProject(t, c,
			`up:`,
			`- go:`,
			`    version: "1.23.6"`,
		)
		c.Cd(t, p.Path)

		lines := c.Run(t, "bud up", context.Timeout(2*time.Minute))
		OutputContains(t, lines, "Golang (1.23.6)", "install golang distribution")
		OutputContains(t, lines, "activated: golang 1.23.6")

		lines = c.Run(t, "go version")
		OutputContains(t, lines, "go version go1.23.6")

		lines = c.Run(t, "echo GO111MODULE=${GO111MODULE}#")
		OutputContains(t, lines, "GO111MODULE=#")

		c.Write(t, "main.go", testingGoMain) // depends on github.com/devbuddy/tiny-test-go-package
		c.Run(t, "go mod init")
		c.Run(t, "go mod tidy", context.Timeout(15*time.Second))
		lines = c.Run(t, "go run main.go", context.Timeout(time.Minute))
		OutputContains(t, lines, "Is it working: true")
	})

	t.Run("gopath_absent", func(t *testing.T) {
		// Less common configuration, Modules but no GOPATH set.

		c.Run(t, "unset GOPATH")

		p := CreateProject(t, c,
			`up:`,
			`- go:`,
			`    version: "1.23.6"`,
			`    modules: true`,
		)
		lines := c.Cd(t, p.Path)
		OutputContains(t, lines, "activated: golang 1.23.6+mod")

		lines = c.Run(t, "bud up", context.Timeout(2*time.Minute))
		OutputContains(t, lines, "◼︎ Golang (1.23.6)")

		lines = c.Run(t, "go version")
		OutputContains(t, lines, "go version go1.23.6")

		lines = c.Run(t, "echo GO111MODULE=${GO111MODULE}#")
		OutputContains(t, lines, "GO111MODULE=on#")

		c.Write(t, "main.go", testingGoMain)
		c.Run(t, "go mod init project2")
		c.Run(t, "go mod tidy", context.Timeout(15*time.Second))
		lines = c.Run(t, "go run main.go", context.Timeout(time.Minute))
		OutputContains(t, lines, "Is it working: true")
	})

	t.Run("legacy_gopath_mode", func(t *testing.T) {
		// Old configuration: Modules disabled.
		// This is named "Legacy Go path mode" in https://go.dev/ref/mod#mod-commands
		// Note: GOPATH-mode `go get` was removed in Go 1.22, so we only verify
		// that bud configures the environment correctly.

		c.Run(t, "export GOPATH=~")

		p := CreateProject(t, c,
			`up:`,
			`- go:`,
			`    version: "1.23.6"`,
			`    modules: false`,
		)
		c.Cd(t, p.Path)

		lines := c.Run(t, "bud up", context.Timeout(2*time.Minute))
		OutputContains(t, lines, "◼︎ Golang (1.23.6)")

		lines = c.Run(t, "go version")
		OutputContains(t, lines, "go version go1.23.6")

		lines = c.Run(t, "echo GO111MODULE=${GO111MODULE}#")
		OutputContains(t, lines, "GO111MODULE=off#")
	})
}
