package database

import (
	"strconv"
	"time"

	"gorm.io/gorm/clause"
)

type Order struct {
	ID               int       `gorm:"column:id;primaryKey"`
	ProductID        int       `gorm:"column:id_prodotto"`
	Product          Product   `gorm:"foreignKey:ProductID"`
	UserID           int       `gorm:"column:id_utente"`
	User             User      `gorm:"foreignKey:UserID"`
	Amount           int       `gorm:"column:quantita"`
	ExpiresAt        time.Time `gorm:"column:richiesto_il"`
	ExpirationString string    `gorm:"-:all"`
}

func (Order) TableName() string { return "ordini" }

func calculateExpirationString(order Order) Order {
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
	return calculateExpirationString(order), err
}

func (db *GormDB) FindAllOrders() ([]Order, error) {
	var orders []Order

	err := db.conn.Preload(clause.Associations).Find(&orders).Error
	for i, o := range orders {
		orders[i] = calculateExpirationString(o)
	}
	return orders, err
}

func (db *GormDB) FindAllOrdersWithUserID(userId int) ([]Order, error) {
	var orders []Order

	err := db.conn.Preload(clause.Associations).Find(&orders).Where("id_utente = ?", userId).Error
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
