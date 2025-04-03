package dsa

import (
	"github.com/ericbutera/project/internal/models"
	"github.com/google/btree"
)

// TODO: make generic (only works with models.Reading)

// Stores readings in a B-tree
// features:
// - maintains sorting by timestamp as data is added
// - prevents duplicate readings
// - quick access to latest reading
type ReadingStore struct {
	tree *btree.BTree
}

func NewReadingStore() *ReadingStore {
	return &ReadingStore{tree: btree.New(2)} // B-tree of degree 2
}

// Add unique readings while ignoring duplicate timestamps
func (rs *ReadingStore) Add(reading models.Reading) {
	// note: models.Reading has a Less method which only compares timestamps (otherwise it would compare values too)
	if rs.tree.Has(reading) {
		return // prevent duplicates
	}
	rs.tree.ReplaceOrInsert(reading)
}

// Get returns all readings in descending order of timestampin descending order of timestamp
func (rs *ReadingStore) Get() []models.Reading {
	var readings []models.Reading
	rs.tree.Ascend(func(i btree.Item) bool {
		readings = append(readings, i.(models.Reading)) //nolint:forcetypeassert
		return true
	})
	return readings
}

// Returns only the latest reading
func (rs *ReadingStore) GetLatest() models.Reading {
	var reading models.Reading
	rs.tree.Ascend(func(i btree.Item) bool {
		reading = i.(models.Reading) //nolint:forcetypeassert
		return false
	})
	return reading
}

// Count returns the number of readings
// Assumption that discarded duplicates don't count towards "cumulative count"
func (rs *ReadingStore) Count() int {
	return rs.tree.Len()
}
