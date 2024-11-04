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
	// TODO: Check order is made by current user
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
			ShowError(w, r, err)
			return
		}

		data.Order = order
	}

	data.ExpiresAtInput = components.Input{"Scadenza", keyOrderRequestedAt, "date", data.Order.ExpiresAt.Format(dateFormat)}

	products, err := appContext.Database(r).FindAllProducts(database.OrderProductByID, true)
	if err != nil {
		ShowError(w, r, err)
		return
	}

	data.ProductAmountInput = components.ProductAmountInput{keyOrderProductID, products, data.Order.ProductID, keyOrderAmount, data.Order.Amount}

	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		ShowError(w, r, err)
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
		ShowError(w, r, err)
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
			ShowError(w, r, err)
			return
		}
	} else {
		id, err := strconv.Atoi(r.FormValue(keyOrderID))
		if err != nil {
			ShowError(w, r, err)
			return
		}

		if delete {
			err = appContext.Database(r).DeleteOrder(id)
			if err != nil {
				ShowError(w, r, err)
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
				ShowError(w, r, err)
				return
			}
		}
	}

	http.Redirect(w, r, DestChef, http.StatusSeeOther)
}

func GetChefOrdersView(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	user, err := appContext.AuthenticatedUser(r)
	if err != nil {
		ShowError(w, r, err)
		return
	}

	orders, err := appContext.Database(r).FindAllOrdersWithUserID(user.ID)
	if err != nil {
		ShowError(w, r, err)
		return
	}

	data := calculateOrdersView(offset, orders)
	data.OrdersViewURL = DestChefOrdersView
	data.OrdersURL = DestChefOrders

	appContext.ExecuteTemplate(w, r, "ordersView", data)
}

func GetAllOrdersView(w http.ResponseWriter, r *http.Request) {
	orders, err := appContext.Database(r).FindAllOrders()
	if err != nil {
		ShowError(w, r, err)
		return
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
