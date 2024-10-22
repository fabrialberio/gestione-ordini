package handlers

import (
	"gestione-ordini/pkg/appContext"
	"log"
	"net/http"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s %s %s Error: %v", r.RemoteAddr, r.Method, r.URL, err)
	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "error.html", err.Error())
}

func RedirectHandler(redirectUrl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
	}
}
