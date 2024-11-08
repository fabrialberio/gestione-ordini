package handlers

import (
	"gestione-ordini/pkg/appContext"
	"log"
	"net/http"
)

func ShowError(w http.ResponseWriter, r *http.Request, err error) {
	LogError(r, err)
	appContext.ExecuteTemplate(w, r, "error.html", err.Error())
}

func LogError(r *http.Request, err error) {
	log.Printf("%s %s %s Error: %v", r.RemoteAddr, r.Method, r.URL, err)
}
