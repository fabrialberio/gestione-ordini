package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"time"
)

func GetChef(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ProductAmountInput components.ProductAmountInput
		ExpiresAtInput     components.Input
	}

	defaultOrder := database.Order{
		Amount:    1,
		ExpiresAt: time.Now(),
	}

	data.ExpiresAtInput = components.Input{"Scadenza", keyOrderRequestedAt, "date", defaultOrder.ExpiresAt.Format(dateFormat)}

	products, err := appContext.Database(r).FindAllProducts(database.OrderProductByID, true)
	if err != nil {
		ShowError(w, r, err)
		return
	}

	data.ProductAmountInput = components.ProductAmountInput{keyOrderProductID, products, defaultOrder.ProductID, keyOrderAmount, defaultOrder.Amount}

	appContext.ExecuteTemplate(w, r, "chef.html", data)
}
