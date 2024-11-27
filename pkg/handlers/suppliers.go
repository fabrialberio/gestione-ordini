package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"strconv"
)

func GetSuppliersTable(w http.ResponseWriter, r *http.Request) {
	query := components.ParseTableQuery(r, 0)

	suppliers, err := appContext.Database(r).FindAllSuppliers(query.OrderBy, query.OrderDesc)
	if err != nil {
		logError(r, err)
	}

	headings := []components.TableHeading{
		{Index: database.OrderSupplierByID, Name: "ID"},
		{Index: database.OrderSupplierByName, Name: "Nome"},
		{Index: database.OrderSupplierByEmail, Name: "Email"},
	}

	rowFunc := func(supplier database.Supplier) components.TableRow {
		return components.TableRow{
			EditURL: DestSuppliers + "/" + strconv.Itoa(supplier.ID),
			Cells: []string{
				strconv.Itoa(supplier.ID),
				supplier.Name,
				supplier.Email,
			},
		}
	}

	appContext.ExecuteTemplate(w, r, "table", components.ComposeTable(query, headings, rowFunc, suppliers))
}

func GetSupplier(w http.ResponseWriter, r *http.Request) {
	isNew := false
	defaultSupplier := database.Supplier{}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		isNew = true
	} else {
		defaultSupplier, err = appContext.Database(r).FindSupplier(id)
		if err != nil {
			ShowItemNotAllowedError(w, r, err)
			return
		}
	}

	data := struct {
		IsNew      bool
		Supplier   database.Supplier
		NameInput  components.Input
		EmailInput components.Input
	}{
		IsNew:    isNew,
		Supplier: defaultSupplier,
		NameInput: components.Input{
			Label:        "Nome",
			Name:         keySupplierName,
			DefaultValue: defaultSupplier.Name,
		},
		EmailInput: components.Input{
			Label:        "Email",
			Name:         keySupplierEmail,
			DefaultValue: defaultSupplier.Email,
		},
	}

	appContext.ExecuteTemplate(w, r, "supplier.html", data)
}

func PostSupplier(w http.ResponseWriter, r *http.Request) {
	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")

	email := r.FormValue(keySupplierEmail)
	name := r.FormValue(keySupplierName)

	if isNew {
		err := appContext.Database(r).CreateSupplier(database.Supplier{
			Email: email,
			Name:  name,
		})
		if err != nil {
			ShowDatabaseQueryError(w, r, err)
			return
		}
	} else {
		id, err := strconv.Atoi(r.FormValue(keySupplierID))
		if err != nil {
			ShowItemInvalidFormError(w, r, err)
			return
		}

		if delete {
			err := appContext.Database(r).DeleteSupplier(id)
			if err != nil {
				ShowItemNotDeletableError(w, r, err)
				return
			}
		} else {
			err := appContext.Database(r).UpdateSupplier(database.Supplier{
				ID:    id,
				Email: email,
				Name:  name,
			})
			if err != nil {
				ShowDatabaseQueryError(w, r, err)
				return
			}
		}
	}

	http.Redirect(w, r, DestSuppliers, http.StatusSeeOther)
}
