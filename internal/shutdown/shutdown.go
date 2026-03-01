package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// GracefulShutdown waits for a SIGINT or SIGTERM signal and returns a context that is canceled when the signal is received.
func GracefulShutdown() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()
	return ctx
}
