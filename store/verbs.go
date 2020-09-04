package store

import (
	"errors"
	"sync"
)

var store = struct {
	s sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func Put(key string, value string) error {
	store.s.Lock()
	store.m[key] = value
	store.s.Unlock()
	return nil
}

var ErrorNoSuchKey = errors.New("no such key")

func Get(key string) (string, error) {
	store.s.RLock()
	value, ok := store.m[key]
	store.s.RUnlock()
	if !ok {
		return "", ErrorNoSuchKey
	}
	return value, nil
}

func Delete(key string) error {
	delete(store.m, key)
	return nil
}
