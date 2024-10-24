package database

import (
	"time"

	"gorm.io/gorm/clause"
)

type Order struct {
	ID          int       `gorm:"column:id;primaryKey"`
	ProductID   int       `gorm:"column:id_prodotto"`
	Product     Product   `gorm:"foreignKey:ProductID"`
	UserID      int       `gorm:"column:id_utente"`
	User        User      `gorm:"foreignKey:UserID"`
	Amount      int       `gorm:"column:quantita"`
	RequestedAt time.Time `gorm:"column:richiesto_il"`
}

func (Order) TableName() string { return "ordini" }

func (db *GormDB) FindOrder(id int) (Order, error) {
	var order Order

	err := db.conn.Take(&order, id).Error
	return order, err
}

func (db *GormDB) FindAllOrders() ([]Order, error) {
	var orders []Order

	err := db.conn.Preload(clause.Associations).Find(&orders).Error
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
