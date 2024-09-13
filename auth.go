package main

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func verifyPassword(username string, password string) (bool, error) {
	passwordHash, err := db.GetPasswordHashByUsername(username)
	if err == sql.ErrNoRows { // User not found
		return false, nil
	} else if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil { // Passwords don't match
		return false, nil
	}

	return true, nil
}
