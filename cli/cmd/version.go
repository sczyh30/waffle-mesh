package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command {
	Use:   "version",
	Short: "Print the Waffle client and server version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.1")
	},
	Args: cobra.NoArgs,
}

func init() {
	WaffleCommand.AddCommand(versionCommand)
}