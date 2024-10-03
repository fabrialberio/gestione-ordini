package main

import (
	"database/sql"
	"embed"
	"fmt"
	"gestione-ordini/database"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

var (
	db *database.Database

	templates = template.Must(template.ParseGlob("templates/*.html"))

	//go:embed public
	publicFS embed.FS
)

func main() {
	CheckEnvVars()

	db = CreateDatabase()
	defer db.Close()

	AddAdminUserIfNotExists()

	mux := http.NewServeMux()
	mux.Handle("/public/", http.FileServerFS(publicFS))
	mux.HandleFunc("/", logRequest(index))
	mux.HandleFunc("/login", logRequest(login))
	mux.HandleFunc("/logout", logRequest(logout))
	mux.HandleFunc("/admin/users", logRequest(adminUsersTable))

	log.Println("Server started on port 8080.")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func CheckEnvVars() {
	envVars := []string{
		"MYSQL_USER",
		"MYSQL_PASSWORD",
		"MYSQL_CONTAINER_NAME",
		"MYSQL_DATABASE",
		"ADMIN_PASSWORD",
	}

	for _, envVar := range envVars {
		if _, ok := os.LookupEnv(envVar); !ok {
			log.Fatalf("Environment variable %s is not set.", envVar)
		}
	}
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

func AddAdminUserIfNotExists() {
	_, err := db.GetUserByUsername("admin")
	if err == sql.ErrNoRows {
		hash, err := hashPassword(os.Getenv("ADMIN_PASSWORD"))
		if err != nil {
			log.Fatalf("Error hashing admin password: %v", err)
		}

		err = db.AddUser(
			database.RoleIDAdministrator,
			"admin",
			hash,
			"Amministratore",
			"",
		)
		if err != nil {
			log.Fatalf("Error creating admin user: %v", err)
		}
		log.Println("Admin user created successfully.")
	} else if err != nil {
		log.Fatalf("Error getting admin user: %v", err)
	}
}
