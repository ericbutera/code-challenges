package db_test

import (
	"testing"
	"time"

	"github.com/ericbutera/project/internal/db"
	"github.com/ericbutera/project/internal/models"
	"github.com/ericbutera/project/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryDB(t *testing.T) {
	t.Run("store readings out of order", func(t *testing.T) {
		data := []*models.Reading{
			{Timestamp: time.Date(2024, 12, 2, 0, 0, 0, 0, time.UTC), Count: 2},
			{Timestamp: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC), Count: 3},
			{Timestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), Count: 1},
		}
		db := setupReadings(t, data)

		readings, err := db.GetReadingsByDevice(test.TestDeviceID)
		require.NoError(t, err)

		require.Len(t, readings, len(data))
		assert.Equal(t, data[2].Timestamp, readings[0].Timestamp)
		assert.Equal(t, data[0].Timestamp, readings[1].Timestamp)
		assert.Equal(t, data[1].Timestamp, readings[2].Timestamp)
	})

	t.Run("store readings with duplicates", func(t *testing.T) {
		data := []*models.Reading{
			{Timestamp: time.Date(2024, 12, 2, 0, 0, 0, 0, time.UTC), Count: 2},
			{Timestamp: time.Date(2024, 12, 2, 0, 0, 0, 0, time.UTC), Count: 3},
			{Timestamp: time.Date(2024, 12, 2, 0, 0, 0, 0, time.UTC), Count: 1},
		}
		db := setupReadings(t, data)

		readings, err := db.GetReadingsByDevice(test.TestDeviceID)
		require.NoError(t, err)

		assert.Len(t, readings, 1)
		assert.Equal(t, data[0].Timestamp, readings[0].Timestamp)
	})

	t.Run("get reading count by device", func(t *testing.T) {
		db := setupReadings(t, []*models.Reading{
			{Timestamp: time.Date(2024, 12, 2, 0, 0, 0, 0, time.UTC), Count: 2},
			{Timestamp: time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC), Count: 3},
		})

		res, err := db.GetReadingCountByDevice(test.TestDeviceID)
		require.NoError(t, err)

		assert.Equal(t, int64(2), res.Count)
	})

	t.Run("get latest reading by device", func(t *testing.T) {
		data := []*models.Reading{
			{Timestamp: time.Date(2024, 12, 2, 0, 0, 0, 0, time.UTC), Count: 2},
			{Timestamp: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), Count: 1},
			{Timestamp: time.Date(2024, 12, 2, 0, 0, 0, 0, time.UTC), Count: 3},
		}
		db := setupReadings(t, data)

		latest, err := db.GetLatestReadingByDevice(test.TestDeviceID)
		require.NoError(t, err)

		assert.Equal(t, data[1].Timestamp, latest.Timestamp)
	})
}

func setupReadings(t *testing.T, data []*models.Reading) *db.InMemoryDB {
	t.Helper()
	db := db.NewInMemoryDB()
	_, err := db.StoreDeviceReadings(test.TestDeviceID, data)
	require.NoError(t, err)
	return db
}
