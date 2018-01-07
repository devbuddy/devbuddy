package cmd

import (
	"fmt"

	color "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/integration"
	"github.com/pior/dad/pkg/project"
)

var cloneCmd = &cobra.Command{
	Use:   "clone [REMOTE]",
	Short: "Clone a project from github.com",
	Run:   cloneRun,
	Args:  OnlyOneArg,
}

func cloneRun(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	checkError(err)

	proj, err := project.NewFromID(args[0], cfg)
	checkError(err)

	if !proj.Exists() {
		err = proj.Clone()
		checkError(err)
	}

	fmt.Println(color.Brown("💡  Jumping to"), color.Green(proj.FullName()))
	err = integration.AddFinalizerCd(proj.Path)
	checkError(err)
}
