package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

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

		t := r.Header.Get("Accept")

		// Return JSON
		if strings.Contains(t, "application/json") {
			obj, _ := json.Marshal(struct {
				Success bool   `json:"success"`
				Error   string `json:"error"`
				Post    string `json:"post"`
				Count   int    `json:"count"`
			}{
				true,
				"",
				post,
				cnt,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(obj)
			return
		}

		// Return plaintext
		//if strings.Contains(t, "text/plain") {
		w.Write([]byte("ok\n"))
		//	return
		//}

		// Return gif to use in <img> tags
		// fmt.Fprintf(w, "%d", cnt)
		// w.Header().Add("Content-Type", "image/gif")
		// gif, _ := base64.StdEncoding.DecodeString(gif_b64)
		// w.Write(gif)
	}

	// Record visit
	hit := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		post := vars["post"]

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = ""
		}

		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(ip+r.Header.Get("User-Agent"))))

		// Log visit
		log.Printf("visit: %s -> %s", hash, post)
		db.RecordVisit(post, hash)

		t := r.Header.Get("Accept")

		// Return JSON
		if strings.Contains(t, "application/json") {
			obj, _ := json.Marshal(struct {
				Success bool   `json:"success"`
				Error   string `json:"error"`
			}{
				true,
				"",
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(obj)
			return
		}

		// Return plaintext
		if strings.Contains(t, "text/plain") {
			w.Write([]byte("ok\n"))
			return
		}

		// Return gif to use in <img> tags
		w.Header().Add("Content-Type", "image/gif")
		gif, _ := base64.StdEncoding.DecodeString(gif_b64)
		w.Write(gif)
	}

	r.Methods("GET").PathPrefix("/visits/{post}/get").HandlerFunc(get)
	r.Methods("GET").PathPrefix("/visits/{post}/hit").HandlerFunc(hit)

	return r
}
