package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/pior/dad/pkg/integration"
)

var rootCmd *cobra.Command

func build() {
	rootCmd = &cobra.Command{
		Use: "dad",
		Run: rootRun,
	}

	rootCmd.PersistentFlags().Bool("shell-init", false, "Shell initialization")
	// rootCmd.PersistentFlags().MarkHidden("shell-init")
	rootCmd.PersistentFlags().Bool("with-completion", false, "Enable completion during initialization")
	// rootCmd.PersistentFlags().MarkHidden("shell-completion")

	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(cdCmd)
	rootCmd.AddCommand(upCmd)
}

func buildTaskCommands() {

}

func rootRun(cmd *cobra.Command, args []string) {
	var err error

	if GetFlagBool(cmd, "shell-init") {
		if GetFlagBool(cmd, "with-completion") {
			err = rootCmd.GenBashCompletion(os.Stdout)
			checkError(err)
		}
		integration.Print()
		os.Exit(0)
	}

	err = cmd.Help()
	checkError(err)
}

func Execute() {
	build()
	buildTaskCommands()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
