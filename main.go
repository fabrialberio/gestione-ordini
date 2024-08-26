package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	time.Sleep(time.Second * 10)

	db, err := sql.Open("mysql", "user:password@tcp(mysql-container:3306)/testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to:", version)
}
