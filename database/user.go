package database

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID           int
	RoleID       int
	Username     string
	PasswordHash string
	Name         string
	Surname      string
	CreatedAt    time.Time
}

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
	query := "SELECT nome FROM ruoli WHERE id = ?"
	var name string

	err := db.conn.QueryRow(query, id).Scan(&name)
	return name, err
}

func (db *Database) GetUser(id int) (*User, error) {
	query := "SELECT id, id_ruolo, username, password_hash, nome, cognome, creato_il FROM utenti WHERE id = ?"
	var user User
	var createdAtString string

	err := db.conn.QueryRow(query, id).Scan(
		&user.ID,
		&user.RoleID,
		&user.Username,
		&user.PasswordHash,
		&user.Name,
		&user.Surname,
		&createdAtString,
	)
	if err != nil {
		return nil, err
	}

	user.CreatedAt, err = time.Parse(DateTimeFormat, createdAtString)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) GetUsers(orderBy int) ([]User, error) {
	var orderByString string
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

	query := "SELECT id, id_ruolo, username, password_hash, nome, cognome, creato_il FROM utenti ORDER BY " + orderByString
	rows, err := db.conn.Query(query)
	if err == sql.ErrNoRows {
		return []User{}, nil
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var createdAtString string

		err = rows.Scan(
			&user.ID,
			&user.RoleID,
			&user.Username,
			&user.PasswordHash,
			&user.Name,
			&user.Surname,
			&createdAtString,
		)
		if err != nil {
			return nil, err
		}

		user.CreatedAt, err = time.Parse(DateTimeFormat, createdAtString)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (db *Database) GetUserByUsername(username string) (*User, error) {
	query := "SELECT id FROM utenti WHERE username = ?"
	var id int

	err := db.conn.QueryRow(query, username).Scan(&id)
	if err != nil {
		return nil, err
	}

	return db.GetUser(id)
}

func (db *Database) AddUser(
	roleId int,
	username string,
	passwordHash string,
	name string,
	surname string,
) error {
	createdAt := time.Now().Format(DateTimeFormat)

	query := "INSERT INTO utenti (id_ruolo, username, password_hash, nome, cognome, creato_il) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := db.conn.Exec(query, roleId, username, passwordHash, name, surname, createdAt)

	return err
}

func (db *Database) UserHasPermission(userId int, permissionId int) (bool, error) {
	query := "SELECT u.username FROM utenti u JOIN ruoli r ON u.id_ruolo = r.id JOIN ruolo_permesso rp ON r.id = rp.id_ruolo JOIN permessi p ON rp.id_permesso = p.id WHERE u.id = ? AND p.id = ?"
	row := db.conn.QueryRow(query, userId, permissionId)
	var username string

	err := row.Scan(&username)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
