package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"log"
	"net/http"
	"strconv"
)

func managerSidebar(selected int) []components.SidebarDest {
	sidebar := []components.SidebarDest{
		{destManagerAllOrders, "fa-users", "Ordini", false},
		{destManagerProducts, "fa-box-open", "Prodotti", false},
	}
	sidebar[selected].Selected = true

	return sidebar
}

func GetManager(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, destManagerAllOrders, http.StatusSeeOther)
}

func GetManagerProducts(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: managerSidebar(1),
	}

	err := appContext.FromRequest(r).Templ.ExecuteTemplate(w, "managerProducts.html", data)
	if err != nil {
		log.Println(err)
	}
}

func GetManagerAllOrders(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: managerSidebar(0),
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "managerAllOrders.html", data)
}

func GetManagerProductsTable(w http.ResponseWriter, r *http.Request) {
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
		data.OrderBy = database.OrderUserByID
	}
	data.OrderDesc = r.URL.Query().Get("orderDesc") == "true"

	data.Products, err = appContext.FromRequest(r).DB.FindAllProducts()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "managerProductsTable.html", data)
}

func GetManagerProduct(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Product             database.Product
		NameInput           components.Input
		ProductTypeSelect   components.Select
		SupplierSelect      components.Select
		UnitOfMeasureSelect components.Select
		IsNew               bool
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		data.IsNew = true
		data.Product = database.Product{}
	} else {
		product, err := appContext.FromRequest(r).DB.FindProduct(id)
		if err != nil {
			HandleError(w, r, err)
			return
		}

		data.Product = product
	}

	data.NameInput = components.Input{"Nome", keyProductName, "text", data.Product.Name}

	productTypes, err := appContext.FromRequest(r).DB.FindAllProductTypes()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.ProductTypeSelect = components.Select{"Tipologia di prodotto", keyProductProductTypeID, data.Product.ProductTypeID, []components.SelectOption{}}
	for _, p := range productTypes {
		data.ProductTypeSelect.Options = append(data.ProductTypeSelect.Options, components.SelectOption{p.ID, p.Name})
	}

	suppliers, err := appContext.FromRequest(r).DB.FindAllSuppliers()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.SupplierSelect = components.Select{"Fornitore", keyProductSupplierID, data.Product.SupplierID, []components.SelectOption{}}
	for _, s := range suppliers {
		data.SupplierSelect.Options = append(data.SupplierSelect.Options, components.SelectOption{s.ID, "Nome del fornitore"})
	}

	unitsOfMeasure, err := appContext.FromRequest(r).DB.FindAllUnitsOfMeasure()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.UnitOfMeasureSelect = components.Select{"Unità di misura", keyProductUnitOfMeasureID, data.Product.UnitOfMeasureID, []components.SelectOption{}}
	for _, u := range unitsOfMeasure {
		data.UnitOfMeasureSelect.Options = append(data.UnitOfMeasureSelect.Options, components.SelectOption{u.ID, u.Symbol})
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "managerProduct.html", data)
}

func PostManagerProduct(w http.ResponseWriter, r *http.Request) {
	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")

	productTypeId, _ := strconv.Atoi(r.FormValue("productTypeId"))
	supplierId, _ := strconv.Atoi(r.FormValue("supplierId"))
	unitOfMeasureId, _ := strconv.Atoi(r.FormValue("unitOfMeasureId"))
	name := r.FormValue("requestedAt")

	if isNew {
		err := appContext.FromRequest(r).DB.CreateProduct(database.Product{
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
		id, err := strconv.Atoi(r.FormValue(keyOrderProductID))
		if err != nil {
			HandleError(w, r, err)
			return
		}

		if delete {
			err = appContext.FromRequest(r).DB.DeleteProduct(id)
			if err != nil {
				HandleError(w, r, err)
				return
			}
		} else {
			err = appContext.FromRequest(r).DB.UpdateProduct(database.Product{
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

	http.Redirect(w, r, destManagerProducts, http.StatusSeeOther)
}
