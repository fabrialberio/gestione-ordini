package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"strconv"
	"time"
)

func GetChefOrder(w http.ResponseWriter, r *http.Request) {
	user, err := appContext.AuthenticatedUser(r)
	if err != nil {
		LogoutError(w, r, err)
		return
	}

	isNew := false
	defaultOrder := database.Order{
		Amount:    1,
		ExpiresAt: time.Now(),
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		isNew = true
	} else {
		defaultOrder, err = appContext.Database(r).FindOrder(id)
		if err != nil {
			ShowItemNotAllowedError(w, r, err)
			return
		}

		if defaultOrder.UserID != user.ID {
			ShowItemNotAllowedError(w, r, err)
			return
		}
	}

	productTypes, err := appContext.Database(r).FindAllProductTypes()
	if err != nil {
		ShowDatabaseQueryError(w, r, err)
		return
	}

	data := struct {
		IsNew          bool
		Order          database.Order
		ProductInput   components.ProductInput
		AmountInputURL string
		AmountInput    components.Input
		ExpiresAtInput components.Input
		UserID         int
	}{
		IsNew: isNew,
		Order: defaultOrder,
		ProductInput: components.ProductInput{
			InitialProductName: defaultOrder.Product.Name,
			ProductSelectName:  keyOrderProductID,
			ProductSearchURL:   DestProductSearch,
			SearchInputName:    keyProductSearchQuery,
			ProductTypesName:   keyProductSearchProductTypes,
			ProductTypes:       productTypes,
		},
		AmountInputURL: DestOrderAmountInput,
		ExpiresAtInput: components.Input{
			Label:        "Data di consegna richiesta",
			Name:         keyOrderRequestedAt,
			Type:         "date",
			DefaultValue: defaultOrder.ExpiresAt.Format(dateFormat),
		},
		UserID: user.ID,
	}

	appContext.ExecuteTemplate(w, r, "chefOrder.html", data)
}

func PostChefOrder(w http.ResponseWriter, r *http.Request) {
	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")

	order, err := parseOrderFromForm(r)
	if err != nil {
		ShowItemInvalidFormError(w, r, err)
		return
	}

	if isNew {
		err := appContext.Database(r).CreateOrder(order)
		if err != nil {
			ShowDatabaseQueryError(w, r, err)
			return
		}
	} else {
		if delete {
			err = appContext.Database(r).DeleteOrder(order.ID)
			if err != nil {
				ShowItemNotDeletableError(w, r, err)
				return
			}
		} else {
			err = appContext.Database(r).UpdateOrder(order)
			if err != nil {
				ShowDatabaseQueryError(w, r, err)
				return
			}
		}
	}

	http.Redirect(w, r, DestChef, http.StatusSeeOther)
}

func parseOrderFromForm(r *http.Request) (database.Order, error) {
	order := database.Order{}

	id, err := strconv.Atoi(r.FormValue(keyOrderID))
	if err == nil {
		order.ID = id
	}

	userId, err := strconv.Atoi(r.FormValue(keyOrderUserID))
	if err != nil {
		// If no userId value is found, use authenticated user ID
		user, err := appContext.AuthenticatedUser(r)
		if err != nil {
			return order, err
		}

		order.UserID = user.ID
	} else {
		order.UserID = userId
	}

	order.ProductID, err = strconv.Atoi(r.FormValue(keyOrderProductID))
	if err != nil {
		return order, err
	}

	order.Amount, err = strconv.Atoi(r.FormValue(keyOrderAmount))
	if err != nil {
		return order, err
	}

	order.ExpiresAt, err = time.Parse(dateFormat, r.FormValue(keyOrderRequestedAt))
	if err != nil {
		return order, err
	}

	return order, nil
}
