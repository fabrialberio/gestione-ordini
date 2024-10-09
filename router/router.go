package router

import (
	"html/template"
	"log"
	"net/http"
)

type Router struct {
	mux   *http.ServeMux
	templ *template.Template
}

type TemplateDataFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

type PostFunc func(w http.ResponseWriter, r *http.Request) error

func htmlError(w http.ResponseWriter, r *http.Request, message string, status int) {
	log.Printf("%s %s %s Error: %v", r.RemoteAddr, r.Method, r.URL, message)
	http.Error(w, message, status)
}

func NewRouter(templ *template.Template) Router {
	return Router{
		mux:   http.NewServeMux(),
		templ: templ,
	}
}

func (rt *Router) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, rt.mux)
}

func (rt *Router) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	rt.mux.HandleFunc(pattern, handlerFunc)
}

func (rt *Router) HandleTemplate(pattern string, templateName string, dataFunc TemplateDataFunc) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			htmlError(w, r, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data, err := dataFunc(w, r)
		if err != nil {
			htmlError(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		rt.templ.ExecuteTemplate(w, templateName, data)
	}

	rt.mux.HandleFunc(pattern, handlerFunc)
}

func (rt *Router) HandlePost(pattern string, postFunc PostFunc) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			htmlError(w, r, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := postFunc(w, r)
		if err != nil {
			htmlError(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	}

	rt.mux.HandleFunc(pattern, handlerFunc)
}
