package cmd

import "github.com/spf13/cobra"

var WaffleCommand = &cobra.Command {
	Use: "waffle",
	Short: "Waffle CLI manages the Waffle service mesh",
}
