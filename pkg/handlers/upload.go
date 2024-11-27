package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"gestione-ordini/pkg/files"
	"mime/multipart"
	"net/http"
	"strconv"
)

const (
	tableProducts = 1
	tableUsers    = 2
)

type uploadForm struct {
	Table   int
	KeepIds bool
	CSVFile multipart.File
}

func GetUpload(w http.ResponseWriter, r *http.Request) {
	user, err := appContext.AuthenticatedUser(r)
	if err != nil {
		LogoutError(w, r, err)
		return
	} else if user.RoleID != database.RoleIDAdministrator {
		ShowItemNotAllowedError(w, r, auth.ErrInvalidRole)
		return
	}

	data := struct {
		Sidebar       []components.SidebarDest
		TableSelect   components.Select
		KeepIdsSelect components.Select
	}{
		Sidebar: currentSidebar(sidebarIndexUpload, true),
		TableSelect: components.Select{
			Name:  "table",
			Label: "Tabella",
			Options: []components.SelectOption{
				{Value: tableProducts, Text: "Prodotti"},
				{Value: tableUsers, Text: "Utenti"},
			},
		},
		KeepIdsSelect: components.Select{
			Name:  "keepIds",
			Label: "Comportamento",
			Options: []components.SelectOption{
				{Value: 0, Text: "Aggiungi"},
				{Value: 1, Text: "Aggiorna in base all'ID"},
			},
			Selected: 0,
		},
	}

	appContext.ExecuteTemplate(w, r, "upload.html", data)
}

func PostUploadPreview(w http.ResponseWriter, r *http.Request) {
	form, err := parseUploadForm(r)
	if err != nil {
		appContext.ExecuteTemplate(w, r, "infoCard", "Nessun file selezionato")
		return
	}

	var headings []components.TableHeading
	query := components.TableQuery{MaxRowCount: 10}

	switch form.Table {
	case tableProducts:
		headings = []components.TableHeading{
			{Name: "ID Tipologia"},
			{Name: "ID Fornitore"},
			{Name: "ID Unità di misura"},
			{Name: "Descrizione"},
			{Name: "Codice"},
		}
	case tableUsers:
		headings = []components.TableHeading{
			{Name: "ID Ruolo"},
			{Name: "Username"},
			{Name: "Password hash"},
			{Name: "Nome"},
			{Name: "Cognome"},
		}
	}

	if form.KeepIds {
		headings = append([]components.TableHeading{{Name: "ID"}}, headings...)
	}

	var table components.Table
	switch form.Table {
	case tableProducts:
		items, err := files.ImportProductsFromCSV(form.CSVFile, form.KeepIds)
		if err != nil {
			logError(r, err)
			appContext.ExecuteTemplate(w, r, "errorCard", err.Error())
			return
		}

		rowFunc := func(p database.Product) components.TableRow {
			return components.TableRow{
				Cells: []components.TableCell{
					{Value: strconv.Itoa(p.ProductTypeID)},
					{Value: strconv.Itoa(p.SupplierID)},
					{Value: strconv.Itoa(p.UnitOfMeasureID)},
					{Value: p.Description},
					{Value: p.Code},
				},
			}
		}

		table = components.ComposeTable(query, headings, rowFunc, items)

	case tableUsers:
		items, err := files.ImportUsersFromCSV(form.CSVFile, form.KeepIds)
		if err != nil {
			logError(r, err)
			appContext.ExecuteTemplate(w, r, "errorCard", err.Error())
			return
		}

		rowFunc := func(u database.User) components.TableRow {
			return components.TableRow{
				Cells: []components.TableCell{
					{Value: strconv.Itoa(u.RoleID)},
					{Value: u.Username},
					{Value: u.PasswordHash},
					{Value: u.Name},
					{Value: u.Surname},
				},
			}
		}

		table = components.ComposeTable(query, headings, rowFunc, items)
	}

	appContext.ExecuteTemplate(w, r, "previewTable", table)
}

func PostUpload(w http.ResponseWriter, r *http.Request) {
	form, err := parseUploadForm(r)
	if err != nil {
		ShowItemInvalidFormError(w, r, err)
		return
	}

	switch form.Table {
	case tableProducts:
		products, err := files.ImportProductsFromCSV(form.CSVFile, form.KeepIds)
		if err != nil {
			ShowItemInvalidFormError(w, r, err)
			return
		}

		if form.KeepIds {
			err = appContext.Database(r).UpdateAllProducts(products)
		} else {
			err = appContext.Database(r).CreateAllProducts(products)
		}
		if err != nil {
			ShowDatabaseQueryError(w, r, err)
			return
		}
	case tableUsers:
		users, err := files.ImportUsersFromCSV(form.CSVFile, form.KeepIds)
		if err != nil {
			ShowItemInvalidFormError(w, r, err)
			return
		}

		if form.KeepIds {
			err = appContext.Database(r).CreateAllUsers(users)
		} else {
			err = appContext.Database(r).UpdateAllUsers(users)
		}
		if err != nil {
			ShowDatabaseQueryError(w, r, err)
			return
		}
	default:
		ShowItemInvalidFormError(w, r, nil)
		return
	}

	http.Redirect(w, r, DestProducts, http.StatusSeeOther)
}

func parseUploadForm(r *http.Request) (uploadForm, error) {
	keepIds := r.FormValue("keepIds") == "1"

	table, err := strconv.Atoi(r.FormValue("table"))
	if err != nil {
		return uploadForm{}, err
	}

	f, _, err := r.FormFile("csvFile")
	if err != nil {
		return uploadForm{}, err
	}

	return uploadForm{
		Table:   table,
		KeepIds: keepIds,
		CSVFile: f,
	}, nil
}
