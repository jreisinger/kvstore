package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jreisinger/kvstore/store"
)

func main() {
	r := mux.NewRouter()

	h, err := store.NewHandler()
	if err != nil {
		log.Fatal(err)
	}

	r.HandleFunc("/v1/{key}", h.Put).Methods("PUT")
	r.HandleFunc("/v1/{key}", h.Get).Methods("GET")
	r.HandleFunc("/v1/{key}", h.Del).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
}
