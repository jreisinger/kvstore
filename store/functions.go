package store

import (
	"errors"
	"sync"
)

// ErrorNoSuchKey means the reuquired key is not in the store.
var ErrorNoSuchKey = errors.New("no such key")

var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func put(key string, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()
	return nil
}

func get(key string) (string, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()
	if !ok {
		return "", ErrorNoSuchKey
	}
	return value, nil
}

func del(key string) error {
	delete(store.m, key)
	return nil
}
