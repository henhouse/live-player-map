package main

import (
	"encoding/json"
	"net/http"
)

type entryHandler struct{}

// handleGET processes all incoming HTTP requests to the base listening address. Only GET is a valid request method, and
// others are rejected. A complication of IP addresses cached in memory will be returned housing latitude and longitude
// locations for map rendering. These are internally updated at a defined time, and therefore spamming will not cause
// any flooding of the 3rd party IP-lookup API, but will always returned the cached data in memory.
func (*entryHandler) handleGET(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Collect list of cached IP address locations, building them into a suitable object for rendering on a map.
	locations := Return()

	// Assign headers.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Encode the location-object array into JSON, and write to response.
	json.NewEncoder(w).Encode(locations)
}
