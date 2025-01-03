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
	start := currentWeekStart().AddDate(0, 0, offset*7)

	user, err := appContext.AuthenticatedUser(r)
	if err != nil {
		LogoutError(w, r, err)
		return
	}

	orders, err := appContext.Database(r).
		FindAllOrdersWithUserIDAndExpiresAtBetween(user.ID, start, start.AddDate(0, 0, 6))
	if err != nil {
		logError(r, err)
	}

	data := calculateOrdersView(start, offset, orders)
	data.OrdersViewURL = DestOwnOrdersView
	data.OrdersURL = DestChefOrders

	appContext.ExecuteTemplate(w, r, "ordersView", data)
}

func GetAllOrdersView(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	start := currentWeekStart().AddDate(0, 0, offset*7)

	orders, err := appContext.Database(r).
		FindAllOrdersWithExpiresAtBetween(start, start.AddDate(0, 0, 6))
	if err != nil {
		logError(r, err)
	}

	data := calculateOrdersView(start, offset, orders)
	data.OrdersViewURL = DestAllOrdersView
	data.OrdersURL = DestOrders

	appContext.ExecuteTemplate(w, r, "ordersView", data)
}

func calculateOrdersView(start time.Time, offset int, orders []database.Order) components.OrdersView {
	days := makeOrdersViewDays(start, orders)

	return components.OrdersView{
		WeekTitle:  "Settimana del " + start.Format("02/01/2006"),
		NextOffset: offset + 1,
		PrevOffset: offset - 1,
		Days:       days,
	}
}

func currentWeekStart() time.Time {
	daysFromMonday := int(time.Now().Weekday())
	return time.Now().AddDate(0, 0, 1-daysFromMonday)
}

func makeOrdersViewDays(start time.Time, orders []database.Order) []components.OrdersViewDay {
	ordersByDay := map[string][]database.Order{}
	days := []components.OrdersViewDay{}
	keyFormat := "2006-01-02"

	for _, o := range orders {
		ordersByDay[o.ExpiresAt.Format(keyFormat)] = append(ordersByDay[o.ExpiresAt.Format(keyFormat)], o)
	}

	for i := 0; i < 7; i++ {
		t := start.AddDate(0, 0, i)

		days = append(days, components.OrdersViewDay{
			Heading: weekdayNames[t.Weekday()] + " " + t.Format("2"),
			Orders:  ordersByDay[t.Format(keyFormat)],
			IsPast:  time.Until(t.AddDate(0, 0, 1)) < 0,
		})
	}

	return days
}
