package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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

	// Configure routes
	r := mux.NewRouter()
	r.Methods("GET").PathPrefix("/visits/{post}/get").HandlerFunc(GetVisits(db))
	r.Methods("GET").PathPrefix("/visits/{post}/hit").HandlerFunc(RecordVisit(db))
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// HTTP server
	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}

	log.Printf("Listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
