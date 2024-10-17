package main

import (
	"embed"
	"fmt"
	"gestione-ordini/database"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	db *database.GormDB

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

	cookMux := http.NewServeMux()
	cookMux.HandleFunc("GET /", HandleGetCook)
	cookMux.HandleFunc("GET /ordersList", CheckPerm(database.PermIDEditOwnOrder, HandleGetCookOrdersList))
	cookMux.HandleFunc("GET /orders/{id}", HandleGetCookOrder)
	cookMux.HandleFunc("POST /orders", HandlePostCookOrder)

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("GET /", HandleGetAdmin)
	adminMux.HandleFunc("GET /usersTable", HandleGetAdminUsersTable)
	adminMux.HandleFunc("GET /users", HandleGetAdminUsers)
	adminMux.HandleFunc("GET /users/{id}", HandleGetAdminUser)
	adminMux.HandleFunc("POST /users", HandlePostAdminUser)

	managerMux := http.NewServeMux()
	managerMux.HandleFunc("GET /", HandleGetManager)
	managerMux.HandleFunc("GET /productsTable", HandleGetManagerProductsTable)
	managerMux.HandleFunc("GET /products/{id}", HandleGetManagerProduct)
	managerMux.HandleFunc("POST /products", HandlePostManagerProduct)

	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleGetIndex)
	mux.Handle("GET /public/", http.FileServerFS(publicFS))
	mux.Handle("/cook/", WithRole(database.RoleIDCook, http.StripPrefix("/cook", cookMux)))
	mux.Handle("/manager/", WithRole(database.RoleIDManager, http.StripPrefix("/manager", managerMux)))
	mux.Handle("/admin/", WithRole(database.RoleIDAdministrator, http.StripPrefix("/admin", adminMux)))

	mux.HandleFunc("POST /login", HandlePostLogin)
	mux.HandleFunc("POST /logout", HandlePostLogout)

	server := http.Server{
		Addr:    ":8080",
		Handler: WithLogging(mux),
	}

	log.Println("Server started on port 8080.")
	log.Fatal(server.ListenAndServe())
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

func createDatabase() *database.GormDB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_CONTAINER_NAME"),
		os.Getenv("MYSQL_DATABASE"),
	)

	db, err := database.New(dsn)
	for n_retries := 0; err != nil; n_retries++ {
		if n_retries == 5 {
			log.Fatalf("Error creating database after %v retries: %v", n_retries, err)
		}

		time.Sleep(5 * time.Second)
		log.Printf("Error creating database, retrying: %v", err)
		db, err = database.New(dsn)
	}
	log.Println("Database created successfully.")

	return db
}

func addAdminUserIfNotExists() {
	_, err := db.FindUserWithUsername("admin")
	if err == database.ErrRecordNotFound {
		hash, err := hashPassword(os.Getenv("ADMIN_PASSWORD"))
		if err != nil {
			log.Fatalf("Error hashing admin password: %v", err)
		}

		err = db.CreateUser(database.User{
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
