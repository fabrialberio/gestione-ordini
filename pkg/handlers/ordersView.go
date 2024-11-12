package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"strconv"
	"time"
)

var weekdayNames = map[time.Weekday]string{
	time.Monday:    "lun",
	time.Tuesday:   "mar",
	time.Wednesday: "mer",
	time.Thursday:  "gio",
	time.Friday:    "ven",
	time.Saturday:  "sab",
	time.Sunday:    "dom",
}

func GetOwnOrdersView(w http.ResponseWriter, r *http.Request) {
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
	data.OrdersViewURL = DestOwnOrdersView
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
	data.OrdersURL = DestOrders
	data.UsersURL = DestUsers

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
