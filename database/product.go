package database

import "gorm.io/gorm/clause"

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
	ID              int           `gorm:"column:id;primaryKey"`
	ProductTypeID   int           `gorm:"column:id_tipologia"`
	ProductType     ProductType   `gorm:"foreignKey:ProductTypeID"`
	SupplierID      int           `gorm:"column:id_fornitore"`
	Supplier        Supplier      `gorm:"foreignKey:SupplierID"`
	UnitOfMeasureID int           `gorm:"column:id_unita_di_misura"`
	UnitOfMeasure   UnitOfMeasure `gorm:"foreignKey:UnitOfMeasureID"`
	Name            string        `gorm:"column:nome"`
}

func (Product) TableName() string { return "prodotti" }

func (db *Database) GetProductTypes() ([]ProductType, error) {
	var productTypes []ProductType

	err := db.conn.Preload(clause.Associations).Find(&productTypes).Error
	return productTypes, err
}

func (db *Database) GetSuppliers() ([]Supplier, error) {
	var suppliers []Supplier

	err := db.conn.Preload(clause.Associations).Find(&suppliers).Error
	return suppliers, err
}

func (db *Database) GetUnitsOfMeasure() ([]UnitOfMeasure, error) {
	var unitsOfMeasure []UnitOfMeasure

	err := db.conn.Preload(clause.Associations).Find(&unitsOfMeasure).Error
	return unitsOfMeasure, err
}

func (db *Database) GetProducts() ([]Product, error) {
	var products []Product

	err := db.conn.Preload(clause.Associations).Find(&products).Error
	return products, err
}

func (db *Database) GetProduct(id int) (Product, error) {
	var product Product

	err := db.conn.Take(&product, id).Error
	return product, err
}

func (db *Database) CreateProduct(product Product) error {
	return db.conn.Create(&product).Error
}

func (db *Database) UpdateProduct(product Product) error {
	columns := []string{"id_tipologia", "id_fornitore", "id_unita_di_misura", "nome"}
	return db.conn.Model(&product).Select(columns).Updates(product).Error
}

func (db *Database) DeleteProduct(id int) error {
	return db.conn.Delete(&Product{}, id).Error
}
