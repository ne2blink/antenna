package storage

import (
	"errors"
	"sync"
)

// NewFunc is a function to create a store
type NewFunc func(options map[string]interface{}) (Store, error)

var factories sync.Map

// Register registers a store
func Register(name string, newFunc NewFunc) {
	factories.Store(name, newFunc)
}

// Get gets a named store
func Get(name string, options map[string]interface{}) (Store, error) {
	v, found := factories.Load(name)
	if !found {
		return nil, errors.New(name + ": not found")
	}
	newFunc := v.(NewFunc)
	return newFunc(options)
}
