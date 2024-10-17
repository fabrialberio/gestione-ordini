package main

import (
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}

func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	})
}

func WithContext(reqCtx RequestContext, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := GetAuthenticatedUser(r)
		reqCtx.AuthenticatedUser = user
		reqCtx.AuthenticationErr = err

		ctx := r.Context()
		ctx = StoreRequestContext(ctx, reqCtx)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func WithRole(roleId int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetRequestContext(r).AuthenticatedUser
		if user == nil {
			HandleError(w, r, ErrNoCookie)
			return
		}

		if user.RoleID != roleId {
			HandleError(w, r, ErrInvalidRole)
			return
		}

		next.ServeHTTP(w, r)
	})
}
