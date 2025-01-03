package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"strconv"
)

func PostProductsTableSearch(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, DestProducts+r.FormValue("q"), http.StatusSeeOther)
}

func GetProductsTable(w http.ResponseWriter, r *http.Request) {
	query := components.ParseTableQuery(r, 21)

	products, err := appContext.Database(r).FindAllProducts(query.OrderBy, query.OrderDesc, query.MaxRowCount+1)
	if err != nil {
		logError(r, err)
	}

	headings := []components.TableHeading{
		{Index: database.OrderProductByID, Name: "ID"},
		{Index: database.OrderProductByDescription, Name: "Descrizione"},
		{Index: database.OrderProductByCode, Name: "Codice"},
		{Index: database.OrderProductByProductType, Name: "Tipologia"},
		{Index: database.OrderProductBySupplier, Name: "Fornitore"},
	}

	appContext.ExecuteTemplate(w, r, "table", components.ComposeTable(query, headings, productRowFunc, products))
}

func productRowFunc(product database.Product) components.TableRow {
	return components.TableRow{
		EditURL: DestProducts + "/" + strconv.Itoa(product.ID),
		Cells: []components.TableCell{
			{Value: strconv.Itoa(product.ID)},
			{Value: product.Description},
			{Value: product.Code},
			{Value: product.ProductType.Name},
			{Value: product.Supplier.Name, URL: DestSuppliers + "/" + strconv.Itoa(product.Supplier.ID)},
		},
	}
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
			ShowItemNotAllowedError(w, r, err)
			return
		}
	}

	productTypes, err := appContext.Database(r).FindAllProductTypes()
	if err != nil {
		ShowDatabaseQueryError(w, r, err)
		return
	}

	productTypeOptions := []components.SelectOption{}
	for _, p := range productTypes {
		productTypeOptions = append(productTypeOptions, components.SelectOption{Value: p.ID, Text: p.Name})
	}

	suppliers, err := appContext.Database(r).FindAllSuppliers(database.OrderSupplierByID, true)
	if err != nil {
		ShowDatabaseQueryError(w, r, err)
		return
	}

	suppilerOptions := []components.SelectOption{}
	for _, s := range suppliers {
		suppilerOptions = append(suppilerOptions, components.SelectOption{Value: s.ID, Text: s.Name})
	}

	unitsOfMeasure, err := appContext.Database(r).FindAllUnitsOfMeasure()
	if err != nil {
		ShowDatabaseQueryError(w, r, err)
		return
	}

	unitOfMeasureOptions := []components.SelectOption{}
	for _, u := range unitsOfMeasure {
		unitOfMeasureOptions = append(unitOfMeasureOptions, components.SelectOption{Value: u.ID, Text: u.Symbol})
	}

	data := struct {
		IsNew               bool
		Product             database.Product
		DescriptionInput    components.Input
		CodeInput           components.Input
		ProductTypeSelect   components.Select
		SupplierSelect      components.Select
		UnitOfMeasureSelect components.Select
	}{
		IsNew:   isNew,
		Product: defaultProduct,
		DescriptionInput: components.Input{
			Label:        "Descrizione",
			Name:         keyProductDescription,
			DefaultValue: defaultProduct.Description,
		},
		CodeInput: components.Input{
			Label:        "Codice",
			Name:         keyProductCode,
			DefaultValue: defaultProduct.Code,
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
	description := r.FormValue(keyProductDescription)
	code := r.FormValue(keyProductCode)

	if isNew {
		err := appContext.Database(r).CreateProduct(database.Product{
			ProductTypeID:   productTypeId,
			SupplierID:      supplierId,
			UnitOfMeasureID: unitOfMeasureId,
			Description:     description,
			Code:            code,
		})
		if err != nil {
			ShowDatabaseQueryError(w, r, err)
			return
		}
	} else {
		id, err := strconv.Atoi(r.FormValue(keyProductID))
		if err != nil {
			ShowItemInvalidFormError(w, r, err)
			return
		}

		if delete {
			err = appContext.Database(r).DeleteProduct(id)
			if err != nil {
				ShowItemNotDeletableError(w, r, err)
				return
			}
		} else {
			err = appContext.Database(r).UpdateProduct(database.Product{
				ID:              id,
				ProductTypeID:   productTypeId,
				SupplierID:      supplierId,
				UnitOfMeasureID: unitOfMeasureId,
				Description:     description,
				Code:            code,
			})
			if err != nil {
				ShowDatabaseQueryError(w, r, err)
				return
			}
		}
	}

	http.Redirect(w, r, DestProducts, http.StatusSeeOther)
}
