package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/zeltbrennt/go-api/internal/server"
	"github.com/zeltbrennt/go-api/internal/store"
)

const port = 8090

// TaskServer godoc
//
//	@titlte Task server
//	@verison 1
//	@description Task Server API
//	@basePath /api/v1
//	@host localhost:8090
func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\b", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer) error {
	// interrupt
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	store := store.NewMockStore()

	logHandler := slog.NewJSONHandler(w, nil) // change this to a fileWriter to keep the logs
	logger := slog.New(logHandler)
	logger.Info(fmt.Sprintf("Server startet and listening on port %d", port))

	taskServer := server.NewTaskService(store, logger)
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: taskServer.Routes(),
	}
	// run server in goroutine
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()
	// warten auf stop signal oder server fehler
	select {
	case <-ctx.Done():
		logger.Info("shutting down server...")
	case err := <-srvErr:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server Shutdown failed %w", err)
	}

	logger.Info("server gracefully stopped")
	return nil
}
