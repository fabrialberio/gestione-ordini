package main

import (
	"gestione-ordini/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

func HandleGetCook(w http.ResponseWriter, r *http.Request) {
	templ.ExecuteTemplate(w, "cook.html", nil)
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

	data.Orders, err = db.FindAllOrdersWithUserID(user.ID)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	templ.ExecuteTemplate(w, "ordersList.html", data)
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
		order, err := db.FindOrder(id)
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

	data.Products, err = db.FindAllProducts()
	if err != nil {
		HandleError(w, r, err)
		return
	}
	log.Println(data.Products)

	templ.ExecuteTemplate(w, "order.html", data)
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
		err := db.CreateOrder(database.Order{
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
			err = db.DeleteOrder(id)
			if err != nil {
				HandleError(w, r, err)
				return
			}
		} else {
			err = db.UpdateOrder(database.Order{
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
