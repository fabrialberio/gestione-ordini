package main

import (
	"log"
	"net/http"
	"text/template"
)

func logHandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next(w, r)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name string
	}{
		Name: "world",
	}

	t := template.Must(template.ParseFiles("templates/index.html"))
	t.Execute(w, data)
}
