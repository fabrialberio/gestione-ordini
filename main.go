package main

import (
	"embed"
	"fmt"
	"gestione-ordini/database"
	"gestione-ordini/router"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	db *database.Database

	//go:embed public
	publicFS embed.FS
)

func main() {
	checkEnvVars()

	templ := template.Must(template.ParseGlob("templates/*.html"))
	template.Must(templ.ParseGlob("templates/**/*.html"))

	db = createDatabase()
	defer db.Close()

	addAdminUserIfNotExists()

	rt := router.New(templ)

	rt.Get("/public/", http.FileServerFS(publicFS).ServeHTTP)

	rt.GetTemplate("/", "login.html", index)
	rt.GetTemplate("/cook", "cook.html", cook)
	rt.GetTemplate("/cook/order", "order.html", cookOrder)
	rt.GetTemplate("/cook/ordersList", "ordersList.html", cookOrdersList)
	rt.GetTemplate("/manager", "manager.html", manager)
	rt.GetTemplate("/manager/product", "product.html", managerProduct)
	rt.GetTemplate("/manager/productsTable", "productsTable.html", managerProductsTable)
	rt.GetTemplate("/admin", "admin.html", admin)
	rt.GetTemplate("/admin/user", "user.html", adminUser)
	rt.GetTemplate("/admin/users", "users.html", adminUsers)
	rt.GetTemplate("/admin/usersTable", "usersTable.html", adminUsersTable)

	rt.Post("/login", login)
	rt.Post("/logout", logout)
	rt.Post("/cook/order/edit", cookOrderEdit)
	rt.Post("/admin/user/edit", adminUserEdit)
	rt.Post("/manager/product/edit", managerProductEdit)

	log.Println("Server started on port 8080.")
	log.Fatal(rt.ListenAndServe(":8080"))
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
