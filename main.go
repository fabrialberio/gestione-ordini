package main

import (
	"fmt"
	"gestione-ordini/database"
	"log"
	"os"
)

func main() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(mysql-container:3306)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_DATABASE"),
	)

	db, err := database.NewDatabase(dsn)
	if err != nil {
		log.Fatalf("Failed to create Database: %v", err)
	}
	defer db.Close()
}
