package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"time"
)

func GetChef(w http.ResponseWriter, r *http.Request) {

	data := struct {
		Order              database.Order
		ProductAmountInput components.ProductAmountInput
		RequestedAtInput   components.Input
	}{
		Order: database.Order{
			Amount:      1,
			RequestedAt: time.Now(),
		},
	}

	data.RequestedAtInput = components.Input{"Richiesto per", keyOrderRequestedAt, "date", data.Order.RequestedAt.Format(dateFormat)}

	products, err := appContext.Database(r).FindAllProducts(database.OrderProductByID, true)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.ProductAmountInput = components.ProductAmountInput{keyOrderProductID, products, data.Order.ProductID, keyOrderAmount, data.Order.Amount}

	appContext.ExecuteTemplate(w, r, "chef.html", data)
}
