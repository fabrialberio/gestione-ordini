package main

import (
	"embed"
	"fmt"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/database"
	"gestione-ordini/pkg/handlers"
	mw "gestione-ordini/pkg/middleware"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed public
var publicFS embed.FS

func main() {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	template.Must(tmpl.ParseGlob("templates/**/*.html"))

	db := createDatabase()
	defer db.Close()
	addAdminUserIfNotExists(db)

	mux := http.NewServeMux()
	setupRoutes(mux)

	server := http.Server{
		Addr:    ":8080",
		Handler: mw.WithContext(db, tmpl, mux),
	}

	log.Println("Server started on port 8080.")
	log.Fatal(server.ListenAndServe())
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

func addAdminUserIfNotExists(db *database.GormDB) {
	_, err := db.FindUserWithUsername("admin")
	if err == database.ErrRecordNotFound {
		hash, err := auth.HashPassword(os.Getenv("ADMIN_PASSWORD"))
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

func setupRoutes(mux *http.ServeMux) {
	chefMux := http.NewServeMux()
	chefMux.HandleFunc("GET "+handlers.DestChef, handlers.GetChef)
	chefMux.HandleFunc("GET "+handlers.DestChefOrders+"{id}", handlers.GetChefOrder)
	chefMux.HandleFunc("POST "+handlers.DestChefOrders, handlers.PostChefOrder)

	consoleMux := http.NewServeMux()
	consoleMux.HandleFunc("GET "+handlers.DestConsole, handlers.GetConsole)
	consoleMux.HandleFunc("GET "+handlers.DestNewOrder, handlers.GetNewOrder)
	consoleMux.HandleFunc("GET "+handlers.DestOrders+"{id}", handlers.GetConsoleOrder)
	consoleMux.HandleFunc("POST "+handlers.DestOrders, handlers.PostConsoleOrder)

	consoleMux.HandleFunc("GET "+handlers.DestAllOrders, handlers.GetAllOrders)
	consoleMux.HandleFunc("GET "+handlers.DestAllOrdersView, handlers.GetAllOrdersView)
	consoleMux.HandleFunc("POST "+handlers.DestOrderSelection, handlers.PostOrderSelection)
	consoleMux.HandleFunc("POST "+handlers.DestOrderSelectionCount, handlers.PostOrderSelectionCount)

	consoleMux.HandleFunc("GET "+handlers.DestProducts, handlers.GetProducts)
	consoleMux.HandleFunc("GET "+handlers.DestProducts+"{id}", handlers.GetProduct)
	consoleMux.HandleFunc("POST "+handlers.DestProducts, handlers.PostProduct)
	consoleMux.HandleFunc("GET "+handlers.DestProductsTable, handlers.GetProductsTable)
	consoleMux.HandleFunc("POST "+handlers.DestProductsTableSearch, handlers.PostProductsTableSearch)

	consoleMux.HandleFunc("GET "+handlers.DestSuppliers, handlers.GetSuppliers)
	consoleMux.HandleFunc("GET "+handlers.DestSuppliers+"{id}", handlers.GetSupplier)
	consoleMux.HandleFunc("POST "+handlers.DestSuppliers, handlers.PostSupplier)
	consoleMux.HandleFunc("GET "+handlers.DestSuppliersTable, handlers.GetSuppliersTable)

	consoleMux.HandleFunc("GET "+handlers.DestUsers, handlers.GetUsers)
	consoleMux.HandleFunc("GET "+handlers.DestUsers+"{id}", handlers.GetUser)
	consoleMux.HandleFunc("POST "+handlers.DestUsers, handlers.PostUser)
	consoleMux.HandleFunc("GET "+handlers.DestUsersTable, handlers.GetUsersTable)

	consoleMux.HandleFunc("GET "+handlers.DestUpload, handlers.GetUpload)
	consoleMux.HandleFunc("POST "+handlers.DestUpload, handlers.PostUpload)
	consoleMux.HandleFunc("POST "+handlers.DestUploadPreview, handlers.PostUploadPreview)

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("POST "+handlers.DestProductSearch, handlers.PostProductSearch)
	apiMux.HandleFunc("POST "+handlers.DestOrderAmountInput, handlers.PostOrderAmountInput)
	apiMux.HandleFunc("GET "+handlers.DestOwnOrdersView, handlers.GetOwnOrdersView)

	mux.HandleFunc("/", handlers.GetIndex)
	mux.HandleFunc("GET "+handlers.DestFirstLogin, handlers.GetFirstLogin)
	mux.HandleFunc("POST "+handlers.DestFirstLogin, handlers.PostFirstLogin)
	mux.HandleFunc("POST /login", handlers.PostLogin)
	mux.HandleFunc("GET /logout", handlers.Logout)
	mux.Handle("GET /public/", http.FileServerFS(publicFS))

	mux.Handle(handlers.DestChef, mw.WithLogging(mw.WithUserCheck(
		func(u *database.User) bool { return u.RoleID == database.RoleIDChef },
		chefMux,
	)))
	mux.Handle(handlers.DestConsole, mw.WithLogging(mw.WithUserCheck(
		func(u *database.User) bool {
			return u.RoleID == database.RoleIDManager || u.RoleID == database.RoleIDAdministrator
		},
		consoleMux,
	)))
	mux.Handle(handlers.DestApi, mw.WithUserCheck(
		func(user *database.User) bool { return true },
		apiMux,
	))
}
