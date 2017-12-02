package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/integration"
	"github.com/pior/dad/pkg/project"
)

var cloneCmd = &cobra.Command{
	Use:  "clone [REMOTE]",
	Run:  cloneRun,
	Args: OnlyOneArg,
}

func cloneRun(cmd *cobra.Command, args []string) {
	proj, err := project.NewFromIdentifier(args[0])
	if err != nil {
		log.Fatalln(err)
	}

	conf := config.Load()

	proj.InferPath(conf)

	if !proj.Exists() {
		err := proj.Clone()
		if err != nil {
			log.Fatalln(err)
		}
	}

	integration.AddFinalizerCd(proj.Path)
}
