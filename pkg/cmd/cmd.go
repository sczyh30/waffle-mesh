package cmd

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitSignal(stop chan struct{}) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
	close(stop)
}
