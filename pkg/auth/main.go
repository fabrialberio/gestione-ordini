package auth

import (
	"errors"
	"gestione-ordini/pkg/database"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNoCookie        = errors.New("cookie not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidJWT      = errors.New("invalid JWT")
	ErrInvalidRole     = errors.New("invalid role")
	ErrInvalidPerm     = errors.New("invalid permission")
)

const sessionCookieName = "jwt"

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func verifyPassword(passwordHash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}

func GetAuthenticatedUser(r *http.Request) (*database.User, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, ErrNoCookie
	}

	claims, err := validateJWT(cookie.Value)
	if err != nil {
		return nil, ErrInvalidJWT
	}

	return &claims.User, nil
}

func SetAuthenticatedUser(w http.ResponseWriter, user *database.User, password string) error {
	ok := verifyPassword(user.PasswordHash, password)
	if !ok {
		return ErrInvalidPassword
	}

	token, err := generateJWT(user.ID, user.RoleID)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		HttpOnly: true,
	})

	return nil
}

func UnsetAuthenticatedUser(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})
}
