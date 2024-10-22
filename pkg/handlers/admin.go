package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"net/http"
)

func adminSidebar(selected int) []components.SidebarDest {
	sidebar := []components.SidebarDest{
		{destAdminUsers, "fa-users", "Utenti", false},
		{destAdminProducts, "fa-box", "Prodotti", false},
		{destAdminSuppliers, "fa-store", "Fornitori", false},
	}
	sidebar[selected].Selected = true

	return sidebar
}

func GetAdmin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, destAdminUsers, http.StatusSeeOther)
}

func GetAdminUsers(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: adminSidebar(0),
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "adminUsers.html", data)
}

func PostAdminUser(w http.ResponseWriter, r *http.Request) {
	postUser(w, r)

	http.Redirect(w, r, destAdminUsers, http.StatusSeeOther)
}

func GetAdminProductsTable(w http.ResponseWriter, r *http.Request) {
	getProductsTable(w, r, destAdminProductsTable, destAdminProducts)
}

func GetAdminProducts(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: adminSidebar(1),
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "adminProducts.html", data)
}

func PostAdminProduct(w http.ResponseWriter, r *http.Request) {
	postProduct(w, r)

	http.Redirect(w, r, destAdminProducts, http.StatusSeeOther)
}

func GetAdminSuppliers(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: adminSidebar(2),
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "adminSuppliers.html", data)
}

func GetAdminSuppliersTable(w http.ResponseWriter, r *http.Request) {
	getSuppliersTable(w, r, destAdminSuppliersTable, destAdminSuppliers)
}

func PostAdminSupplier(w http.ResponseWriter, r *http.Request) {
	postSupplier(w, r)

	http.Redirect(w, r, destAdminSuppliers, http.StatusSeeOther)
}
