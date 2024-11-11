package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"gestione-ordini/pkg/files"
	"net/http"
	"strconv"
)

func PostUploadPreview(w http.ResponseWriter, r *http.Request) {
	f, _, err := r.FormFile("csvFile")
	if err != nil {
		logError(r, err)
		return
	}

	products, err := files.ImportProductsFromCSV(f)
	if err != nil {
		logError(r, err)
		appContext.ExecuteTemplate(w, r, "uploadError", err.Error())
		return
	}

	rows := [][]string{}
	for _, p := range products {
		rows = append(rows, []string{
			strconv.Itoa(p.ID),
			strconv.Itoa(p.ProductTypeID),
			strconv.Itoa(p.SupplierID),
			strconv.Itoa(p.UnitOfMeasureID),
			p.Description,
			p.Code,
		})
	}

	data := components.PreviewTable{
		Table: components.Table{
			TableURL: DestUploadPreview,
			OrderBy:  -1,
			Headings: []components.TableHeading{
				{Name: "ID"},
				{Name: "ID Tipologia"},
				{Name: "ID Fornitore"},
				{Name: "ID Unità di misura"},
				{Name: "Descrizione"},
				{Name: "Codice"},
			},
		},
		Rows: rows,
	}

	appContext.ExecuteTemplate(w, r, "previewTable", data)
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
		Sidebar []components.SidebarDest
	}{
		Sidebar: currentSidebar(4, true),
	}

	appContext.ExecuteTemplate(w, r, "upload.html", data)
}

func PostUpload(w http.ResponseWriter, r *http.Request) {

}
