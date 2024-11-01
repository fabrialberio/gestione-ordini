package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"strconv"
)

func GetSuppliersTable(w http.ResponseWriter, r *http.Request) {
	var err error
	data := components.SuppliersTable{
		Table: components.Table{
			TableURL: DestSuppliersTable,
		},
	}

	data.Table.Headings = []components.TableHeading{
		{database.OrderSupplierByID, "ID"},
		{database.OrderSupplierByName, "Nome"},
		{database.OrderSupplierByEmail, "Email"},
	}

	data.Table.OrderBy, err = strconv.Atoi(r.URL.Query().Get("orderBy"))
	if err != nil {
		data.Table.OrderBy = database.OrderSupplierByID
	}
	data.Table.OrderDesc = r.URL.Query().Get("orderDesc") == "true"

	data.Suppliers, err = appContext.Database(r).FindAllSuppliers(data.Table.OrderBy, data.Table.OrderDesc)
	if err != nil {
		ShowError(w, r, err)
		return
	}

	appContext.ExecuteTemplate(w, r, "suppliersTable", data)
}

func GetSupplier(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Supplier   database.Supplier
		NameInput  components.Input
		EmailInput components.Input
		IsNew      bool
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		data.IsNew = true
		data.Supplier = database.Supplier{}
	} else {
		supplier, err := appContext.Database(r).FindSupplier(id)
		if err != nil {
			ShowError(w, r, err)
			return
		}
		data.Supplier = supplier
	}

	data.NameInput = components.Input{"Nome", keySupplierName, "text", data.Supplier.Name}
	data.EmailInput = components.Input{"Email", keySupplierEmail, "text", data.Supplier.Email}

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
