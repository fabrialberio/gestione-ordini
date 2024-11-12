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
	start, _ := time.Parse(dateFormat, r.FormValue(keyOrderSelectionStart))
	end, _ := time.Parse(dateFormat, r.FormValue(keyOrderSelectionEnd))
	supplierId, _ := strconv.Atoi(r.FormValue(keyOrderSelectionSupplierID))
	allSuppliers := supplierId == 0

	orders, err := getFilteredOrders(r, start, end, supplierId)
	if err != nil {
		ShowDatabaseQueryError(w, r, err)
		return
	}

	// TODO: Move this to exporters
	filename := "ordini_" + start.Format("2006-01-02") + "_" + end.Format("2006-01-02")
	if !allSuppliers {
		supplier, _ := appContext.Database(r).FindSupplier(supplierId)
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
	start, _ := time.Parse(dateFormat, r.FormValue(keyOrderSelectionStart))
	end, _ := time.Parse(dateFormat, r.FormValue(keyOrderSelectionEnd))
	supplierId, _ := strconv.Atoi(r.FormValue(keyOrderSelectionSupplierID))

	orders, err := getFilteredOrders(r, start, end, supplierId)
	if err != nil {
		logError(r, err)
		return
	}

	w.Write([]byte(strconv.Itoa(len(orders))))
}

// TODO: Add parseOrderSelection function

func getFilteredOrders(r *http.Request, start time.Time, end time.Time, supplierId int) ([]database.Order, error) {
	orders, err := appContext.Database(r).FindAllOrdersWithExpiresAtBetween(start, end)
	if err != nil {
		return nil, err
	}

	var filteredOrders []database.Order
	for _, o := range orders {
		if o.Product.SupplierID == supplierId || supplierId == 0 {
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
