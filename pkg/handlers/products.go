package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"log"
	"net/http"
	"strconv"
)

func GetProductsTable(w http.ResponseWriter, r *http.Request) {
	orderBy, err := strconv.Atoi(r.URL.Query().Get("orderBy"))
	if err != nil {
		orderBy = database.OrderProductByID
	}
	orderDesc := r.URL.Query().Get("orderDesc") == "true"

	products, err := appContext.Database(r).FindAllProducts(orderBy, orderDesc)
	if err != nil {
		ShowError(w, r, err)
		return
	}

	data := components.ProductsTable{
		Table: components.Table{
			TableURL:  DestProductsTable,
			OrderBy:   orderBy,
			OrderDesc: orderDesc,
			Headings: []components.TableHeading{
				{Index: database.OrderProductByID, Name: "ID"},
				{Index: database.OrderProductByName, Name: "Nome"},
				{Index: database.OrderProductByProductType, Name: "Tipologia"},
				{Index: database.OrderProductBySupplier, Name: "Fornitore"},
				{Index: database.OrderProductByUnitOfMeasure, Name: "Unità"},
			},
		},
		Products: products,
	}

	appContext.ExecuteTemplate(w, r, "productsTable", data)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	isNew := false
	defaultProduct := database.Product{}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		isNew = true
	} else {
		defaultProduct, err = appContext.Database(r).FindProduct(id)
		if err != nil {
			ShowError(w, r, err)
			return
		}
	}

	productTypes, err := appContext.Database(r).FindAllProductTypes()
	if err != nil {
		ShowError(w, r, err)
		return
	}

	productTypeOptions := []components.SelectOption{}
	for _, p := range productTypes {
		productTypeOptions = append(productTypeOptions, components.SelectOption{Value: p.ID, Text: p.Name})
	}

	suppliers, err := appContext.Database(r).FindAllSuppliers(database.OrderSupplierByID, true)
	if err != nil {
		ShowError(w, r, err)
		return
	}

	suppilerOptions := []components.SelectOption{}
	for _, s := range suppliers {
		suppilerOptions = append(suppilerOptions, components.SelectOption{Value: s.ID, Text: s.Name})
	}

	unitsOfMeasure, err := appContext.Database(r).FindAllUnitsOfMeasure()
	if err != nil {
		ShowError(w, r, err)
		return
	}

	unitOfMeasureOptions := []components.SelectOption{}
	for _, u := range unitsOfMeasure {
		unitOfMeasureOptions = append(unitOfMeasureOptions, components.SelectOption{Value: u.ID, Text: u.Symbol})
	}

	data := struct {
		IsNew               bool
		Product             database.Product
		NameInput           components.Input
		ProductTypeSelect   components.Select
		SupplierSelect      components.Select
		UnitOfMeasureSelect components.Select
	}{
		IsNew:   isNew,
		Product: defaultProduct,
		NameInput: components.Input{
			Label:        "Nome",
			Name:         keyProductName,
			Type:         "text",
			DefaultValue: defaultProduct.Name,
		},
		ProductTypeSelect: components.Select{
			Label:    "Tipologia di prodotto",
			Name:     keyProductProductTypeID,
			Selected: defaultProduct.ProductTypeID,
			Options:  productTypeOptions,
		},
		SupplierSelect: components.Select{
			Label:    "Fornitore",
			Name:     keyProductSupplierID,
			Selected: defaultProduct.SupplierID,
			Options:  suppilerOptions,
		},
		UnitOfMeasureSelect: components.Select{
			Label:    "Unità di misura",
			Name:     keyProductUnitOfMeasureID,
			Selected: defaultProduct.UnitOfMeasureID,
			Options:  unitOfMeasureOptions,
		},
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
			ShowError(w, r, err)
			return
		}
	} else {
		id, err := strconv.Atoi(r.FormValue(keyProductID))
		if err != nil {
			ShowError(w, r, err)
			return
		}

		if delete {
			err = appContext.Database(r).DeleteProduct(id)
			if err != nil {
				ShowError(w, r, err)
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
				ShowError(w, r, err)
				return
			}
		}
	}

	http.Redirect(w, r, DestProducts, http.StatusSeeOther)
}

func PostProductSearch(w http.ResponseWriter, r *http.Request) {
	log.Println("Searching products...")

	appContext.ExecuteTemplate(w, r, "productSearchDialog", nil)
}
