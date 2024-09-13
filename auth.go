package main

import (
	"database/sql"
	"fmt"
	"net/http"

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

func login(username string, password string, w http.ResponseWriter) error {
	ok, err := verifyPassword(username, password)
	if err != nil {
		return fmt.Errorf("errore verifica password: %v", err)
	}

	if !ok {
		return fmt.Errorf("password errata")
	}

	token, err := generateJWT(username)
	if err != nil {
		return fmt.Errorf("errore generazione JWT: %v", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
	})

	return nil
}

func logout(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})
}
