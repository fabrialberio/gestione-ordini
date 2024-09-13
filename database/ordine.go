package database

import "time"

type Ordine struct {
	ID            int
	IDProdotto    int
	IDUtente      int
	Quantita      int
	DataRichiesta time.Time
}
