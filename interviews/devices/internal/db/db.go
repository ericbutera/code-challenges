package db

import (
	"errors"
	"time"

	"github.com/ericbutera/project/internal/models"
)

var ErrNotFound = errors.New("not found")

type StoreDeviceReadingsResult struct{}

type DeviceReadingsCount struct {
	Count int64
}

type DeviceLatestReading struct {
	Timestamp time.Time
}

// DB is an interface any persistent storage must implement.
// This allows swapping out in-memory storage for a real database (for instance using sqlite3 in-memory with gorm).
// Note: methods normally require ctx to handle timeout, cancellation & observability
type DB interface {
	StoreDeviceReadings(deviceID string, readings []*models.Reading) (*StoreDeviceReadingsResult, error)
	GetReadingCountByDevice(deviceID string) (*DeviceReadingsCount, error)
	GetLatestReadingByDevice(deviceID string) (*DeviceLatestReading, error)
	GetReadingsByDevice(deviceID string) ([]models.Reading, error)
}
