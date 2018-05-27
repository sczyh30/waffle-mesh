package cmd

import "github.com/spf13/cobra"

var logCommand = &cobra.Command {
	Use:   "log",
	Short: "Get logs of Waffle Brain",
	Run: func(cmd *cobra.Command, args []string) {

	},
	Args: cobra.NoArgs,
}

func init() {
	WaffleCommand.AddCommand(logCommand)
}

