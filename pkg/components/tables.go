package components

import "gestione-ordini/pkg/database"

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

type PreviewTable struct {
	TableURL    string
	MaxRowCount int
	Headings    []TableHeading
	Rows        [][]string
}

type TableHeading struct {
	Index int
	Name  string
}
