package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"strconv"
	"time"
)

const dateFormat = "2006-01-02"

var weekdayNames = map[time.Weekday]string{
	time.Monday:    "lun",
	time.Tuesday:   "mar",
	time.Wednesday: "mer",
	time.Thursday:  "gio",
	time.Friday:    "ven",
	time.Saturday:  "sab",
	time.Sunday:    "dom",
}

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

func GetChefOrdersView(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	user, err := appContext.AuthenticatedUser(r)
	if err != nil {
		LogoutError(w, r, err)
		return
	}

	orders, err := appContext.Database(r).FindAllOrdersWithUserID(user.ID)
	if err != nil {
		logError(r, err)
	}

	data := calculateOrdersView(offset, orders)
	data.OrdersViewURL = DestChefOrdersView
	data.OrdersURL = DestChefOrders

	appContext.ExecuteTemplate(w, r, "ordersView", data)
}

func GetAllOrdersView(w http.ResponseWriter, r *http.Request) {
	orders, err := appContext.Database(r).FindAllOrders()
	if err != nil {
		logError(r, err)
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	data := calculateOrdersView(offset, orders)
	data.OrdersViewURL = DestAllOrdersView

	appContext.ExecuteTemplate(w, r, "ordersView", data)
}

func calculateOrdersView(offset int, orders []database.Order) components.OrdersView {
	start := currentWeekStart().Add(time.Hour * 24 * 7 * time.Duration(offset))
	days := makeOrdersViewDays(start, orders)

	return components.OrdersView{
		WeekTitle:  "Settimana del " + start.Format("02/01/2006"),
		NextOffset: offset + 1,
		PrevOffset: offset - 1,
		Days:       days,
	}
}

func currentWeekStart() time.Time {
	daysFromMonday := time.Duration(time.Now().Weekday() - 1)
	return time.Now().Add(time.Hour * 24 * -daysFromMonday)
}

func makeOrdersViewDays(start time.Time, orders []database.Order) []components.OrdersViewDay {
	ordersByDay := map[string][]database.Order{}
	days := []components.OrdersViewDay{}
	keyFormat := "2006-01-02"

	for _, o := range orders {
		ordersByDay[o.ExpiresAt.Format(keyFormat)] = append(ordersByDay[o.ExpiresAt.Format(keyFormat)], o)
	}

	for i := 0; i < 7; i++ {
		t := start.Add(time.Hour * 24 * time.Duration(i))

		days = append(days, components.OrdersViewDay{
			Heading: weekdayNames[t.Weekday()] + " " + t.Format("2"),
			Orders:  ordersByDay[t.Format(keyFormat)],
			IsPast:  time.Until(t.Add(time.Hour*24)) < 0,
		})
	}

	return days
}
