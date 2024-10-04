package database

import (
	"fmt"
	"time"

	"gorm.io/gorm/clause"
)

type User struct {
	ID           int       `gorm:"column:id;primaryKey"`
	RoleID       int       `gorm:"column:id_ruolo"`
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
	RoleIDCook int = iota
	RoleIDManager
	RoleIDAdministrator
)

const (
	PermIDViewProducts int = iota
	PermIDViewAllOrders
	PermIDEditProducts
	PermIDEditOwnOrder
	PermIDEditAllOrders
	PermIDEditUsers
)

const (
	UserOrderByID int = iota
	UserOrderByRole
	UserOrderByUsername
	UserOrderByName
	UserOrderBySurname
	UserOrderByCreatedAt
)

func (db *Database) GetRoleName(id int) (string, error) {
	var role Role

	if db.conn.Take(&role, id).Error != nil {
		return "", fmt.Errorf("role with ID %d not found", id)
	}

	return role.Name, nil
}

func (db *Database) GetUser(id int) (*User, error) {
	var user User

	if db.conn.Take(&user, id).Error != nil {
		return nil, fmt.Errorf("user with ID %d not found", id)
	}

	return &user, nil
}

func (db *Database) GetUsers(orderBy int) ([]User, error) {
	var orderByString string
	var users []User

	switch orderBy {
	case UserOrderByID:
		orderByString = "id"
	case UserOrderByRole:
		orderByString = "id_ruolo"
	case UserOrderByUsername:
		orderByString = "username"
	case UserOrderByName:
		orderByString = "nome"
	case UserOrderBySurname:
		orderByString = "cognome"
	case UserOrderByCreatedAt:
		orderByString = "creato_il"
	default:
		return nil, fmt.Errorf("invalid orderBy value: %d", orderBy)
	}

	err := db.conn.Find(&users).Order(clause.OrderByColumn{
		Column: clause.Column{Name: orderByString},
		Desc:   true,
	}).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (db *Database) GetUserByUsername(username string) (*User, error) {
	var user *User

	err := db.conn.Model(&User{}).Take(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *Database) AddUser(
	roleId int,
	username string,
	passwordHash string,
	name string,
	surname string,
) error {
	err := db.conn.Model(&User{}).Create(&User{
		RoleID:       roleId,
		Username:     username,
		PasswordHash: passwordHash,
		Name:         name,
		Surname:      surname,
		CreatedAt:    time.Now(),
	}).Error
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UserHasPerm(userId int, permissionId int) (bool, error) {
	rows, err := db.conn.Model(&User{}).
		Select("utenti.id").
		Joins("JOIN ruoli ON utenti.id_ruolo = ruoli.id").
		Joins("JOIN ruolo_permesso ON ruoli.id = ruolo_permesso.id_ruolo").
		Joins("JOIN permessi ON ruolo_permesso.id_permesso = permessi.id").
		Where("utenti.id = ? AND permessi.id = ?", userId, permissionId).
		Rows()

	if err != nil {
		return false, err
	} else if rows.Next() {
		return true, nil
	}

	return false, nil
}
