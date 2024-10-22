package appContext

import (
	"context"
	"gestione-ordini/pkg/database"
	"html/template"
	"net/http"
)

type AppContext struct {
	DB                *database.GormDB
	Templ             *template.Template
	AuthenticatedUser *database.User
	AuthenticationErr error
}

type ctxKey string

const appContextKey ctxKey = "app.context"

func FromRequest(r *http.Request) *AppContext {
	return r.Context().Value(appContextKey).(*AppContext)
}

func WithContext(ctx context.Context, reqCtx *AppContext) context.Context {
	return context.WithValue(ctx, appContextKey, reqCtx)
}
