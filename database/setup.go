package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func DBSet() *sql.DB {
	connStr := "postgresql://<username>:<password>@<database_ip>/todos?sslmode=disable"
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

func UserData(tableName string) {
	rows, err := DB.Query("SELECT * FROM " + tableName)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		// Scan the row into variables here
		// For example:
		// var id int
		// var name string
		// err = rows.Scan(&id, &name)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Println(id, name)
	}
}

func ProductData(tableName string) {
	rows, err := DB.Query("SELECT * FROM " + tableName)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		// Scan the row into variables here
		// For example:
		// var id int
		// var name string
		// var price float64
		// err = rows.Scan(&id, &name, &price)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Println(id, name, price)
	}
}
