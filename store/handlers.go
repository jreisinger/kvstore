// Package store implements the key/value store functionality.
package store

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler represents an HTTP handler.
type Handler struct{}

// Put creates or updates an entry in the key/value store. It's idempotent.
func (h Handler) Put(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Get retrieves a value for a given key.
func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := Get(key)
	if errors.Is(err, ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil { // some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(value))
}
