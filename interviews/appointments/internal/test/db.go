package test

import (
	"testing"
	"time"

	"github.com/ericbutera/appointments/internal/data"
	"github.com/ericbutera/appointments/internal/db"
	"github.com/ericbutera/appointments/internal/repo"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const (
	TestTrainerID uint = 1
	TestUserID    uint = 1
)

type DBSetup struct {
	DB       *gorm.DB
	Repo     *repo.Repo
	Location *time.Location
}

func (s *DBSetup) Seed(t *testing.T) {
	t.Helper()

	user := repo.User{Name: "User"}
	MustCreate(t, s.DB, &user)

	trainer := repo.Trainer{Name: "Trainer"}
	MustCreate(t, s.DB, &trainer)

	appointments, err := data.GetFileJSONAs[[]*repo.Appointment]("appointments.json")
	require.NoError(t, err)
	for _, a := range appointments {
		MustCreate(t, s.DB, a)
	}
}

func MustCreate(t *testing.T, db *gorm.DB, model any) {
	t.Helper()
	require.NoError(t, db.Create(model).Error)
}

func NewSetup(t *testing.T) *DBSetup {
	t.Helper()

	location, err := time.LoadLocation(repo.BusinessLocation)
	require.NoError(t, err)

	db := NewDB(t)
	r, err := repo.New(db, location)
	require.NoError(t, err)

	db.Exec("DELETE FROM appointments") // TODO: fix - another lovely hack to clear data between tests. options are to use test suite or transaction per test

	return &DBSetup{
		DB:       db,
		Repo:     r,
		Location: location,
	}
}

func NewDB(t *testing.T) *gorm.DB {
	t.Helper()
	d, err := db.New(&db.Config{
		AutoMigrate: true,
		LogQueries:  false,
	})
	require.NoError(t, err)
	return d
}
