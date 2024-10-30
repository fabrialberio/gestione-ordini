package components

import "gestione-ordini/pkg/database"

type SidebarDest struct {
	DestURL     string
	FasIconName string
	Label       string
	Selected    bool
}

type Input struct {
	Label        string
	Name         string
	Type         string
	DefaultValue string
}

type ProductAmountInput struct {
	ProductSelectName       string
	Products                []database.Product
	SelectedProduct         int
	AmountInputName         string
	AmountInputDefaultValue int
}

type Select struct {
	Label    string
	Name     string
	Selected int
	Options  []SelectOption
}

type SelectOption struct {
	Value int
	Text  string
}

type Table struct {
	TableURL  string
	OrderBy   int
	OrderDesc bool
	Headings  []TableHeading
}

type ProductsTable struct {
	Table
	Products []database.Product
}

type UsersTable struct {
	Table
	Users []database.User
}

type SuppliersTable struct {
	Table
	Suppliers []database.Supplier
}

type TableHeading struct {
	Index int
	Name  string
}

type OrdersView struct {
	OrdersURL string
	WeekTitle string
	Days      []OrdersViewDay
}

type OrdersViewDay struct {
	Heading string
	Orders  []database.Order
}
