package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

func PostOrderSelection(w http.ResponseWriter, r *http.Request) {
	start, _ := time.Parse(dateFormat, r.FormValue(keyOrderSelectionStart))
	end, _ := time.Parse(dateFormat, r.FormValue(keyOrderSelectionEnd))
	supplierId, _ := strconv.Atoi(r.FormValue(keyOrderSelectionSupplierID))

	orders, err := appContext.Database(r).FindAllOrdersWithExpiresAtBetween(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var filteredOrders []database.Order
	for _, o := range orders {
		if o.Product.SupplierID == supplierId || supplierId == 0 {
			filteredOrders = append(filteredOrders, o)
		}
	}

	log.Println("Selected:\n", len(filteredOrders))

	// TODO: Actually download a list / csv file

	http.Redirect(w, r, DestAllOrders, http.StatusSeeOther)
}
