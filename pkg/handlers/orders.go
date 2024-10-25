package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"strconv"
	"time"
)

const dateFormat = "2006-01-02"

func GetChefOrdersList(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Orders []database.Order
	}

	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.Orders, err = appContext.Database(r).FindAllOrdersWithUserID(user.ID)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	appContext.ExecuteTemplate(w, r, "chefOrdersList.html", data)
}

func GetChefOrder(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Order              database.Order
		ProductAmountInput components.ProductAmountInput
		ExpiresAtInput     components.Input
		UserID             int
		IsNew              bool
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		data.IsNew = true
		data.Order = database.Order{
			Amount:    1,
			ExpiresAt: time.Now(),
		}
	} else {
		order, err := appContext.Database(r).FindOrder(id)
		if err != nil {
			HandleError(w, r, err)
			return
		}

		data.Order = order
	}

	data.ExpiresAtInput = components.Input{"Scadenza", keyOrderRequestedAt, "date", data.Order.ExpiresAt.Format(dateFormat)}

	products, err := appContext.Database(r).FindAllProducts(database.OrderProductByID, true)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.ProductAmountInput = components.ProductAmountInput{keyOrderProductID, products, data.Order.ProductID, keyOrderAmount, data.Order.Amount}

	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	data.UserID = user.ID

	appContext.ExecuteTemplate(w, r, "chefOrder.html", data)
}

func PostChefOrder(w http.ResponseWriter, r *http.Request) {
	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")
	user, err := appContext.AuthenticatedUser(r)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	productId, _ := strconv.Atoi(r.FormValue(keyOrderProductID))
	amount, _ := strconv.Atoi(r.FormValue(keyOrderAmount))
	requestedAt, _ := time.Parse(dateFormat, r.FormValue(keyOrderRequestedAt))

	if isNew {
		err := appContext.Database(r).CreateOrder(database.Order{
			ProductID: productId,
			UserID:    user.ID,
			Amount:    amount,
			ExpiresAt: requestedAt,
		})
		if err != nil {
			HandleError(w, r, err)
			return
		}
	} else {
		id, err := strconv.Atoi(r.FormValue(keyOrderID))
		if err != nil {
			HandleError(w, r, err)
			return
		}

		if delete {
			err = appContext.Database(r).DeleteOrder(id)
			if err != nil {
				HandleError(w, r, err)
				return
			}
		} else {
			err = appContext.Database(r).UpdateOrder(database.Order{
				ID:        id,
				ProductID: productId,
				UserID:    user.ID,
				Amount:    amount,
				ExpiresAt: requestedAt,
			})
			if err != nil {
				HandleError(w, r, err)
				return
			}
		}
	}

	http.Redirect(w, r, DestChef, http.StatusSeeOther)
}

func GetAllOrdersTable(w http.ResponseWriter, r *http.Request) {
	var err error
	var data components.OrdersTable

	data.Headings = []components.TableHeading{
		{Name: "ID"},
		{Name: "Prodotto"},
		{Name: "Utente"},
		{Name: "Quantità"},
		{Name: "Scadenza"},
	}

	data.Orders, err = appContext.Database(r).FindAllOrders()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	appContext.ExecuteTemplate(w, r, "allOrdersTable", data)
}
