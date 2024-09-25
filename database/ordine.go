package database

import "time"

type Ordine struct {
	ID          int
	IDProdotto  int
	IDUtente    int
	Quantita    int
	RichiestoIl time.Time
}
