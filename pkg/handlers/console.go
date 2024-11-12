package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
)

const (
	sidebarIndexNewOrder int = iota
	sidebarIndexAllOrders
	sidebarIndexProducts
	sidebarIndexSuppliers
	sidebarIndexUsers
	sidebarIndexUpload
)

func currentSidebar(selected int, isAdmin bool) []components.SidebarDest {
	sidebar := []components.SidebarDest{
		{
			DestURL:     DestNewOrder,
			FasIconName: "fa-circle-plus",
			Label:       "Nuovo ordine",
		},
		{
			DestURL:     DestAllOrders,
			FasIconName: "fa-calendar-check",
			Label:       "Tutti gli ordini",
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

	if isAdmin {
		sidebar = append(sidebar,
			components.SidebarDest{
				DestURL:     DestUsers,
				FasIconName: "fa-users",
				Label:       "Utenti",
			},
			components.SidebarDest{
				DestURL:     DestUpload,
				FasIconName: "fa-file-upload",
				Label:       "Importa",
			},
		)
	} else if selected > sidebarIndexSuppliers {
		return sidebar
	}

	sidebar[selected].Selected = true
	return sidebar
}

func GetConsole(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, DestAllOrders, http.StatusSeeOther)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	user, err := appContext.AuthenticatedUser(r)
	if err != nil {
		LogoutError(w, r, err)
		return
	}

	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: currentSidebar(sidebarIndexProducts, user.RoleID == database.RoleIDAdministrator),
	}

	appContext.ExecuteTemplate(w, r, "products.html", data)
}

func GetSuppliers(w http.ResponseWriter, r *http.Request) {
	user, err := appContext.AuthenticatedUser(r)
	if err != nil {
		LogoutError(w, r, err)
		return
	}

	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: currentSidebar(sidebarIndexSuppliers, user.RoleID == database.RoleIDAdministrator),
	}

	appContext.ExecuteTemplate(w, r, "suppliers.html", data)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	user, err := appContext.AuthenticatedUser(r)
	if err != nil {
		LogoutError(w, r, err)
		return
	} else if user.RoleID != database.RoleIDAdministrator {
		ShowItemNotAllowedError(w, r, auth.ErrInvalidRole)
		return
	}

	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: currentSidebar(sidebarIndexUsers, true),
	}

	appContext.ExecuteTemplate(w, r, "users.html", data)
}
