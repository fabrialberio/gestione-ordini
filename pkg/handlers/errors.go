package handlers

import (
	"gestione-ordini/pkg/appContext"
	"log"
	"net/http"
)

func ErrorRedirect(w http.ResponseWriter, r *http.Request, err error) {
	LogError(r, err)
	http.Redirect(w, r, "/error", http.StatusSeeOther)
}

func GetError(w http.ResponseWriter, r *http.Request) {
	appContext.ExecuteTemplate(w, r, "error.html", "Errore interno")
}

func ShowError(w http.ResponseWriter, r *http.Request, err error) {
	LogError(r, err)
	appContext.ExecuteTemplate(w, r, "error.html", err.Error())
}

func LogError(r *http.Request, err error) {
	log.Printf("%s %s %s Error: %v", r.RemoteAddr, r.Method, r.URL, err)
}
