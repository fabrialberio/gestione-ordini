package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"log"
	"net/http"
)

func managerSidebar(selected int) []components.SidebarDest {
	sidebar := []components.SidebarDest{
		{destManagerAllOrders, "fa-users", "Ordini", false},
		{destManagerProducts, "fa-box-open", "Prodotti", false},
	}
	sidebar[selected].Selected = true

	return sidebar
}

func GetManager(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, destManagerAllOrders, http.StatusSeeOther)
}

func GetManagerProductsTable(w http.ResponseWriter, r *http.Request) {
	getProductsTable(w, r, destManagerProductsTable)
}

func GetManagerProducts(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: managerSidebar(1),
	}

	err := appContext.FromRequest(r).Templ.ExecuteTemplate(w, "managerProducts.html", data)
	if err != nil {
		log.Println(err)
	}
}

func GetManagerAllOrders(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: managerSidebar(0),
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "managerAllOrders.html", data)
}

func PostManagerProduct(w http.ResponseWriter, r *http.Request) {
	postProduct(w, r)

	http.Redirect(w, r, destManagerProducts, http.StatusSeeOther)
}
