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
		appContext.ExecuteTemplate(w, r, "uploadPlaceholder", nil)
		return
	}

	var headings []components.TableHeading
	var rows [][]string

	switch form.Table {
	case tableProducts:
		products, err := files.ImportProductsFromCSV(form.CSVFile, form.KeepIds)
		if err != nil {
			logError(r, err)
			appContext.ExecuteTemplate(w, r, "uploadError", err.Error())
			return
		}

		headings, rows = composeProductsPreview(form, products)
	case tableUsers:
		users, err := files.ImportUsersFromCSV(form.CSVFile, form.KeepIds)
		if err != nil {
			logError(r, err)
			appContext.ExecuteTemplate(w, r, "uploadError", err.Error())
			return
		}

		headings, rows = composeUsersPreview(form, users)
	default:
		appContext.ExecuteTemplate(w, r, "uploadPlaceholder", nil)
		return
	}

	// TODO: Tables with 1 row not shown??

	data := components.PreviewTable{
		TableURL:    DestUploadPreview,
		MaxRowCount: 8,
		Headings:    headings,
		Rows:        rows,
	}

	appContext.ExecuteTemplate(w, r, "previewTable", data)
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

func composeProductsPreview(form uploadForm, products []database.Product) (headings []components.TableHeading, rows [][]string) {
	rows = [][]string{}
	for _, p := range products {
		row := []string{
			strconv.Itoa(p.ProductTypeID),
			strconv.Itoa(p.SupplierID),
			strconv.Itoa(p.UnitOfMeasureID),
			p.Description,
			p.Code,
		}

		if form.KeepIds {
			row = append([]string{strconv.Itoa(p.ID)}, row...)
		}

		rows = append(rows, row)
	}

	headings = []components.TableHeading{
		{Name: "ID Tipologia"},
		{Name: "ID Fornitore"},
		{Name: "ID Unità di misura"},
		{Name: "Descrizione"},
		{Name: "Codice"},
	}

	if form.KeepIds {
		headings = append([]components.TableHeading{{Name: "ID"}}, headings...)
	}

	return headings, rows
}

func composeUsersPreview(form uploadForm, users []database.User) (headings []components.TableHeading, rows [][]string) {
	rows = [][]string{}
	for _, u := range users {
		row := []string{
			strconv.Itoa(u.RoleID),
			u.Username,
			u.PasswordHash,
			u.Name,
			u.Surname,
		}

		if form.KeepIds {
			row = append([]string{strconv.Itoa(u.ID)}, row...)
		}

		rows = append(rows, row)
	}

	headings = []components.TableHeading{
		{Name: "ID Ruolo"},
		{Name: "Username"},
		{Name: "Password hash"},
		{Name: "Nome"},
		{Name: "Cognome"},
	}

	if form.KeepIds {
		headings = append([]components.TableHeading{{Name: "ID"}}, headings...)
	}

	return headings, rows
}
