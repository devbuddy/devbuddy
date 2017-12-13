package cmd

import (
	"fmt"

	color "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/integration"
	"github.com/pior/dad/pkg/project"
)

var cdCmd = &cobra.Command{
	Use:   "cd [PROJECT]",
	Short: "Jump to a local project",
	Run:   cdRun,
	Args:  OnlyOneArg,
}

func cdRun(cmd *cobra.Command, args []string) {
	conf := config.Load()

	proj, err := project.FindBestMatch(args[0], conf)
	checkError(err)

	fmt.Println(color.Brown("💡  Jumping to"), color.Green(proj.FullName()))

	integration.AddFinalizerCd(proj.Path)
}
