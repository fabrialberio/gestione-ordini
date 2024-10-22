package main

import (
	"embed"
	"fmt"
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/database"
	"gestione-ordini/pkg/handlers"
	"gestione-ordini/pkg/middleware"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed public
var publicFS embed.FS

func main() {
	templ := template.Must(template.ParseGlob("templates/*.html"))
	template.Must(templ.ParseGlob("templates/**/*.html"))

	db := createDatabase()
	defer db.Close()
	addAdminUserIfNotExists(db)

	appCtx := appContext.AppContext{
		DB:    db,
		Templ: templ,
	}

	cookMux := http.NewServeMux()
	cookMux.HandleFunc("GET "+handlers.DestCook, handlers.GetCook)
	cookMux.HandleFunc("GET "+handlers.DestCookOrdersList, handlers.GetCookOrdersList)
	cookMux.HandleFunc("GET "+handlers.DestCookOrders+"{id}", handlers.GetCookOrder)
	cookMux.HandleFunc("POST "+handlers.DestCookOrders, handlers.PostCookOrder)

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("GET "+handlers.DestAdmin, handlers.RedirectHandler(handlers.DestAdminUsers))
	adminMux.HandleFunc("GET "+handlers.DestAdminUsers, handlers.GetAdminUsers)
	adminMux.HandleFunc("GET "+handlers.DestAdminUsers+"{id}", handlers.GetUser(handlers.DestAdminUsers))
	adminMux.HandleFunc("POST "+handlers.DestAdminUsers, handlers.PostAdminUser)
	adminMux.HandleFunc("GET "+handlers.DestAdminUsersTable, handlers.GetUsersTable)
	adminMux.HandleFunc("GET "+handlers.DestAdminProducts, handlers.GetAdminProducts)
	adminMux.HandleFunc("GET "+handlers.DestAdminProducts+"{id}", handlers.GetProduct(handlers.DestAdminProducts))
	adminMux.HandleFunc("POST "+handlers.DestAdminProducts, handlers.PostProduct(handlers.DestAdminProducts))
	adminMux.HandleFunc("GET "+handlers.DestAdminProductsTable, handlers.GetProductsTable(handlers.DestAdminProductsTable, handlers.DestAdminProducts))
	adminMux.HandleFunc("GET "+handlers.DestAdminSuppliers, handlers.GetAdminSuppliers)
	adminMux.HandleFunc("GET "+handlers.DestAdminSuppliers+"{id}", handlers.GetSupplier(handlers.DestAdminSuppliers))
	adminMux.HandleFunc("POST "+handlers.DestAdminSuppliers, handlers.PostAdminSupplier)
	adminMux.HandleFunc("GET "+handlers.DestAdminSuppliersTable, handlers.GetAdminSuppliersTable)

	managerMux := http.NewServeMux()
	managerMux.HandleFunc("GET "+handlers.DestManager, handlers.RedirectHandler(handlers.DestManagerAllOrders))
	managerMux.HandleFunc("GET "+handlers.DestManagerAllOrders, handlers.GetManagerAllOrders)
	managerMux.HandleFunc("GET "+handlers.DestManagerProducts, handlers.GetManagerProducts)
	managerMux.HandleFunc("GET "+handlers.DestManagerProducts+"{id}", handlers.GetProduct(handlers.DestManagerProducts))
	managerMux.HandleFunc("POST "+handlers.DestManagerProducts, handlers.PostProduct(handlers.DestManagerProducts))
	managerMux.HandleFunc("GET "+handlers.DestManagerProductsTable, handlers.GetProductsTable(handlers.DestManagerProductsTable, handlers.DestManagerProducts))

	mux := http.NewServeMux()
	mux.Handle("GET /public/", http.FileServerFS(publicFS))
	mux.Handle(handlers.DestCook, middleware.WithRole(database.RoleIDCook, cookMux))
	mux.Handle(handlers.DestAdmin, middleware.WithRole(database.RoleIDAdministrator, adminMux))
	mux.Handle(handlers.DestManager, middleware.WithRole(database.RoleIDManager, managerMux))

	mux.HandleFunc("/", handlers.GetIndex)
	mux.HandleFunc("POST /login", handlers.PostLogin)
	mux.HandleFunc("POST /logout", handlers.PostLogout)

	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.WithLogging(middleware.WithContext(&appCtx, mux)),
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
