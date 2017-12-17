package integration

import (
	"fmt"
	"os"

	color "github.com/logrusorgru/aurora"
)

func Hook() {
	// notify("Yo! This is a message from Dad")
}

func notify(msg string) {
	fmt.Fprintf(os.Stderr, "👴🏽  %s\n", color.Cyan(msg))
}
