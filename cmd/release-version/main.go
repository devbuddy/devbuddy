package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/devbuddy/devbuddy/internal/release"
)

func main() {
	kind := flag.String("kind", "", "release kind: minor, patch, rc, custom")
	custom := flag.String("version", "", "custom version when kind=custom")
	flag.Parse()

	tags, err := readLines(os.Stdin)
	if err != nil {
		fatal(err)
	}

	version, err := release.NextVersion(tags, release.Kind(*kind), *custom)
	if err != nil {
		fatal(err)
	}

	fmt.Println(version)
}

func readLines(file *os.File) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}
