package main

import (
	"context"
	"errors"
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

type ServerConfig struct {
	store           store.Storer
	logger          *slog.Logger
	host            string
	port            int
	shutdownTimeout time.Duration
}

// TaskServer godoc
//
//	@titlte Task server
//	@verison 1
//	@description Task Server API
//	@basePath /api/v1
//	@host localhost:8090
func main() {
	ctx := context.Background()

	config := ServerConfig{
		host:            "localhost",
		port:            8090,
		shutdownTimeout: 5 * time.Second,
	}

	if err := setupAndRun(ctx, config, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func createServer(config ServerConfig) *http.Server {
	taskService := server.NewTaskService(config.store, config.logger)
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.host, config.port),
		Handler: taskService.Routes(),
	}
}

func setupAndRun(ctx context.Context, config ServerConfig, w io.Writer) error {
	// interrupt
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	store := store.NewMockStore()
	config.store = store
	logHandler := slog.NewJSONHandler(w, nil) // change this to a fileWriter to keep the logs
	logger := slog.New(logHandler)
	config.logger = logger

	srv := createServer(config)
	// run server in goroutine
	srvErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			srvErr <- err
		}
		close(srvErr)
	}()
	logger.Info(fmt.Sprintf("Server started on %s", srv.Addr))
	// warten auf stop signal oder server fehler

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-srvErr:
		return fmt.Errorf("server error: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, config.shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		if closErr := srv.Close(); closErr != nil {
			return errors.Join(err, closErr)
		}
	}

	logger.Info("server gracefully stopped")
	return nil
}
