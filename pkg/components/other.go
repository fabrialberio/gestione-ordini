package components

import "gestione-ordini/pkg/database"

type SidebarDest struct {
	DestURL     string
	FasIconName string
	Label       string
	Selected    bool
}

type OrdersView struct {
	OrdersURL     string
	OrdersViewURL string
	WeekTitle     string
	NextOffset    int
	PrevOffset    int
	Days          []OrdersViewDay
}

type OrdersViewDay struct {
	Heading string
	Orders  []database.Order
	IsPast  bool
}
