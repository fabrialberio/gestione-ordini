package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/database"
	"gestione-ordini/pkg/exporters"
	"net/http"
	"strconv"
	"time"
)

func PostOrderSelection(w http.ResponseWriter, r *http.Request) {
	start, _ := time.Parse(dateFormat, r.FormValue(keyOrderSelectionStart))
	end, _ := time.Parse(dateFormat, r.FormValue(keyOrderSelectionEnd))
	supplierId, _ := strconv.Atoi(r.FormValue(keyOrderSelectionSupplierID))
	allSuppliers := supplierId == 0

	orders, err := getFilteredOrders(r, start, end, supplierId)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	filename := "ordini_" + start.Format("2006-01-02") + "_" + end.Format("2006-01-02")
	if !allSuppliers {
		supplier, _ := appContext.Database(r).FindSupplier(supplierId)
		filename += "_" + supplier.Name
	}

	if r.Form.Has("csv") {
		csv := exporters.ExportToCSV(orders)
		downloadFile(w, filename+".csv", "text/csv", csv)
	} else if r.Form.Has("list") {
		list := exporters.ExportToList(orders)
		downloadFile(w, filename+".txt", "text/plain", list)
	}
}

func GetOrderSelectionCount(w http.ResponseWriter, r *http.Request) {
	start, _ := time.Parse(dateFormat, r.URL.Query().Get(keyOrderSelectionStart))
	end, _ := time.Parse(dateFormat, r.URL.Query().Get(keyOrderSelectionEnd))
	supplierId, _ := strconv.Atoi(r.URL.Query().Get(keyOrderSelectionSupplierID))

	orders, err := getFilteredOrders(r, start, end, supplierId)
	if err != nil {
		return
	}

	w.Write([]byte(strconv.Itoa(len(orders))))
}

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
