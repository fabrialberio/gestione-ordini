package main

import (
	"fmt"
	"gestione-ordini/database"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	db := CreateDatabase()
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", logHandler(index))

	log.Println("Server started on port 8080.")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func CreateDatabase() *database.Database {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_CONTAINER_NAME"),
		os.Getenv("MYSQL_DATABASE"),
	)

	db, err := database.NewDatabase(dsn)
	for n_retries := 0; err != nil; n_retries++ {
		if n_retries == 5 {
			log.Fatalf("Error creating database after %v retries: %v", n_retries, err)
		}

		time.Sleep(5 * time.Second)
		log.Printf("Error creating database, retrying: %v", err)
		db, err = database.NewDatabase(dsn)
	}
	log.Println("Database created successfully.")

	return db
}
