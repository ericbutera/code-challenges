package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ericbutera/project/internal/api"
	"github.com/ericbutera/project/internal/db"
	"github.com/ericbutera/project/internal/logger"
	"github.com/ericbutera/project/internal/repo"
	"github.com/samber/lo"
)

func main() {
	start()
}

func start() {
	slog := logger.New()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db := db.NewInMemoryDB()

	repo := lo.Must(repo.New(db))
	server := lo.Must(api.New(repo))

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- server.Start()
	}()

	select {
	case err := <-srvErr:
		quit(ctx, err)
	case <-ctx.Done():
		slog.Info("shutting down")
		stop()
	}
}

func quit(ctx context.Context, err error) {
	slog.ErrorContext(ctx, err.Error())
	os.Exit(1)
}
