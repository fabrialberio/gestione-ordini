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
		{
			DestURL:     DestAllOrders,
			FasIconName: "fa-calendar-check",
			Label:       "Ordini",
		},
		{
			DestURL:     DestProducts,
			FasIconName: "fa-box",
			Label:       "Prodotti",
		},
		{
			DestURL:     DestSuppliers,
			FasIconName: "fa-truck",
			Label:       "Fornitori",
		},
	}

	user, _ := appContext.AuthenticatedUser(r)
	if user != nil && user.RoleID == database.RoleIDAdministrator {
		sidebar = append(sidebar, components.SidebarDest{
			DestURL:     DestUsers,
			FasIconName: "fa-users",
			Label:       "Utenti",
		})
	}

	sidebar[selected].Selected = true

	return sidebar
}

func GetConsole(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, DestAllOrders, http.StatusSeeOther)
}

func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	defaultStart := currentWeekStart()
	weekDuration := time.Hour * 24 * 6

	suppliers, err := appContext.Database(r).FindAllSuppliers(database.OrderSupplierByID, true)
	if err != nil {
		ShowError(w, r, err)
		return
	}

	supplierOptions := []components.SelectOption{{
		Value: 0,
		Text:  "Tutti i fornitori",
	}}
	for _, s := range suppliers {
		supplierOptions = append(supplierOptions, components.SelectOption{Value: s.ID, Text: s.Name})
	}

	data := struct {
		Sidebar        []components.SidebarDest
		StartDateInput components.Input
		EndDateInput   components.Input
		SupplierSelect components.Select
	}{
		Sidebar: sidebarDestinations(r, 0),
		StartDateInput: components.Input{
			Label:        "Da",
			Name:         keyOrderSelectionStart,
			Type:         "date",
			DefaultValue: defaultStart.Format(dateFormat),
		},
		EndDateInput: components.Input{
			Label:        "A",
			Name:         keyOrderSelectionEnd,
			Type:         "date",
			DefaultValue: defaultStart.Add(weekDuration).Format(dateFormat),
		},
		SupplierSelect: components.Select{
			Label:    "Fornitore",
			Name:     keyOrderSelectionSupplierID,
			Selected: 0,
			Options:  supplierOptions,
		},
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
