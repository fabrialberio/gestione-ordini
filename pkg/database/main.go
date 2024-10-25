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
	db, err := gorm.Open(mysql.Open(dsn+"?parseTime=true"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	return &GormDB{db}, nil
}

func (db *GormDB) Close() {}
