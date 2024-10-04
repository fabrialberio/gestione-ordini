package database

import "time"

type Order struct {
	ID          int       `gorm:"column:id;primaryKey"`
	ProductID   int       `gorm:"column:id_prodotto"`
	UserID      int       `gorm:"column:id_utente"`
	Amount      int       `gorm:"column:quantita"`
	RequestedAt time.Time `gorm:"column:richiesto_il"`
}
