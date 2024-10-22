package appContext

import (
	"context"
	"gestione-ordini/pkg/database"
	"html/template"
	"log"
	"net/http"
)

type appContext struct {
	db                *database.GormDB
	tmpl              *template.Template
	authenticatedUser *database.User
	authenticationErr error
}

type ctxKey string

const appContextKey ctxKey = "app.context"

func fromRequest(r *http.Request) *appContext {
	return r.Context().Value(appContextKey).(*appContext)
}

func New(db *database.GormDB, tmpl *template.Template, authenticatedUser *database.User, authenticationErr error) appContext {
	return appContext{db: db, tmpl: tmpl, authenticatedUser: authenticatedUser, authenticationErr: authenticationErr}
}

func NewContext(ctx context.Context, appContext appContext) context.Context {
	return context.WithValue(ctx, appContextKey, appContext)
}

func ExecuteTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	appCtx := fromRequest(r)
	err := appCtx.tmpl.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		log.Println(err)
	}
}

func Database(r *http.Request) *database.GormDB {
	return fromRequest(r).db
}

func AuthenticatedUser(r *http.Request) (*database.User, error) {
	return fromRequest(r).authenticatedUser, fromRequest(r).authenticationErr
}
