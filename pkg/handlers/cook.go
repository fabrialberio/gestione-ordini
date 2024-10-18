package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetCook(w http.ResponseWriter, r *http.Request) {
	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "cook.html", nil)
}

func GetCookOrdersList(w http.ResponseWriter, r *http.Request) {
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

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "ordersList.html", data)
}

func GetCookOrder(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Order    database.Order
		Products []database.Product
		UserID   int
		IsNew    bool
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		data.IsNew = true
		data.Order = database.Order{
			Amount: 1,
		}
	} else {
		order, err := appContext.FromRequest(r).DB.FindOrder(id)
		if err != nil {
			HandleError(w, r, err)
			return
		}

		data.Order = order
	}

	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	data.UserID = user.ID

	data.Products, err = appContext.FromRequest(r).DB.FindAllProducts()
	if err != nil {
		HandleError(w, r, err)
		return
	}
	log.Println(data.Products)

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "order.html", data)
}

func PostCookOrder(w http.ResponseWriter, r *http.Request) {
	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")

	productId, _ := strconv.Atoi(r.FormValue("productId"))
	userId, _ := strconv.Atoi(r.FormValue("userId"))
	amount, _ := strconv.Atoi(r.FormValue("amount"))
	requestedAt, _ := time.Parse("2006-01-02", r.FormValue("requestedAt"))

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
		id, err := strconv.Atoi(r.FormValue("id"))
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

	http.Redirect(w, r, "/cook", http.StatusSeeOther)
}
