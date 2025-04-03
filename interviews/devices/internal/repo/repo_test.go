package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/ericbutera/project/internal/db"
	"github.com/ericbutera/project/internal/models"
	"github.com/ericbutera/project/internal/repo"
	"github.com/ericbutera/project/internal/test"
	"github.com/stretchr/testify/require"
)

var testTime = time.Date(2025, 0o1, 13, 8, 0, 0, 0, time.UTC)

func TestStoreReadings(t *testing.T) {
	data := []*models.Reading{
		{Timestamp: testTime, Count: 2},
	}
	s := setupTest(t)

	s.db.EXPECT().
		StoreDeviceReadings(test.TestDeviceID, data).
		Return(&db.StoreDeviceReadingsResult{}, nil)

	_, err := s.repo.StoreReadings(context.Background(), test.TestDeviceID, data)
	require.NoError(t, err)
}

// TODO: GetCountByDevice
// TODO: GetLatestReadingByDevice

type testSetup struct {
	db   *db.MockDB
	repo *repo.Repo
}

func setupTest(t *testing.T) *testSetup {
	t.Helper()

	db := new(db.MockDB)

	r, err := repo.New(db)
	require.NoError(t, err)

	return &testSetup{
		db:   db,
		repo: r,
	}
}
