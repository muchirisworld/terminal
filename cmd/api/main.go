package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/muchirisworld/terminal/internal/app"
	"github.com/muchirisworld/terminal/internal/shutdown"
)

func main() {
	ctx := shutdown.GracefulShutdown()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	a, err := app.New(ctx)
	if err != nil {
		return err
	}

	go func() {
		if err := a.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.Logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), a.Config.ShutdownTimeout)
	defer shutdownCancel()

	return a.Close(shutdownCtx)
}
