package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type GormDB struct {
	conn *gorm.DB
}

func New(dsn string) (*GormDB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn + "?parseTime=true",
		DefaultStringSize: 255,
	}), &gorm.Config{})
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

	return &GormDB{db}, nil
}

func (db *GormDB) Close() {}
