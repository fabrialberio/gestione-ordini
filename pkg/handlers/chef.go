package handlers

import (
	"gestione-ordini/pkg/appContext"
	"net/http"
)

func GetChef(w http.ResponseWriter, r *http.Request) {
	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "chef.html", nil)
}
