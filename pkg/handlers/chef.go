package handlers

import (
	"gestione-ordini/pkg/appContext"
	"net/http"
)

func GetChef(w http.ResponseWriter, r *http.Request) {
	appContext.ExecuteTemplate(w, r, "chef.html", nil)
}
