package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/tasks"
	"github.com/pior/dad/pkg/termui"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Ensure the project is up and running",
	Run:   upRun,
	Args:  noArgs,
}

func upRun(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	checkError(err)

	ui := termui.NewUI(cfg)

	proj, err := project.FindCurrent()
	checkError(err)

	success, err := tasks.RunAll(cfg, proj, ui)
	checkError(err)
	if !success {
		os.Exit(1)
	}
}
