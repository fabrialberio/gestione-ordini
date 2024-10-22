package database

import "gorm.io/gorm/clause"

type Supplier struct {
	ID    int    `gorm:"column:id;primaryKey"`
	Email string `gorm:"column:email;size:255"`
	Name  string `gorm:"column:nome;size:255"`
}

func (Supplier) TableName() string { return "fornitori" }

func (db *GormDB) FindAllSuppliers() ([]Supplier, error) {
	var suppliers []Supplier

	err := db.conn.Preload(clause.Associations).Find(&suppliers).Error
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
