package main

import (
	"log"
	"net/http"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s %s %s Error: %v", r.RemoteAddr, r.Method, r.URL, err)
	GetRequestContext(r).Templ.ExecuteTemplate(w, "error.html", err.Error())
}
