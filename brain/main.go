package main

import (
	"fmt"
	"flag"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/sczyh30/waffle-mesh/pkg/cmd"
	"github.com/sczyh30/waffle-mesh/brain/bootstrap"
)

var brainArgs bootstrap.BrainArgs

var command = &cobra.Command{
	Use: "waffle-brain",
	Short: "Waffle Brain works as the control plane center of the Waffle service mesh.",
	RunE: func(c *cobra.Command, args []string) error {
		stop := make(chan struct{})

		brainServer, err := bootstrap.NewServer(brainArgs)
		if err != nil {
			return fmt.Errorf("failed to create Waffle brain server: %v", err)
		}

		// Start the brain server.
		err = brainServer.Start(stop)
		if err != nil {
			return fmt.Errorf("failed to start Waffle Brain server: %v", err)
		}

		cmd.WaitSignal(stop)
		log.Println("Stopping the Waffle Brain server...")
		return nil
	},
}

func main() {
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	command.PersistentFlags().Uint32Var(&brainArgs.XdsProviderPort, "xdsPort",
		bootstrap.DefaultXdsProviderPort, "Port of discovery service server (gRPC)")
	command.PersistentFlags().Uint32Var(&brainArgs.MetricsServerPort, "metricsPort",
		bootstrap.DefaultMetricsServerPort, "Port of metrics server")

	flag.CommandLine.VisitAll(func(gf *flag.Flag) {
		command.PersistentFlags().AddGoFlag(gf)
	})
}


