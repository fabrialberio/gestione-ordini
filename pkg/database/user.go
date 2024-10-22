package database

import (
	"fmt"
	"time"

	"gorm.io/gorm/clause"
)

type User struct {
	ID           int       `gorm:"column:id;primaryKey"`
	RoleID       int       `gorm:"column:id_ruolo"`
	Role         Role      `gorm:"foreignKey:RoleID"`
	Username     string    `gorm:"column:username"`
	PasswordHash string    `gorm:"column:password_hash"`
	Name         string    `gorm:"column:nome"`
	Surname      string    `gorm:"column:cognome"`
	CreatedAt    time.Time `gorm:"column:creato_il"`
}

func (User) TableName() string { return "utenti" }

type Role struct {
	ID   int64  `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:nome;size:255"`
}

func (Role) TableName() string { return "ruoli" }

const (
	RoleIDCook int = iota + 1
	RoleIDManager
	RoleIDAdministrator
)

const (
	PermIDOrders int = iota + 1
	PermIDOwnOrder
	PermIDProductTypes
	PermIDSuppliers
	PermIDUnitsOfMeasure
	PermIDProducts
	PermIDUsers
	PermIDOwnUser
)

const (
	OrderUserByID int = iota + 1
	OrderUserByRole
	OrderUserByUsername
	OrderUserByName
	OrderUserBySurname
	OrderUserByCreatedAt
)

func (db *GormDB) FindAllRoles() ([]Role, error) {
	var roles []Role

	err := db.conn.Find(&roles).Error
	return roles, err
}

func (db *GormDB) FindUser(id int) (User, error) {
	var user User

	err := db.conn.Take(&user, id).Error
	return user, err
}

func (db *GormDB) FindAllUsers(orderBy int, orderDesc bool) ([]User, error) {
	var orderByString string
	var users []User

	switch orderBy {
	case OrderUserByID:
		orderByString = "id"
	case OrderUserByRole:
		orderByString = "id_ruolo"
	case OrderUserByUsername:
		orderByString = "username"
	case OrderUserByName:
		orderByString = "nome"
	case OrderUserBySurname:
		orderByString = "cognome"
	case OrderUserByCreatedAt:
		orderByString = "creato_il"
	default:
		return nil, fmt.Errorf("invalid orderBy value: %d", orderBy)
	}

	err := db.conn.Preload(clause.Associations).Order(clause.OrderByColumn{
		Column: clause.Column{Name: orderByString},
		Desc:   orderDesc,
	}).Find(&users).Error
	return users, err
}

func (db *GormDB) FindUserWithUsername(username string) (*User, error) {
	var user *User

	err := db.conn.Model(&User{}).Take(&user, "username = ?", username).Error
	return user, err
}

func (db *GormDB) CreateUser(user User) error {
	return db.conn.Create(&user).Error
}

func (db *GormDB) UpdateUser(user User) error {
	columns := []string{"id_ruolo", "username", "nome", "cognome"}
	if user.PasswordHash != "" {
		columns = append(columns, "password_hash")
	}

	return db.conn.Model(&user).Select(columns).Updates(user).Error
}

func (db *GormDB) DeleteUser(id int) error {
	return db.conn.Delete(&User{}, id).Error
}

func (db *GormDB) UserHasPerm(userId int, permId int) (bool, error) {
	rows, err := db.conn.Model(&User{}).
		Select("utenti.id").
		Joins("JOIN ruoli ON utenti.id_ruolo = ruoli.id").
		Joins("JOIN ruolo_permesso ON ruoli.id = ruolo_permesso.id_ruolo").
		Joins("JOIN permessi ON ruolo_permesso.id_permesso = permessi.id").
		Where("utenti.id = ? AND permessi.id = ?", userId, permId).
		Rows()

	if err != nil {
		return false, err
	} else if rows.Next() {
		return true, nil
	}

	return false, nil
}
