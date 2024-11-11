package handlers

import (
	"fmt"
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/database"
	"net/http"
)

type firstLoginData struct {
	InitialUsername      string
	InitialPassword      string
	PasswordConfirmError bool
	UsernameError        bool
}

func GetFirstLogin(w http.ResponseWriter, r *http.Request) {
	appContext.ExecuteTemplate(w, r, "firstLogin.html", firstLoginData{})
}

func PostFirstLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("passwordConfirm")

	data := firstLoginData{
		InitialUsername: username,
		InitialPassword: password,
	}

	if password != passwordConfirm {
		data.PasswordConfirmError = true
		appContext.ExecuteTemplate(w, r, "firstLogin.html", data)
		return
	}

	user, err := getEligibleUser(r, username)
	if err != nil {
		data.UsernameError = true
		appContext.ExecuteTemplate(w, r, "firstLogin.html", data)
		return
	}

	user.PasswordHash, err = auth.HashPassword(password)
	if err != nil {
		ShowDatabaseQueryError(w, r, err)
		return
	}

	err = appContext.Database(r).UpdateUser(user)
	if err != nil {
		ShowDatabaseQueryError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getEligibleUser(r *http.Request, username string) (database.User, error) {
	user, err := appContext.Database(r).FindUserWithUsername(username)
	if err != nil {
		return database.User{}, err
	}

	if user.PasswordHash != "" {
		return database.User{}, fmt.Errorf("password already set")
	}

	return user, nil
}
