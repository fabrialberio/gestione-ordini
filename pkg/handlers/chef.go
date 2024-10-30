package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"time"
)

func GetChef(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ProductAmountInput components.ProductAmountInput
		ExpiresAtInput     components.Input
		OrdersView         components.OrdersView
	}

	defaultOrder := database.Order{
		Amount:    1,
		ExpiresAt: time.Now(),
	}

	data.ExpiresAtInput = components.Input{"Scadenza", keyOrderRequestedAt, "date", defaultOrder.ExpiresAt.Format(dateFormat)}

	products, err := appContext.Database(r).FindAllProducts(database.OrderProductByID, true)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.ProductAmountInput = components.ProductAmountInput{keyOrderProductID, products, defaultOrder.ProductID, keyOrderAmount, defaultOrder.Amount}

	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	orders, err := appContext.Database(r).FindAllOrdersWithUserID(user.ID)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.OrdersView = composeOrdersView(time.Now(), DestChefOrders, orders)

	appContext.ExecuteTemplate(w, r, "chef.html", data)
}
