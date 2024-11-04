package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"time"
)

func GetChef(w http.ResponseWriter, r *http.Request) {
	products, err := appContext.Database(r).FindAllProducts(database.OrderProductByID, true)
	if err != nil {
		ShowError(w, r, err)
		return
	}

	data := struct {
		ProductAmountInput components.ProductAmountInput
		ExpiresAtInput     components.Input
	}{
		ProductAmountInput: components.ProductAmountInput{
			ProductSelectName:       keyOrderProductID,
			Products:                products,
			SelectedProduct:         0,
			AmountInputName:         keyOrderAmount,
			AmountInputDefaultValue: 1,
		},
		ExpiresAtInput: components.Input{
			Label:        "Scadenza",
			Name:         keyOrderRequestedAt,
			Type:         "date",
			DefaultValue: time.Now().Format(dateFormat),
		},
	}

	appContext.ExecuteTemplate(w, r, "chef.html", data)
}
