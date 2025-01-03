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
	start := currentWeekStart().Add(time.Hour * 24 * 7 * time.Duration(offset))

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
	start := currentWeekStart().Add(time.Hour * 24 * 7 * time.Duration(offset))

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
