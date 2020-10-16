package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Tiniest gif I could find in 3 seconds
// http://probablyprogramming.com/2009/03/15/the-tiniest-gif-ever
const gif_b64 string = "R0lGODlhAQABAIABAP///wAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=="

func VisitsRouter(db *Database) *mux.Router {
	r := mux.NewRouter()

	// Get visit
	get := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		post := vars["post"]

		cnt := db.GetVisitCount(post)
		log.Printf("get visit: %d %s", cnt, post)

		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, "%d", cnt)
	}

	// Record visit
	hit := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		post := vars["post"]
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(r.RemoteAddr+r.Header.Get("User-Agent"))))

		// Log visit
		log.Printf("visit: %s -> %s", hash, post)
		db.RecordVisit(post, hash)

		// Return a tiny gif
		w.Header().Add("Content-Type", "image/gif")
		gif, _ := base64.StdEncoding.DecodeString(gif_b64)
		w.Write(gif)
		return
	}

	r.Methods("GET").PathPrefix("/visits/{post}/get").HandlerFunc(get)
	r.Methods("GET").PathPrefix("/visits/{post}/hit").HandlerFunc(hit)

	return r
}
