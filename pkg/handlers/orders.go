package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

const dateFormat = "2006-01-02"

func GetOrdersList(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Orders []database.Order
	}

	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.Orders, err = appContext.FromRequest(r).DB.FindAllOrdersWithUserID(user.ID)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "chefOrdersList.html", data)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Order            database.Order
		AmountInput      components.Input
		RequestedAtInput components.Input
		ProductSelect    components.Select
		UserID           int
		IsNew            bool
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		data.IsNew = true
		data.Order = database.Order{
			Amount:      1,
			RequestedAt: time.Now(),
		}
	} else {
		order, err := appContext.FromRequest(r).DB.FindOrder(id)
		if err != nil {
			HandleError(w, r, err)
			return
		}

		data.Order = order
	}

	data.AmountInput = components.Input{"Quantità", keyOrderAmount, "number", strconv.Itoa(data.Order.Amount)}
	data.RequestedAtInput = components.Input{"Richiesto per", keyOrderRequestedAt, "date", data.Order.RequestedAt.Format(dateFormat)}

	products, err := appContext.FromRequest(r).DB.FindAllProducts(database.OrderProductByID, true)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.ProductSelect = components.Select{"Prodotto", keyOrderProductID, data.Order.ProductID, []components.SelectOption{}}
	for _, p := range products {
		data.ProductSelect.Options = append(data.ProductSelect.Options, components.SelectOption{p.ID, p.Name})
	}

	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	data.UserID = user.ID

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "chefOrder.html", data)
}

func PostOrder(w http.ResponseWriter, r *http.Request) {
	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")

	productId, _ := strconv.Atoi(r.FormValue(keyOrderProductID))
	userId, _ := strconv.Atoi(r.FormValue(keyUserID))
	amount, _ := strconv.Atoi(r.FormValue(keyOrderAmount))
	requestedAt, _ := time.Parse(dateFormat, r.FormValue(keyOrderRequestedAt))

	log.Println(requestedAt)

	if isNew {
		err := appContext.FromRequest(r).DB.CreateOrder(database.Order{
			ProductID:   productId,
			UserID:      userId,
			Amount:      amount,
			RequestedAt: requestedAt,
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
			err = appContext.FromRequest(r).DB.DeleteOrder(id)
			if err != nil {
				HandleError(w, r, err)
				return
			}
		} else {
			err = appContext.FromRequest(r).DB.UpdateOrder(database.Order{
				ID:        id,
				ProductID: productId,
				UserID:    userId,
				Amount:    amount,
			})
			if err != nil {
				HandleError(w, r, err)
				return
			}
		}
	}

	http.Redirect(w, r, DestChef, http.StatusSeeOther)
}
