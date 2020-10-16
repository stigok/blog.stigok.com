package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

var db *Database

func main() {
	dbname := os.Getenv("DATABASE_NAME")
	if dbname == "" {
		dbname = MEMORY_DB
	}
	log.Printf("using database %s", dbname)

	addr := os.Getenv("LISTEN_ADDRESS")
	if addr == "" {
		addr = "0.0.0.0:8000"
	}

	// Setup database
	db = NewDatabase(dbname)

	// HTTP server
	srv := &http.Server{
		Handler:      VisitsRouter(db),
		Addr:         addr,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}

	log.Printf("Listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
