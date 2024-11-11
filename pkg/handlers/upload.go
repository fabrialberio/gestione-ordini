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

type uploadForm struct {
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
		KeepIdsSelect components.Select
	}{
		Sidebar: currentSidebar(4, true),
		KeepIdsSelect: components.Select{
			Name:  "keepIds",
			Label: "Mantieni ID",
			Options: []components.SelectOption{
				{Value: 1, Text: "Sì"},
				{Value: 0, Text: "No"},
			},
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

	products, err := files.ImportProductsFromCSV(form.CSVFile, form.KeepIds)
	if err != nil {
		logError(r, err)
		appContext.ExecuteTemplate(w, r, "uploadError", err.Error())
		return
	}

	rows := [][]string{}
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

	headings := []components.TableHeading{
		{Name: "ID Tipologia"},
		{Name: "ID Fornitore"},
		{Name: "ID Unità di misura"},
		{Name: "Descrizione"},
		{Name: "Codice"},
	}

	if form.KeepIds {
		headings = append([]components.TableHeading{{Name: "ID"}}, headings...)
	}

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

	products, err := files.ImportProductsFromCSV(form.CSVFile, form.KeepIds)
	if err != nil {
		ShowItemInvalidFormError(w, r, err)
		return
	}

	err = appContext.Database(r).CreateAllProducts(products)
	if err != nil {
		ShowDatabaseQueryError(w, r, err)
		return
	}

	http.Redirect(w, r, DestProducts, http.StatusSeeOther)
}

func parseUploadForm(r *http.Request) (uploadForm, error) {
	keepIds := r.FormValue("keepIds") == "1"

	f, _, err := r.FormFile("csvFile")
	if err != nil {
		return uploadForm{}, err
	}

	return uploadForm{
		KeepIds: keepIds,
		CSVFile: f,
	}, nil
}
