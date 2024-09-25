package database

import "time"

type Order struct {
	ID          int
	ProductID   int
	UserID      int
	Amount      int
	RequestedAt time.Time
}
