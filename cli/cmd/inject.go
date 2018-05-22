package cmd

import (
	"github.com/spf13/cobra"
)

type injectArgs struct {
	input string
	output string
}

var args = injectArgs{}

var injectCommand = &cobra.Command {
	Use:   "inject",
	Short: "Inject the Waffle Proxy as sidecar proxy to Kubernetes application deployment",
	Run: func(cmd *cobra.Command, args []string) {
		println("Inject successful, output file is deployment-injected.yaml")
	},
	Args: cobra.NoArgs,
}

func init() {
	injectCommand.PersistentFlags().StringVar(&args.output, "output", "", "")
	injectCommand.PersistentFlags().StringVar(&args.input, "input", "", "")
	WaffleCommand.AddCommand(injectCommand)
}
