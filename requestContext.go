package main

import (
	"context"
	"gestione-ordini/database"
	"html/template"
	"net/http"
)

type RequestContext struct {
	DB                *database.GormDB
	Templ             *template.Template
	AuthenticatedUser *database.User
	AuthenticationErr error
}

type ctxKey string

const requestContextKey ctxKey = "requestContext"

func GetRequestContext(r *http.Request) RequestContext {
	return r.Context().Value(requestContextKey).(RequestContext)
}

func StoreRequestContext(ctx context.Context, reqCtx RequestContext) context.Context {
	return context.WithValue(ctx, requestContextKey, reqCtx)
}
