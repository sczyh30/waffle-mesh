package main

import (
	"fmt"
	"flag"
	"os"

	"github.com/spf13/cobra"
	"github.com/sczyh30/waffle-mesh/proxy/server"
	"github.com/sczyh30/waffle-mesh/pkg/cmd"
)

var proxyArgs server.ProxyArgs

var command = &cobra.Command{
	Use: "waffle-proxy",
	Short: "Waffle Proxy",
	RunE: func(c *cobra.Command, args []string) error {
		stop := make(chan struct{})

		proxyServer, err := server.NewProxy(proxyArgs)
		if err != nil {
			return fmt.Errorf("failed to create proxy: %v", err)
		}

		// Start the proxy server.
		err = proxyServer.StartProxy(stop)
		if err != nil {
			return fmt.Errorf("failed to start Waffle Proxy server: %v", err)
		}

		cmd.WaitSignal(stop)
		return nil
	},
}

func main() {
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	command.PersistentFlags().StringVar(&proxyArgs.BrainServerHost, "brainServerHost",
		server.DefaultBrainServerHost, "Host name of Waffle Brain Server")
	command.PersistentFlags().Uint32Var(&proxyArgs.GrpcPort, "xdsPort",
		server.DefaultGrpcPort, "Port of discovery service server (gRPC)")
	command.PersistentFlags().Uint32Var(&proxyArgs.MetricsPort, "metricsPort",
		server.DefaultMetricsPort, "Port of metrics server")
	command.PersistentFlags().Uint32Var(&proxyArgs.ListenerPort, "listenerPort",
		server.DefaultListenerPort, "Port of listener")

	flag.CommandLine.VisitAll(func(gf *flag.Flag) {
		command.PersistentFlags().AddGoFlag(gf)
	})
}

