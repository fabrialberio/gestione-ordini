package components

import "gestione-ordini/pkg/database"

type Input struct {
	Label        string
	Name         string
	Type         string
	DefaultValue string
}

type ProductInput struct {
	InitialProductDescription string
	ProductSelectName         string
	ProductSearchURL          string
	SearchInputName           string
	ProductTypesName          string
	ProductTypes              []database.ProductType
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
