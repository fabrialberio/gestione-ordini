package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"log"
	"net/http"
)

func managerSidebar(selected int) []components.SidebarDest {
	sidebar := []components.SidebarDest{
		{DestManagerAllOrders, "fa-users", "Ordini", false},
		{DestManagerProducts, "fa-box-open", "Prodotti", false},
	}
	sidebar[selected].Selected = true

	return sidebar
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
