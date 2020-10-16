package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var db *Database

func main() {
	// Setup database
	db = NewDatabase("mydb.sqlite.db")

	// Routes
	r := mux.NewRouter()
	r.PathPrefix("/visits/").Handler(VisitsRouter(db))

	// HTTP server
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}

	log.Printf("Listening on %s", srv.Addr)
	srv.ListenAndServe()
}
