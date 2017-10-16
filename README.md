### Live Player Map

![Online map example](https://cdn.discordapp.com/attachments/184738763589681153/369575727709945876/unknown.png)

This little program will provide a simple API to return all latitude/longitude locations of currently online players on a game realm (designed for use with WoW private servers). No IP/account references are sent from the API to the client, just coordinates.
Additionally supplied is an `index.html` file to display these locations on a map, given the user supplies his or her own [Mapbox](https://www.mapbox.com) API credentials.

To prevent abuse, it will check every minute to update its data, not on each `GET` request. There are no data stores, but instead caches the data in memory as it is assumed the program will be left running continuously. If used on a very large server with thousands of users needing to be queried frequently, a state service like MongoDB could easily be implemented, but it seems unnecessary. With the ability to query 150 IPs per minute, and with caching, it should keep up relatively well.

### Requirements
1. [Golang](https://golang.com) 1.x (built on 1.9)
2. An MySQL `account` table with online accounts (real or fake). TrinityCore/MaNGOS/OregonCore, etc. Literally could be used for anything with online account tables with IP addresses in them, not just WoW servers.

### Configuration
1. Rename `config.go.dist` to `config.go`
2. Edit newly renamed `config.go` file to set its MySQL credentials to the game server you're querying online accounts for.
3. (Optional) To use an example map, edit the `index.html` file and change the [Mapbox](https://www.mapbox.com) credentials to your own. Change `localhost:8080` to a new URL if you are running the Go program from a different location than localhost.

### Usage
1. Run `go get` in the Terminal/Command Prompt inside the program directory.
2. Run `go build && ./live-player-map` in the Terminal (`go build && ./live-player-map.exe` for Windows)
3. Run `index.html` in the browser with updated credentials to view data. 

### Notes
This program utilizes the free [ip-api](http://ip-api.com/) API to lookup IP address information. Due to its limitation of 150 
requests per minute, we have an internal counter to not exceed that threshold, and once hit, waits 1 minute before resuming.
