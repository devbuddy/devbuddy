package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/project"
)

var cloneCmd = &cobra.Command{
	Use:  "clone [REMOTE-PROJECT]",
	Run:  cloneRun,
	Args: cobra.ExactArgs(1),
}

func cloneRun(cmd *cobra.Command, args []string) {
	proj, err := project.NewFromIdentifier(args[0])
	if err != nil {
		log.Fatalln(err)
	}

	conf := config.Load()

	path, err := proj.Clone(conf)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("cd %s\n", path)
}
