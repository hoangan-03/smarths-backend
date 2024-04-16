package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func DBSet() *sql.DB {
	connStr := "postgresql://homeowner:yevdpqFNRVXGXOfm9sGjPg@home-stay-6549.6xw.aws-ap-southeast-1.cockroachlabs.cloud:26257/homestay?sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println("Failed to connect to the database")
		return nil
	}
	fmt.Println("Successfully connected to the database")
	return db
}

var DB *sql.DB = DBSet()
