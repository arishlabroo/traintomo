package internal

import (
	"sync"
)

const timeFormat string = "03:04 PM"

//go:generate mockgen -destination=./mocks/db.go -source=./db.go

// DB interface to store and retrieve schedule
type DB interface {
	Set(key string, value interface{})
	Fetch(key string) interface{}
	Keys() []string

	// Note: The interface is defined conforming to the requirements in the assignment.
	// More idiomatically, I would define this interface as below
	/*
		Set(key string, value interface{}) error
		Fetch(key string) (bool, interface{})
		Keys() []string
	*/
}

// InMemoryDB is a DB implementation which stores and retrieves data from memory
type InMemoryDB struct {
	rwMtx sync.RWMutex
	store map[string]interface{}
}

// NewInMemoryDB returns a new instance of the in memory db
func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		store: make(map[string]interface{}),
	}
}

// Set ...
func (d *InMemoryDB) Set(key string, value interface{}) {
	d.rwMtx.Lock()
	defer d.rwMtx.Unlock()

	if len(key) < 1 {
		return
	}
	d.store[key] = value
}

// Fetch ...
func (d *InMemoryDB) Fetch(key string) interface{} {
	d.rwMtx.RLock()
	defer d.rwMtx.RUnlock()
	if value, found := d.store[key]; found {
		return value
	}

	//ignore: go nil interface
	return nil

}

// Keys ...
func (d *InMemoryDB) Keys() []string {
	d.rwMtx.RLock()
	defer d.rwMtx.RUnlock()
	keys := make([]string, 0, len(d.store))
	for k := range d.store {
		keys = append(keys, k)
	}
	return keys
}
