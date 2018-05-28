package cmd

import (
	"github.com/spf13/cobra"
	"os/exec"
	"fmt"
)

var getCommand = &cobra.Command {
	Use:   "get",
	Short: "Get Waffle resource from Kubernetes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("please specify the resource kind")
		}
		c := exec.Command("kubectl", "get", args[0])
		output, err := c.Output()
		fmt.Printf("%s\n", string(output))
		return err
	},
}

func init() {
	WaffleCommand.AddCommand(getCommand)
}
