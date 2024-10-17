package main

import (
	"gestione-ordini/database"
	"log"
	"net/http"
	"strconv"
)

func HandleGetManager(w http.ResponseWriter, r *http.Request) {
	templ.ExecuteTemplate(w, "manager.html", nil)
}

func HandleGetManagerProductsTable(w http.ResponseWriter, r *http.Request) {
	var err error
	var data struct {
		OrderBy   int
		OrderDesc bool
		Headers   interface{}
		Products  []database.Product
	}

	data.Headers = []struct {
		Index int
		Name  string
	}{}

	data.OrderBy, err = strconv.Atoi(r.URL.Query().Get("orderBy"))
	if err != nil {
		data.OrderBy = database.UserOrderByID
	}
	data.OrderDesc = r.URL.Query().Get("orderDesc") == "true"

	data.Products, err = db.FindAllProducts()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	templ.ExecuteTemplate(w, "productsTable.html", data)
}

func HandleGetManagerProduct(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Product        database.Product
		ProductTypes   []database.ProductType
		Suppliers      []database.Supplier
		UnitsOfMeasure []database.UnitOfMeasure
		IsNew          bool
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		data.IsNew = true
		data.Product = database.Product{}
	} else {
		product, err := db.FindProduct(id)
		if err != nil {
			HandleError(w, r, err)
			return
		}

		data.Product = product
	}

	data.ProductTypes, err = db.FindAllProductTypes()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.Suppliers, err = db.FindAllSuppliers()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.UnitsOfMeasure, err = db.FindAllUnitsOfMeasure()
	log.Println(data.UnitsOfMeasure)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	templ.ExecuteTemplate(w, "product.html", data)
}

func HandlePostManagerProduct(w http.ResponseWriter, r *http.Request) {
	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")

	productTypeId, _ := strconv.Atoi(r.FormValue("productTypeId"))
	supplierId, _ := strconv.Atoi(r.FormValue("supplierId"))
	unitOfMeasureId, _ := strconv.Atoi(r.FormValue("unitOfMeasureId"))
	name := r.FormValue("requestedAt")

	if isNew {
		err := db.CreateProduct(database.Product{
			ProductTypeID:   productTypeId,
			SupplierID:      supplierId,
			UnitOfMeasureID: unitOfMeasureId,
			Name:            name,
		})
		if err != nil {
			HandleError(w, r, err)
			return
		}
	} else {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			HandleError(w, r, err)
			return
		}

		if delete {
			err = db.DeleteProduct(id)
			if err != nil {
				HandleError(w, r, err)
				return
			}
		} else {
			err = db.UpdateProduct(database.Product{
				ID:              id,
				ProductTypeID:   productTypeId,
				SupplierID:      supplierId,
				UnitOfMeasureID: unitOfMeasureId,
				Name:            name,
			})
			if err != nil {
				HandleError(w, r, err)
				return
			}
		}
	}

	http.Redirect(w, r, "/manager", http.StatusSeeOther)
}
