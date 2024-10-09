package main

import (
	"log"
	"net/http"
)

func logRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		h(w, r)
	}
}

func htmlError(w http.ResponseWriter, status int) {
	message := ""

	switch status {
	case http.StatusForbidden:
		message = "Non autorizzato"
	case http.StatusInternalServerError:
		message = "Errore interno"
	case http.StatusMethodNotAllowed:
		message = "Metodo non consentito"
	}

	http.Error(w, message, status)
}

func checkPerm(w http.ResponseWriter, r *http.Request, permId int) bool {
	claims, err := getSessionCookie(r)
	if err != nil {
		htmlError(w, http.StatusUnauthorized)
		return false
	}

	if ok, err := db.UserHasPerm(claims.UserID, permId); err != nil || !ok {
		htmlError(w, http.StatusUnauthorized)
		return false
	}

	return true
}

func checkRole(w http.ResponseWriter, r *http.Request, roleId int) bool {
	claims, err := getSessionCookie(r)
	if err != nil {
		htmlError(w, http.StatusUnauthorized)
		return false
	}

	if claims.RoleID != roleId {
		htmlError(w, http.StatusUnauthorized)
		return false
	}

	return true
}
