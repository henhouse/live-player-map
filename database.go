package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *Database

type Database struct {
	Session *sql.DB
}

func Connect() *sql.DB {
	var connStr string
	connStr += DB_USERNAME
	connStr += ":" + DB_PASSWORD
	connStr += "@tcp(" + DB_HOST + ":" + DB_PORT + ")"
	connStr += "/" + DB_DATABASE

	session, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
		log.Fatal("Couldn't connect to DB.")
	}

	// sql.Open above does not return connection failure, so we must ping to ensure the session is valid.
	err = session.Ping()
	if err != nil {
		log.Fatal(err)
		log.Fatal("Couldn't connect to DB.")
	}

	return session
}

// OnlineIPs is an array of strings containing the IP addresses of all online players. It is refreshed via the
// UpdateOnlineIPs function.
var OnlineIPs []string

// UpdateOnlineIPs queries the player account database for online accounts, returning their IP addresses. Each of these
// is appended to the OnlineIPs array which is reset just before.
func (db *Database) UpdateOnlineIPs() error {
	rows, err := db.Session.Query("SELECT last_ip FROM account WHERE online = 1")
	if err != nil {
		return err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("Failed to get columns", err)
		return err
	}

	// Result is your slice string.
	rawResult := make([][]byte, len(cols))
	result := []string{}

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			fmt.Println("Failed to scan row", err)
			return err
		}

		for _, raw := range rawResult {
			if raw != nil {
				//result[string(raw)] = new(IPLocation)
				result = append(result, string(raw))
			}
		}
	}

	// Clear out existing array for online IPs.
	OnlineIPs = nil

	// Assign newer online data to OnlineIPs array.
	OnlineIPs = result

	return nil
}
