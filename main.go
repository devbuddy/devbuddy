package main

import (
	"github.com/devbuddy/devbuddy/pkg/cmd"
)

var Version = "devel" // replaced by build flag

func main() {
	cmd.Execute(Version)
}
