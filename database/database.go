package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Conn *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("mysql", "user:password@tcp(mysql-container:3306)/testdb")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	log.Println("Successfully connected to the datanase.")
	return &Database{db}, nil
}

func (db *Database) Close() error {
	return db.Conn.Close()
}
