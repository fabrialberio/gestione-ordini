package main

import "net/http"

func cook(w http.ResponseWriter, r *http.Request) {
	templ.ExecuteTemplate(w, "cook.html", nil)
}
