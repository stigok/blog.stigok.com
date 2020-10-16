package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var db *Database
var tinyGif []byte

func main() {
	// Tiniest gif I could find in 3 seconds
	// http://probablyprogramming.com/2009/03/15/the-tiniest-gif-ever
	gif, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIABAP///wAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==")
	tinyGif = gif

	// Setup database
	db = NewDatabase("mydb.sqlite.db")

	// Main router
	r := mux.NewRouter()
	r.Methods("GET").PathPrefix("/visits/{post}").HandlerFunc(visitsHandler)

	// HTTP server
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listening on %s", srv.Addr)
	srv.ListenAndServe()
}

func visitsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post := vars["post"]
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(r.RemoteAddr+r.Header.Get("User-Agent"))))

	// Log visit
	log.Printf("visit: %s -> %s", hash, post)
	db.RecordVisit(post, hash)

	// Return a tiny gif
	w.Header().Add("Content-Type", "image/gif")
	w.Write(tinyGif)
	return
}
