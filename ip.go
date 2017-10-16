package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// IPLocation defines the JSON structure of an object returned by the ip-api.com API.
type IPLocation struct {
	Address      string `json:"query"`
	Status       string
	Country      string
	CountryCode  string
	Region       string
	RegionName   string
	PostalCode   string  `json:"zip"`
	Latitude     float32 `json:"lat"`
	Longitude    float32 `json:"lon"`
	ISP          string
	Organization string `json:"org"`
	//AsName       string `json:"as"`
}

// IPTable is a map containing IPLocation objects held by their IP address as a key. This object is internally filled
// in memory and serves as our storage for cached IP address entries. This method is opt-ed for over an internal storage
// option such as MongoDB simply to prevent further unnecessary dependencies.
var IPTable map[string]*IPLocation

// APIAddr is the location of the 3rd-party IP lookup API we will use.
const APIAddr = "http://ip-api.com/json"

// MaxLookupsPerMinute is the maximum amount of references we can put on the API per minute.
const MaxLookupsPerMinute = 150

// Location refers to a single user location in the world based on their IP address. This is the data returned to the
// client via our API to populate a live map.
type Location struct {
	Addr string
	Lat  float32
	Lon  float32
}

// lastUpdate holds a time value references the last time an update completed successfully.
var lastUpdate time.Time

// Update handles the main logic for updating our internal cache with IP address information. We retrieve a new list of
// online account IP addresses in the game server, and check to see if we have already looked up that IP address before.
// If we have, we skip over it and query the API for information on those we do not have, caching them as well.
func Update() {
	// Only allow Update to perform if the last completed update was more than 2 minutes ago.
	if time.Since(lastUpdate) < time.Minute*2 {
		return
	}

	// Retrieve a new list of currently-online account IP addresses from the game server.
	err := db.UpdateOnlineIPs()
	if err != nil {
		log.Println(err)
		return
	}

	// Ensure we returned at least one result.
	if len(OnlineIPs) < 1 {
		return
	}

	log.Println("Retrieved online account IPs, amount:", len(OnlineIPs))

	// ipLookupCount is our internal counter to ensure we do not look up too many IPs too quickly and get blacklisted
	// from sending too many requests. (Currently allowed 150 per minute)
	ipLookupCount := 0

	// Loop through each online IP.
	for _, addr := range OnlineIPs {
		// Check to see if we have this IP address cached already, if not, retrieve it via IP lookup API.
		if IPTable[addr] != nil {
			log.Printf("Already have %v in memory, skipping...\n", addr)
			continue
		} else {
			// Check to see if we have referenced near 150 IPs on the lookup API. If we have, to prevent us from getting
			// banned, we will set an internal wait for 1 minute to reset the flood guard.
			if ipLookupCount >= MaxLookupsPerMinute {
				waitTimer := time.NewTimer(time.Minute)

				// Blocks further processing until channel informs us the timer has expired.
				<-waitTimer.C

				// Reset counter back to 0
				ipLookupCount = 0
			}

			// Query 3rd-party IP lookup API for information on address, collect response.
			resp, err := http.Get(APIAddr + "/" + addr)
			if err != nil {
				log.Fatal("Unable to communicate with IP API.")
				continue
			}

			// Increment lookup counter.
			ipLookupCount++

			// Read in response body
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Failure to read packet data.")
				continue
			}
			resp.Body.Close()

			// Unmarshal response body into new IPLocation object.
			ipLoc := new(IPLocation)
			json.Unmarshal(b, ipLoc)

			log.Println("Caching IP:", ipLoc.Address)

			// Cache new IP address object to the global IPTable map.
			IPTable[addr] = ipLoc
		}
	}

	// Update the lastUpdate time to be now.
	lastUpdate = time.Now()
}

// Updater calls upon the Update function, and once it is complete, waits a few seconds, and invokes itself again to
// check if it is time to update again.
func Updater() {
	// Perform Update immediately (will be rejected if time is too soon)
	Update()

	// Check every 5 seconds if we need to update IPs.
	updateChecker := time.NewTimer(time.Second * 5)

	// Locks until timer above has passed.
	<-updateChecker.C

	// Call upon our self to update.
	Updater()
}

// Return compiles and returns an array of Location objects back containing the corresponding latitude, and longitude
// locations of addresses stored within a the IPTable map. Only this data is pertinent to rendering on a live map.
func Return() []Location {
	// Initial our return object.
	ret := []Location{}

	// Loop through all addresses stored in the IPTable map, appending a new Location to the return array with its data.
	for _, loc := range IPTable {
		ret = append(ret, Location{Addr: loc.Address, Lat: loc.Latitude, Lon: loc.Longitude})
	}

	return ret
}
