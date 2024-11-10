package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
)

func GetUpload(w http.ResponseWriter, r *http.Request) {
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
		Sidebar: currentSidebar(4, true),
	}

	appContext.ExecuteTemplate(w, r, "upload.html", data)
}

func PostUpload(w http.ResponseWriter, r *http.Request) {

}
