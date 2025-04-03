package db

import (
	"github.com/ericbutera/project/internal/dsa"
	"github.com/ericbutera/project/internal/models"
)

// In Memory store for devices and readings.
type InMemoryDB struct {
	readings dsa.Map[string, *dsa.ReadingStore] // deviceID -> btree
	// TODO: store devices in a map; devices  dsa.Map[UUID, *models.Device]
	// TODO: if we need all counts including duplicates; readingsCounters dsa.Map[UUID, int64] // { deviceID: readings-count }
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{}
}

func (db *InMemoryDB) StoreDeviceReadings(deviceID string, readings []*models.Reading) (*StoreDeviceReadingsResult, error) {
	// use separate maps for counts
	// latest reading can simply be the first reading in the sorted slice
	// note: readings can be out of order!
	store := db.getReadingsStore(deviceID)
	for _, reading := range readings {
		store.Add(*reading)
	}
	return &StoreDeviceReadingsResult{}, nil
}

func (db *InMemoryDB) GetReadingsByDevice(deviceID string) ([]models.Reading, error) {
	store := db.getReadingsStore(deviceID)
	return store.Get(), nil
}

// Note: this doesn't increase the count of duplicate readings
func (db *InMemoryDB) GetReadingCountByDevice(deviceID string) (*DeviceReadingsCount, error) {
	store := db.getReadingsStore(deviceID)
	return &DeviceReadingsCount{
		Count: int64(store.Count()),
	}, nil
}

func (db *InMemoryDB) GetLatestReadingByDevice(deviceID string) (*DeviceLatestReading, error) {
	store := db.getReadingsStore(deviceID)
	latest := store.GetLatest()
	if latest.Timestamp.IsZero() {
		return nil, ErrNotFound
	}
	return &DeviceLatestReading{
		Timestamp: latest.Timestamp,
	}, nil
}

// TODO: rename `getReadingsStore` to `loadOrCreate`
// TODO: getReadingsStore should return ErrNotFound if deviceID not found
func (db *InMemoryDB) getReadingsStore(deviceID string) *dsa.ReadingStore {
	var store *dsa.ReadingStore
	store, ok := db.readings.Load(deviceID)
	if !ok {
		store = dsa.NewReadingStore()
		db.readings.Store(deviceID, store)
	}
	return store
}
