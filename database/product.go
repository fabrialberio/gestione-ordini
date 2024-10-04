package database

type ProductType struct {
	ID   int    `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:nome;size:255"`
}

func (ProductType) TableName() string { return "tipologie_prodotto" }

type Supplier struct {
	ID int `gorm:"column:id;primaryKey"`
}

func (Supplier) TableName() string { return "fornitori" }

type UnitOfMeasure struct {
	ID     int    `gorm:"column:id;primaryKey"`
	Symbol string `gorm:"column:simbolo;size:10"`
}

func (UnitOfMeasure) TableName() string { return "unita_di_misura" }

type Product struct {
	ID              int    `gorm:"column:id;primaryKey"`
	ProductTypeID   int    `gorm:"column:id_tipologia"`
	SupplierID      int    `gorm:"column:id_fornitore"`
	UnitOfMeasureID int    `gorm:"column:id_unita_di_misura"`
	Name            string `gorm:"column:nome"`
}

func (Product) TableName() string { return "prodotti" }
