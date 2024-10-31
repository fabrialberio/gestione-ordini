package database

import (
	"strconv"
	"time"
)

type Order struct {
	ID               int       `gorm:"column:id;primaryKey"`
	ProductID        int       `gorm:"column:id_prodotto"`
	Product          Product   `gorm:"-:all"`
	UserID           int       `gorm:"column:id_utente"`
	User             User      `gorm:"-:all"`
	Amount           int       `gorm:"column:quantita"`
	AmountString     string    `gorm:"-:all"`
	ExpiresAt        time.Time `gorm:"column:richiesto_il"`
	ExpirationString string    `gorm:"-:all"`
}

func (Order) TableName() string { return "ordini" }

func (db *GormDB) completeOrder(order Order) Order {
	order.Product, _ = db.FindProduct(order.ProductID)
	order.User, _ = db.FindUser(order.UserID)

	order.AmountString = strconv.Itoa(order.Amount) + " " + order.Product.UnitOfMeasure.Symbol

	expiresInHours := time.Until(order.ExpiresAt).Hours()

	if expiresInHours < 0 {
		order.ExpirationString = "scaduto"
	} else if expiresInHours/24 < 1 {
		order.ExpirationString = "oggi"
	} else if expiresInHours/24 < 2 {
		order.ExpirationString = "domani"
	} else {
		order.ExpirationString = strconv.Itoa(int(expiresInHours/24)) + " giorni"
	}

	return order
}

func (db *GormDB) FindOrder(id int) (Order, error) {
	var order Order

	err := db.conn.Take(&order, id).Error
	return db.completeOrder(order), err
}

func (db *GormDB) FindAllOrders() ([]Order, error) {
	var orders []Order

	err := db.conn.Find(&orders).Error
	for i, o := range orders {
		orders[i] = db.completeOrder(o)
	}
	return orders, err
}

func (db *GormDB) FindAllOrdersWithUserID(userId int) ([]Order, error) {
	var orders []Order

	err := db.conn.Where(&Order{UserID: userId}).Find(&orders).Error
	for i, o := range orders {
		orders[i] = db.completeOrder(o)
	}
	return orders, err
}

func (db *GormDB) FindAllOrdersWithExpiresAtBetween(start, end time.Time) ([]Order, error) {
	var orders []Order

	err := db.conn.Where("richiesto_il BETWEEN ? AND ?", start, end).Find(&orders).Error
	for i, o := range orders {
		orders[i] = db.completeOrder(o)
	}
	return orders, err
}

func (db *GormDB) CreateOrder(order Order) error {
	return db.conn.Create(&order).Error
}

func (db *GormDB) UpdateOrder(order Order) error {
	columns := []string{"id_prodotto", "id_utente", "quanita", "richiesto_il"}
	return db.conn.Model(&order).Select(columns).Updates(order).Error
}

func (db *GormDB) DeleteOrder(id int) error {
	return db.conn.Delete(&Order{}, id).Error
}
