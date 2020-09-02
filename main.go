package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", PutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", GetHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
