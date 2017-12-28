package main

import (
	"github.com/pior/dad/pkg/cmd"
)

var Version = "devel" // replaced by build flag

func main() {
	cmd.Execute(Version)
}
