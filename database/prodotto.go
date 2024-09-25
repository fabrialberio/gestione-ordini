package database

type ProductType struct {
	ID   int
	Name string
}

type Supplier struct {
	ID int
}

type UnitOfMeasure struct {
	ID     int
	Symbol string
}

type Product struct {
	ID              int
	ProductTypeID   int
	SupplierID      int
	UnitOfMeasureID int
	Name            string
}
