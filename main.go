package main

import (
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

	templ *template.Template

	//go:embed public
	publicFS embed.FS
)

func main() {
	checkEnvVars()

	templ = template.Must(template.ParseGlob("templates/*.html"))
	template.Must(templ.ParseGlob("templates/**/*.html"))

	db = createDatabase()
	defer db.Close()

	addAdminUserIfNotExists()

	mux := http.NewServeMux()
	mux.Handle("/public/", http.FileServerFS(publicFS))
	mux.HandleFunc("/", logRequest(index))
	mux.HandleFunc("/login", logRequest(login))
	mux.HandleFunc("/logout", logRequest(logout))

	mux.HandleFunc("/cook", logRequest(index))
	mux.HandleFunc("/manager", logRequest(index))

	mux.HandleFunc("/admin", logRequest(admin))
	mux.HandleFunc("/admin/users/edit", logRequest(usersEdit))
	mux.HandleFunc("/admin/usersPage", logRequest(usersPage))
	mux.HandleFunc("/admin/usersTable", logRequest(usersTable))
	mux.HandleFunc("/admin/users/applyEdit", logRequest(usersApplyEdit))

	log.Println("Server started on port 8080.")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func logRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		h(w, r)
	}
}

func checkEnvVars() {
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

func createDatabase() *database.Database {
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

func addAdminUserIfNotExists() {
	_, err := db.GetUserByUsername("admin")
	if err == database.ErrRecordNotFound {
		hash, err := hashPassword(os.Getenv("ADMIN_PASSWORD"))
		if err != nil {
			log.Fatalf("Error hashing admin password: %v", err)
		}

		err = db.AddUser(database.User{
			RoleID:       database.RoleIDAdministrator,
			Username:     "admin",
			PasswordHash: hash,
			Name:         "Amministratore",
			Surname:      "",
		})
		if err != nil {
			log.Fatalf("Error creating admin user: %v", err)
		}
		log.Println("Admin user created successfully.")
	} else if err != nil {
		log.Fatalf("Error getting admin user: %v", err)
	}
}
