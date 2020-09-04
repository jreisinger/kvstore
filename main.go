package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jreisinger/kvstore/store"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", store.PutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", store.GetHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
