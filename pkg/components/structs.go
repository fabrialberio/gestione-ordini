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
	SearchDialog            ProductSearchDialog
	ProductSelectName       string
	SelectedProduct         int
	AmountInputName         string
	AmountInputDefaultValue int
}

type ProductSearchDialog struct {
	ProductSearchURL string
	SearchInputName  string
	ProductTypesName string
	ProductTypes     []database.ProductType
}

type ProductSearchResult struct {
	ProductSelectName string
	Product           database.Product
	IsSelected        bool
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
