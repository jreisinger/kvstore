// Package store implements the key/value store functionality.
package store

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jreisinger/kvstore/transactions"
)

// Handler represents an HTTP handler.
type Handler struct {
	transact transactions.TransactionLogger
}

// NewHandler initializes transactions logger and returns a handler.
func NewHandler() (Handler, error) {
	var err error

	transact, err := transactions.NewFileTransactionLogger("transaction.log")
	if err != nil {
		return Handler{}, fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := transact.ReadEvents()
	ok, e := true, transactions.Event{}

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case transactions.EventDelete:
				err = Delete(e.Key)
			case transactions.EventPut:
				err = Put(e.Key, e.Value)
			}
		}
	}

	transact.Run()

	return Handler{transact: transact}, err
}

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

	h.transact.WritePut(key, string(value))

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
