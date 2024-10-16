package main

import (
	"gestione-ordini/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

func cook(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, checkRole(r, database.RoleIDCook)
}

func cookOrdersList(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if err := checkPerm(r, database.PermIDEditOwnOrder); err != nil {
		return nil, err
	}

	var data struct {
		Orders []database.Order
	}

	claims, err := getSessionCookie(r)
	if err != nil {
		return nil, err
	}

	data.Orders, err = db.FindAllOrdersWithUserID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func cookOrder(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if err := checkPerm(r, database.PermIDEditOwnOrder); err != nil {
		return nil, err
	}

	var data struct {
		Order    database.Order
		Products []database.Product
		UserID   int
		IsNew    bool
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		data.IsNew = true
		data.Order = database.Order{
			Amount: 1,
		}
	} else {
		order, err := db.FindOrder(id)
		if err != nil {
			return nil, err
		}

		data.Order = order
	}

	claims, err := getSessionCookie(r)
	if err != nil {
		return nil, err
	}
	data.UserID = claims.UserID

	data.Products, err = db.FindAllProducts()
	if err != nil {
		return nil, err
	}
	log.Println(data.Products)

	return data, nil
}

func cookOrderEdit(w http.ResponseWriter, r *http.Request) error {
	if err := checkPerm(r, database.PermIDEditOwnOrder); err != nil {
		return err
	}

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
			return err
		}
	} else {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		if delete {
			err = db.DeleteOrder(id)
			if err != nil {
				return err
			}
		} else {
			err = db.UpdateOrder(database.Order{
				ID:        id,
				ProductID: productId,
				UserID:    userId,
				Amount:    amount,
			})
			if err != nil {
				return err
			}
		}
	}

	http.Redirect(w, r, "/cook", http.StatusSeeOther)
	return nil
}
