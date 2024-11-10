package handlers

import (
	"gestione-ordini/pkg/appContext"
	"net/http"
)

func GetUploadProducts(w http.ResponseWriter, r *http.Request) {
	appContext.ExecuteTemplate(w, r, "uploadProducts.html", nil)
}

func PostUploadProducts(w http.ResponseWriter, r *http.Request) {

}
