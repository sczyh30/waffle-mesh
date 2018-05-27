package cmd

import "github.com/spf13/cobra"

var getCommand = &cobra.Command {
	Use:   "get",
	Short: "Get Waffle resource from Kubernetes",
	Run: func(cmd *cobra.Command, args []string) {

	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	WaffleCommand.AddCommand(getCommand)
}
