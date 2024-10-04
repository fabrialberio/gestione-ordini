package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	conn *gorm.DB
}

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

func NewDatabase(dsn string) (*Database, error) {
	db, err := gorm.Open(mysql.Open(dsn+"?parseTime=true"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
	})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	db.Table("tipologie_prodotto").AutoMigrate(&ProductType{})
	db.Table("fornitori").AutoMigrate(&Supplier{})
	db.Table("unita_di_misura").AutoMigrate(&UnitOfMeasure{})
	db.Table("prodotti").AutoMigrate(&Product{})
	db.Table("ordini").AutoMigrate(&Order{})
	db.Table("utenti").AutoMigrate(&User{})
	db.Table("ruoli").AutoMigrate(&Role{})

	return &Database{db}, nil
}

func (db *Database) Close() error {
	// TODO: Remove
	return nil
}
