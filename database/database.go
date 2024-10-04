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

	db.AutoMigrate(&ProductType{})
	db.AutoMigrate(&Supplier{})
	db.AutoMigrate(&UnitOfMeasure{})
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&Order{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Role{})

	return &Database{db}, nil
}

func (db *Database) Close() error {
	// TODO: Remove
	return nil
}
