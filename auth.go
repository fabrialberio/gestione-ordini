package main

import (
	"crypto/rsa"
	"fmt"
	"gestione-ordini/database"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	privateKeyBytes, err := os.ReadFile("private.key")
	if err != nil {
		log.Fatalf("Error reading private key: %v", err)
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		log.Fatalf("Error parsing private key: %v", err)
	}

	publicKeyBytes, err := os.ReadFile("public.key")
	if err != nil {
		log.Fatalf("Error reading public key: %v", err)
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		log.Fatalf("Error parsing public key: %v", err)
	}
}

type UserClaims struct {
	UserID int `json:"userId"`
	RoleID int `json:"roleId"`
	jwt.RegisteredClaims
}

func generateJWT(userId int, roleId int) (string, error) {
	claims := UserClaims{
		UserID: userId,
		RoleID: roleId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
			Issuer:    "gestione-ordini",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateJWT(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func verifyPassword(user *database.User, password string) bool {
	if user == nil {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

const SessionCookieName = "jwt"

func getSessionCookie(r *http.Request) (*UserClaims, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return nil, err
	}

	claims, err := validateJWT(cookie.Value)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func setSessionCookie(w http.ResponseWriter, userId int, roleId int) error {
	token, err := generateJWT(userId, roleId)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		HttpOnly: true,
	})

	return nil
}

func unsetSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})
}
