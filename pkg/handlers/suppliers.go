package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"strconv"
)

func GetSuppliersTable(w http.ResponseWriter, r *http.Request) {
	orderBy, err := strconv.Atoi(r.URL.Query().Get("orderBy"))
	if err != nil {
		orderBy = database.OrderSupplierByID
	}
	orderDesc := r.URL.Query().Get("orderDesc") == "true"

	suppliers, err := appContext.Database(r).FindAllSuppliers(orderBy, orderDesc)
	if err != nil {
		ShowError(w, r, err)
		return
	}

	data := components.SuppliersTable{
		Table: components.Table{
			TableURL:  DestSuppliersTable,
			OrderBy:   orderBy,
			OrderDesc: orderDesc,
			Headings: []components.TableHeading{
				{Index: database.OrderSupplierByID, Name: "ID"},
				{Index: database.OrderSupplierByName, Name: "Nome"},
				{Index: database.OrderSupplierByEmail, Name: "Email"},
			},
		},
		Suppliers: suppliers,
	}

	appContext.ExecuteTemplate(w, r, "suppliersTable", data)
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
			ShowError(w, r, err)
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
			Type:         "text",
			DefaultValue: defaultSupplier.Name,
		},
		EmailInput: components.Input{
			Label:        "Email",
			Name:         keySupplierEmail,
			Type:         "text",
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
			ShowError(w, r, err)
			return
		}
	} else {
		id, err := strconv.Atoi(r.FormValue(keySupplierID))
		if err != nil {
			ShowError(w, r, err)
			return
		}

		if delete {
			err := appContext.Database(r).DeleteSupplier(id)
			if err != nil {
				ShowError(w, r, err)
				return
			}
		} else {
			err := appContext.Database(r).UpdateSupplier(database.Supplier{
				ID:    id,
				Email: email,
				Name:  name,
			})
			if err != nil {
				ShowError(w, r, err)
				return
			}
		}
	}

	http.Redirect(w, r, DestSuppliers, http.StatusSeeOther)
}
