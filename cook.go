package main

import (
	"gestione-ordini/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

func HandleGetCook(w http.ResponseWriter, r *http.Request) {
	GetRequestContext(r).Templ.ExecuteTemplate(w, "cook.html", nil)
}

func HandleGetCookOrdersList(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Orders []database.Order
	}

	user, err := GetAuthenticatedUser(r)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.Orders, err = GetRequestContext(r).DB.FindAllOrdersWithUserID(user.ID)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	GetRequestContext(r).Templ.ExecuteTemplate(w, "ordersList.html", data)
}

func HandleGetCookOrder(w http.ResponseWriter, r *http.Request) {
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
		order, err := GetRequestContext(r).DB.FindOrder(id)
		if err != nil {
			HandleError(w, r, err)
			return
		}

		data.Order = order
	}

	user, err := GetAuthenticatedUser(r)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	data.UserID = user.ID

	data.Products, err = GetRequestContext(r).DB.FindAllProducts()
	if err != nil {
		HandleError(w, r, err)
		return
	}
	log.Println(data.Products)

	GetRequestContext(r).Templ.ExecuteTemplate(w, "order.html", data)
}

func HandlePostCookOrder(w http.ResponseWriter, r *http.Request) {
	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")

	productId, _ := strconv.Atoi(r.FormValue("productId"))
	userId, _ := strconv.Atoi(r.FormValue("userId"))
	amount, _ := strconv.Atoi(r.FormValue("amount"))
	requestedAt, _ := time.Parse("2006-01-02", r.FormValue("requestedAt"))

	log.Println(requestedAt)

	if isNew {
		err := GetRequestContext(r).DB.CreateOrder(database.Order{
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
			err = GetRequestContext(r).DB.DeleteOrder(id)
			if err != nil {
				HandleError(w, r, err)
				return
			}
		} else {
			err = GetRequestContext(r).DB.UpdateOrder(database.Order{
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
