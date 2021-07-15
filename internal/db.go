package internal

import (
	"sync"
)

const timeFormat string = "03:04 PM"

//go:generate mockgen -destination=./mocks/db.go -source=./db.go
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

type InMemoryDB struct {
	rwMtx sync.RWMutex
	store map[string]interface{}
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		store: make(map[string]interface{}),
	}
}

func (d *InMemoryDB) Set(key string, value interface{}) {
	d.rwMtx.Lock()
	defer d.rwMtx.Unlock()

	if len(key) < 1 {
		return
	}
	d.store[key] = value
}

func (d *InMemoryDB) Fetch(key string) interface{} {
	d.rwMtx.RLock()
	defer d.rwMtx.RUnlock()
	if value, found := d.store[key]; found {
		return value
	}

	//ignore: go nil interface
	return nil

}

func (d *InMemoryDB) Keys() []string {
	d.rwMtx.RLock()
	defer d.rwMtx.RUnlock()
	keys := make([]string, 0, len(d.store))
	for k := range d.store {
		keys = append(keys, k)
	}
	return keys
}
