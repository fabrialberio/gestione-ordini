package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"strconv"
)

func GetProductsTable(w http.ResponseWriter, r *http.Request) {
	var err error
	data := components.ProductsTable{
		Table: components.Table{
			TableURL: DestProductsTable,
		},
	}

	data.Table.Headings = []components.TableHeading{
		{database.OrderProductByID, "ID"},
		{database.OrderProductByName, "Nome"},
		{database.OrderProductByProductType, "Tipologia"},
		{database.OrderProductBySupplier, "Fornitore"},
		{database.OrderProductByUnitOfMeasure, "Unità"},
	}

	data.Table.OrderBy, err = strconv.Atoi(r.URL.Query().Get("orderBy"))
	if err != nil {
		data.Table.OrderBy = database.OrderProductByID
	}
	data.Table.OrderDesc = r.URL.Query().Get("orderDesc") == "true"

	data.Products, err = appContext.Database(r).FindAllProducts(data.Table.OrderBy, data.Table.OrderDesc)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	appContext.ExecuteTemplate(w, r, "productsTable.html", data)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
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
		product, err := appContext.Database(r).FindProduct(id)
		if err != nil {
			HandleError(w, r, err)
			return
		}

		data.Product = product
	}

	data.NameInput = components.Input{"Nome", keyProductName, "text", data.Product.Name}

	productTypes, err := appContext.Database(r).FindAllProductTypes()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.ProductTypeSelect = components.Select{"Tipologia di prodotto", keyProductProductTypeID, data.Product.ProductTypeID, []components.SelectOption{}}
	for _, p := range productTypes {
		data.ProductTypeSelect.Options = append(data.ProductTypeSelect.Options, components.SelectOption{p.ID, p.Name})
	}

	suppliers, err := appContext.Database(r).FindAllSuppliers(database.OrderSupplierByID, true)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.SupplierSelect = components.Select{"Fornitore", keyProductSupplierID, data.Product.SupplierID, []components.SelectOption{}}
	for _, s := range suppliers {
		data.SupplierSelect.Options = append(data.SupplierSelect.Options, components.SelectOption{s.ID, s.Name})
	}

	unitsOfMeasure, err := appContext.Database(r).FindAllUnitsOfMeasure()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.UnitOfMeasureSelect = components.Select{"Unità di misura", keyProductUnitOfMeasureID, data.Product.UnitOfMeasureID, []components.SelectOption{}}
	for _, u := range unitsOfMeasure {
		data.UnitOfMeasureSelect.Options = append(data.UnitOfMeasureSelect.Options, components.SelectOption{u.ID, u.Symbol})
	}

	appContext.ExecuteTemplate(w, r, "product.html", data)
}

func PostProduct(w http.ResponseWriter, r *http.Request) {
	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")

	productTypeId, _ := strconv.Atoi(r.FormValue(keyProductProductTypeID))
	supplierId, _ := strconv.Atoi(r.FormValue(keyProductSupplierID))
	unitOfMeasureId, _ := strconv.Atoi(r.FormValue(keyProductUnitOfMeasureID))
	name := r.FormValue(keyProductName)

	if isNew {
		err := appContext.Database(r).CreateProduct(database.Product{
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
		id, err := strconv.Atoi(r.FormValue(keyProductID))
		if err != nil {
			HandleError(w, r, err)
			return
		}

		if delete {
			err = appContext.Database(r).DeleteProduct(id)
			if err != nil {
				HandleError(w, r, err)
				return
			}
		} else {
			err = appContext.Database(r).UpdateProduct(database.Product{
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

	http.Redirect(w, r, DestProducts, http.StatusSeeOther)
}
