package router

import (
	"html/template"
	"log"
	"net/http"
)

const (
	errorTemplateName = "error.html"
)

type Router struct {
	mux   *http.ServeMux
	templ *template.Template
}

type TemplateDataFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

type PostFunc func(w http.ResponseWriter, r *http.Request) error

func (rt *Router) htmlError(w http.ResponseWriter, r *http.Request, message string) {
	log.Printf("%s %s %s Error: %v", r.RemoteAddr, r.Method, r.URL, message)
	rt.templ.ExecuteTemplate(w, errorTemplateName, message)
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

func (rt *Router) Get(pattern string, getFunc http.HandlerFunc) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			rt.htmlError(w, r, "Method not allowed")
			return
		}

		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		getFunc(w, r)
	}

	rt.mux.HandleFunc(pattern, handlerFunc)
}

func (rt *Router) Post(pattern string, postFunc PostFunc) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			rt.htmlError(w, r, "Method not allowed")
			return
		}

		err := postFunc(w, r)
		if err != nil {
			rt.htmlError(w, r, err.Error())
			return
		}

		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	}

	rt.mux.HandleFunc(pattern, handlerFunc)
}

func (rt *Router) GetTemplate(pattern string, templateName string, dataFunc TemplateDataFunc) {
	getFunc := func(w http.ResponseWriter, r *http.Request) {
		data, err := dataFunc(w, r)
		if err != nil {
			rt.htmlError(w, r, err.Error())
			return
		}

		rt.templ.ExecuteTemplate(w, templateName, data)
	}

	rt.Get(pattern, getFunc)
}
