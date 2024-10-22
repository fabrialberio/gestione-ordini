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

	if appContext.FromRequest(r).AuthenticatedUser.RoleID == database.RoleIDAdministrator {
		sidebar = append(sidebar, components.SidebarDest{DestUsers, "fa-users", "Utenti", false})
	}

	sidebar[selected].Selected = true

	return sidebar
}

func GetConsole(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, DestAllOrders, http.StatusSeeOther)
}

func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: sidebarDestinations(r, 0),
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "allOrders.html", data)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: sidebarDestinations(r, 1),
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "products.html", data)
}

func GetSuppliers(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: sidebarDestinations(r, 2),
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "suppliers.html", data)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: sidebarDestinations(r, 3),
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "users.html", data)
}
