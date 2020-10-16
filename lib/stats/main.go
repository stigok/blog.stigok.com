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

	// Setup database
	db = NewDatabase(dbname)

	// HTTP server
	srv := &http.Server{
		Handler:      VisitsRouter(db),
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}

	log.Printf("Listening on %s", srv.Addr)
	srv.ListenAndServe()
}
