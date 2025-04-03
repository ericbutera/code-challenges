package models

// TODO: use separate models for transport, domain logic and persistence layers
// TODO: reading data could be represented in a timeseries

import (
	"time"

	"github.com/google/btree"
)

type Reading struct {
	Timestamp time.Time `binding:"required"       description:"Device reported timestamp of reading" example:"2021-01-01T00:00:00-05:00" json:"timestamp"`
	Count     int       `binding:"required,min=0" description:"Reading data"                         example:"17"                        json:"count"`
	// TODO: compare device reported Timestamp to server CreatedAt to determine time integrity; DDIA: unreliable clocks
}

// used by dsa.ReadingStore to ensure unique readings
func (r Reading) Less(than btree.Item) bool {
	return r.Timestamp.After(than.(Reading).Timestamp) //nolint:forcetypeassert
}

// TODO: in case we want to save device with extra data
// type Device struct {
// 	ID string `binding:"required,uuid4" description:"Device ID" example:"00000000-0000-0000-0000-000000000000" json:"id"`
// }
