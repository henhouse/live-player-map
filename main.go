package main

import (
	"log"
	"net/http"
)

func main() {
	// Establish connection to account database.
	db = &Database{Session: Connect()}
	defer db.Session.Close()

	// Initialize memory for our global IP address table for caching.
	IPTable = make(map[string]*IPLocation)

	// On startup, immediately begin collecting all online IP address information and caching it. This goroutine will
	// repeatedly run to update IP addresses every minute since a last successful update.
	go Updater()

	// Send incoming requests to our HTTP handler.
	http.HandleFunc("/", new(entryHandler).handleGET)

	port := "8080"
	log.Println("listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
