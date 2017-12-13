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

func OnlyOneArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expecting 1 argument")
	}
	return nil
}

func checkError(err error) {
	if err != nil {
		fmt.Println(color.Red("Error:"), err)
		os.Exit(-1)
	}
}
