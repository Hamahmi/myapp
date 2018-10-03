package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbUser     = "root"
	dbPassword = "secret"
	dbName     = "myapp"
)

func main() {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   int
			name string
		)

		rows.Scan(&id, &name)

		log.Printf("%d: %s", id, name)
	}
}
