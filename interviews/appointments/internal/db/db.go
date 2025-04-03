package db

import (
	"time"

	"github.com/ericbutera/appointments/internal/repo"
	gorm_logrus "github.com/onrik/gorm-logrus"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	models = []any{
		repo.Appointment{},
		repo.User{},
		repo.Trainer{},
	}
)

type Config struct {
	LogQueries  bool
	AutoMigrate bool
}

func NewDefaultConfig() *Config {
	return &Config{
		LogQueries:  true,
		AutoMigrate: true,
	}
}

func New(config *Config) (*gorm.DB, error) {
	// TODO: options pattern
	// TODO: add support for multiple adapters (postgres)
	// TODO: use a real migration tool like golang-migrate

	opts := &gorm.Config{}
	if config.LogQueries { // TODO: convert to slog
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetLevel(logrus.DebugLevel)
		opts.Logger = gorm_logrus.New()
	}

	d, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		opts,
	)
	if err != nil {
		return nil, err
	}

	d.NowFunc = func() time.Time {
		return time.Now().UTC()
	}

	if config.AutoMigrate {
		if err := d.AutoMigrate(models...); err != nil {
			return nil, err
		}
	}

	return d, nil
}
