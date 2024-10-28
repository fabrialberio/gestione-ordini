package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"time"
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

type supplierOrders struct {
	Supplier   database.Supplier
	OrdersList components.OrdersList
}

func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar            []components.SidebarDest
		ExpiredOrdersList  components.OrdersList
		SupplierOrders     map[database.Supplier]components.OrdersList
		NextWeekOrdersList components.OrdersList
	}{
		Sidebar:        sidebarDestinations(r, 0),
		SupplierOrders: map[database.Supplier]components.OrdersList{},
	}

	orders, err := appContext.Database(r).FindAllOrders()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	oneWeek := time.Hour * 24 * 7

	for _, o := range orders {
		if time.Until(o.ExpiresAt) < 0 {
			data.ExpiredOrdersList.Orders = append(data.ExpiredOrdersList.Orders, o)
		} else if time.Until(o.ExpiresAt) > oneWeek {
			data.NextWeekOrdersList.Orders = append(data.NextWeekOrdersList.Orders, o)
		} else {
			orders := data.SupplierOrders[o.Product.Supplier].Orders
			orders = append(orders, o)
			data.SupplierOrders[o.Product.Supplier] = components.OrdersList{Orders: orders}
		}
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
