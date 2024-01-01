package cmd

import (
	"fmt"
	"log"
	"os"

	color "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func GetFlagBool(cmd *cobra.Command, flag string) bool {
	b, err := cmd.Flags().GetBool(flag)
	if err != nil {
		log.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}
	return b
}

func onlyOneArg(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expecting one argument")
	}
	return nil
}

func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("expecting no arguments at all for command '%s'", cmd.Name())
	}
	return nil
}

func zeroOrOneArg(_ *cobra.Command, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("expecting zero or one argument")
	}
	return nil
}

func checkError(err error) {
	if err != nil {
		exitWithMessage(err.Error())
	}
}

func exitWithMessage(msg string) {
	fmt.Println(color.Red("Error:"), msg)
	os.Exit(1)
}
