package database

import (
	"fmt"

	"gorm.io/gorm/clause"
)

type Supplier struct {
	ID    int    `gorm:"column:id;primaryKey"`
	Email string `gorm:"column:email;size:255"`
	Name  string `gorm:"column:nome;size:255"`
}

func (Supplier) TableName() string { return "fornitori" }

const (
	OrderSupplierByID = iota
	OrderSupplierByEmail
	OrderSupplierByName
)

func (db *GormDB) FindAllSuppliers(orderBy int, orderDesc bool) ([]Supplier, error) {
	var orderByString string
	var suppliers []Supplier

	switch orderBy {
	case OrderSupplierByID:
		orderByString = "id"
	case OrderSupplierByEmail:
		orderByString = "email"
	case OrderSupplierByName:
		orderByString = "nome"
	default:
		return nil, fmt.Errorf("invalid orderBy value: %d", orderBy)
	}

	err := db.conn.Preload(clause.Associations).Order(clause.OrderByColumn{
		Column: clause.Column{Name: orderByString},
		Desc:   orderDesc,
	}).Find(&suppliers).Error
	return suppliers, err
}

func (db *GormDB) CreateSupplier(s Supplier) error {
	return db.conn.Create(&s).Error
}

func (db *GormDB) UpdateSupplier(s Supplier) error {
	return db.conn.Save(&s).Error
}

func (db *GormDB) DeleteSupplier(id int) error {
	return db.conn.Delete(&Supplier{}, id).Error
}
