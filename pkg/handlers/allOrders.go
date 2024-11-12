package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"gestione-ordini/pkg/files"
	"net/http"
	"strconv"
	"time"
)

type orderSelection struct {
	Start        time.Time
	End          time.Time
	SupplierID   int
	AllSuppliers bool
}

func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	defaultStart := currentWeekStart()
	weekDuration := time.Hour * 24 * 6

	suppliers, err := appContext.Database(r).FindAllSuppliers(database.OrderSupplierByID, true)
	if err != nil {
		ShowDatabaseQueryError(w, r, err)
		return
	}

	supplierOptions := []components.SelectOption{{
		Value: 0,
		Text:  "Tutti i fornitori",
	}}
	for _, s := range suppliers {
		supplierOptions = append(supplierOptions, components.SelectOption{Value: s.ID, Text: s.Name})
	}

	user, err := appContext.AuthenticatedUser(r)
	if err != nil {
		LogoutError(w, r, err)
		return
	}

	data := struct {
		Sidebar        []components.SidebarDest
		StartDateInput components.Input
		EndDateInput   components.Input
		SupplierSelect components.Select
	}{
		Sidebar: currentSidebar(sidebarIndexAllOrders, user.RoleID == database.RoleIDAdministrator),
		StartDateInput: components.Input{
			Label:        "Da",
			Name:         keyOrderSelectionStart,
			DefaultValue: defaultStart.Format(dateFormat),
		},
		EndDateInput: components.Input{
			Label:        "A",
			Name:         keyOrderSelectionEnd,
			DefaultValue: defaultStart.Add(weekDuration).Format(dateFormat),
		},
		SupplierSelect: components.Select{
			Label:    "Fornitore",
			Name:     keyOrderSelectionSupplierID,
			Selected: 0,
			Options:  supplierOptions,
		},
	}

	appContext.ExecuteTemplate(w, r, "allOrders.html", data)
}

func PostOrderSelection(w http.ResponseWriter, r *http.Request) {
	selection, err := parseOrderSelectionForm(r)
	if err != nil {
		ShowItemInvalidFormError(w, r, err)
		return
	}

	orders, err := getFilteredOrders(r, selection)
	if err != nil {
		ShowDatabaseQueryError(w, r, err)
		return
	}

	filename := "ordini_" + selection.Start.Format("2006-01-02") + "_" + selection.End.Format("2006-01-02")
	if !selection.AllSuppliers {
		supplier, _ := appContext.Database(r).FindSupplier(selection.SupplierID)
		filename += "_" + supplier.Name
	}

	if r.Form.Has("csv") {
		csv := files.ExportToCSV(orders)
		downloadFile(w, filename+".csv", "text/csv", csv)
	} else if r.Form.Has("list") {
		list := files.ExportToList(orders)
		downloadFile(w, filename+".txt", "text/plain", list)
	}
}

func PostOrderSelectionCount(w http.ResponseWriter, r *http.Request) {
	selection, _ := parseOrderSelectionForm(r)

	orders, err := getFilteredOrders(r, selection)
	if err != nil {
		logError(r, err)
		return
	}

	w.Write([]byte(strconv.Itoa(len(orders))))
}

func parseOrderSelectionForm(r *http.Request) (sel orderSelection, err error) {
	sel.Start, err = time.Parse(dateFormat, r.FormValue(keyOrderSelectionStart))
	if err != nil {
		return
	}
	sel.End, err = time.Parse(dateFormat, r.FormValue(keyOrderSelectionEnd))
	if err != nil {
		return
	}
	sel.SupplierID, err = strconv.Atoi(r.FormValue(keyOrderSelectionSupplierID))
	if err != nil {
		return
	}
	sel.AllSuppliers = sel.SupplierID == 0

	return
}

func getFilteredOrders(r *http.Request, sel orderSelection) ([]database.Order, error) {
	orders, err := appContext.Database(r).FindAllOrdersWithExpiresAtBetween(sel.Start, sel.End)
	if err != nil {
		return nil, err
	}

	var filteredOrders []database.Order
	for _, o := range orders {
		if o.Product.SupplierID == sel.SupplierID || sel.AllSuppliers {
			filteredOrders = append(filteredOrders, o)
		}
	}

	return filteredOrders, nil
}

func downloadFile(w http.ResponseWriter, filename string, contentType string, content []byte) {
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", contentType)
	w.Write(content)
}
