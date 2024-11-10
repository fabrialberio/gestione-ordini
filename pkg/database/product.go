package database

import (
	"fmt"

	"gorm.io/gorm/clause"
)

type ProductType struct {
	ID   int    `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:nome;size:255"`
}

func (ProductType) TableName() string { return "tipologie_prodotto" }

type UnitOfMeasure struct {
	ID     int    `gorm:"column:id;primaryKey"`
	Symbol string `gorm:"column:simbolo;size:10"`
}

func (UnitOfMeasure) TableName() string { return "unita_di_misura" }

type Product struct {
	ID              int           `gorm:"column:id;primaryKey"`
	ProductTypeID   int           `gorm:"column:id_tipologia"`
	ProductType     ProductType   `gorm:"-:all"`
	SupplierID      int           `gorm:"column:id_fornitore"`
	Supplier        Supplier      `gorm:"-:all"`
	UnitOfMeasureID int           `gorm:"column:id_unita_di_misura"`
	UnitOfMeasure   UnitOfMeasure `gorm:"-:all"`
	Description     string        `gorm:"column:descrizione"`
	Code            string        `gorm:"column:codice"`
}

func (Product) TableName() string { return "prodotti" }

const (
	OrderProductByID = iota
	OrderProductByProductType
	OrderProductBySupplier
	OrderProductByUnitOfMeasure
	OrderProductByDescription
	OrderProductByCode
)

func (db *GormDB) completeProduct(product Product) Product {
	product.ProductType, _ = db.FindProductType(product.ProductTypeID)
	product.Supplier, _ = db.FindSupplier(product.SupplierID)
	product.UnitOfMeasure, _ = db.FindUnitOfMeasure(product.UnitOfMeasureID)
	return product
}

func (db *GormDB) FindProductType(id int) (ProductType, error) {
	var productType ProductType

	err := db.conn.Take(&productType, id).Error
	return productType, err
}

func (db *GormDB) FindAllProductTypes() ([]ProductType, error) {
	var productTypes []ProductType

	err := db.conn.Find(&productTypes).Error
	return productTypes, err
}

func (db *GormDB) FindUnitOfMeasure(id int) (UnitOfMeasure, error) {
	var unitOfMeasure UnitOfMeasure

	err := db.conn.Take(&unitOfMeasure, id).Error
	return unitOfMeasure, err
}

func (db *GormDB) FindAllUnitsOfMeasure() ([]UnitOfMeasure, error) {
	var unitsOfMeasure []UnitOfMeasure

	err := db.conn.Find(&unitsOfMeasure).Error
	return unitsOfMeasure, err
}

func (db *GormDB) FindAllProducts(orderBy int, orderDesc bool) ([]Product, error) {
	var orderByString string
	var products []Product

	switch orderBy {
	case OrderProductByID:
		orderByString = "id"
	case OrderProductByProductType:
		orderByString = "id_tipologia"
	case OrderProductBySupplier:
		orderByString = "id_fornitore"
	case OrderProductByUnitOfMeasure:
		orderByString = "id_unita_di_misura"
	case OrderProductByDescription:
		orderByString = "descrizione"
	case OrderProductByCode:
		orderByString = "codice"
	default:
		return nil, fmt.Errorf("invalid orderBy value: %d", orderBy)
	}

	err := db.conn.Order(clause.OrderByColumn{
		Column: clause.Column{Name: orderByString},
		Desc:   orderDesc,
	}).Find(&products).Error

	for i, p := range products {
		products[i] = db.completeProduct(p)
	}
	return products, err
}

func (db *GormDB) FindProduct(id int) (Product, error) {
	var product Product

	err := db.conn.Take(&product, id).Error
	return db.completeProduct(product), err
}

func (db *GormDB) CreateProduct(product Product) error {
	return db.conn.Create(&product).Error
}

func (db *GormDB) UpdateProduct(product Product) error {
	columns := []string{"id_tipologia", "id_fornitore", "id_unita_di_misura", "descrizione", "codice"}
	return db.conn.Model(&product).Select(columns).Updates(product).Error
}

func (db *GormDB) DeleteProduct(id int) error {
	return db.conn.Delete(&Product{}, id).Error
}
