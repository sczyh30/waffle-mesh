package main

import (
	"os"

	"github.com/sczyh30/waffle-mesh/cli/cmd"
)

func main() {
	if err := cmd.WaffleCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
