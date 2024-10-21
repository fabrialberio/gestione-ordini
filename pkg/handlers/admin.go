package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"net/http"
)

func adminSidebar(selected int) []components.SidebarDest {
	sidebar := []components.SidebarDest{
		{destAdminUsers, "fa-users", "Utenti", false},
		{destAdminProducts, "fa-box-open", "Prodotti", false},
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
