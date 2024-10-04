package database

type ProductType struct {
	ID   int    `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:nome;size:255"`
}

type Supplier struct {
	ID int `gorm:"column:id;primaryKey"`
}

type UnitOfMeasure struct {
	ID     int    `gorm:"column:id;primaryKey"`
	Symbol string `gorm:"column:simbolo;size:10"`
}

type Product struct {
	ID              int    `gorm:"column:id;primaryKey"`
	ProductTypeID   int    `gorm:"column:id_tipologia"`
	SupplierID      int    `gorm:"column:id_fornitore"`
	UnitOfMeasureID int    `gorm:"column:id_unita_di_misura"`
	Name            string `gorm:"column:nome"`
}
