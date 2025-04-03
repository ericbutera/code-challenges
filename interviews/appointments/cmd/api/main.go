package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ericbutera/appointments/internal/api/server"
	"github.com/ericbutera/appointments/internal/data"
	"github.com/ericbutera/appointments/internal/db"
	"github.com/ericbutera/appointments/internal/logger"
	"github.com/ericbutera/appointments/internal/repo"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func main() {
	start()
}

func start() {
	slog := logger.New()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db := lo.Must(db.New(db.NewDefaultConfig()))
	seed(db)

	location := lo.Must(time.LoadLocation(repo.BusinessLocation))
	repo := lo.Must(repo.New(db, location))
	srv := lo.Must(server.New(repo, location))

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.Start()
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

func seed(db *gorm.DB) {
	// Note: use real migration tool in production
	mustCreate(db, &repo.User{Name: "User"})
	mustCreate(db, &repo.Trainer{Name: "Trainer"})

	appointments := lo.Must(data.GetFileJSONAs[[]*repo.Appointment]("appointments.json"))
	for _, a := range appointments {
		mustCreate(db, a)
	}
}

func mustCreate(db *gorm.DB, model any) {
	lo.Must0(db.Create(model).Error)
}
