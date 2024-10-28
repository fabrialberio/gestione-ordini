package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
)

func sidebarDestinations(r *http.Request, selected int) []components.SidebarDest {
	sidebar := []components.SidebarDest{
		{DestAllOrders, "fa-clipboard-check", "Ordini", false},
		{DestProducts, "fa-box", "Prodotti", false},
		{DestSuppliers, "fa-truck", "Fornitori", false},
	}

	user, _ := appContext.AuthenticatedUser(r)
	if user != nil && user.RoleID == database.RoleIDAdministrator {
		sidebar = append(sidebar, components.SidebarDest{DestUsers, "fa-users", "Utenti", false})
	}

	sidebar[selected].Selected = true

	return sidebar
}

func GetConsole(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, DestAllOrders, http.StatusSeeOther)
}

func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	var err error
	data := struct {
		Sidebar    []components.SidebarDest
		OrdersList components.OrdersList
	}{
		Sidebar: sidebarDestinations(r, 0),
	}

	data.OrdersList.Orders, err = appContext.Database(r).FindAllOrders()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	appContext.ExecuteTemplate(w, r, "allOrders.html", data)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: sidebarDestinations(r, 1),
	}

	appContext.ExecuteTemplate(w, r, "products.html", data)
}

func GetSuppliers(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: sidebarDestinations(r, 2),
	}

	appContext.ExecuteTemplate(w, r, "suppliers.html", data)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: sidebarDestinations(r, 3),
	}

	appContext.ExecuteTemplate(w, r, "users.html", data)
}
