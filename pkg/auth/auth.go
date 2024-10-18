package auth

import (
	"crypto/rsa"
	"fmt"
	"gestione-ordini/pkg/database"
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

var (
	ErrNoCookie    = fmt.Errorf("cookie not found")
	ErrInvalidJWT  = fmt.Errorf("invalid JWT")
	ErrInvalidRole = fmt.Errorf("invalid role")
	ErrInvalidPerm = fmt.Errorf("invalid permission")
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

type userClaims struct {
	database.User
	jwt.RegisteredClaims
}

func generateJWT(userId int, roleId int) (string, error) {
	claims := userClaims{
		User: database.User{
			ID:     userId,
			RoleID: roleId,
		},
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

func validateJWT(tokenString string) (*userClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*userClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func VerifyPassword(passwordHash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}

const SessionCookieName = "jwt"

func GetAuthenticatedUser(r *http.Request) (*database.User, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return nil, ErrNoCookie
	}

	claims, err := validateJWT(cookie.Value)
	if err != nil {
		return nil, ErrInvalidJWT
	}

	return &claims.User, nil
}

func SetAuthenticatedUser(w http.ResponseWriter, userId int, roleId int) error {
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

func UnsetAuthenticatedUser(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})
}
