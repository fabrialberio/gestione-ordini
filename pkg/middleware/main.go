package middleware

import (
	appContext "gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/handlers"
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

func WithContext(reqCtx *appContext.AppContext, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetAuthenticatedUser(r)
		reqCtx.AuthenticatedUser = user
		reqCtx.AuthenticationErr = err

		ctx := r.Context()
		ctx = appContext.WithContext(ctx, reqCtx)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func WithRole(roleId int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := appContext.FromRequest(r).AuthenticatedUser
		if user == nil {
			handlers.HandleError(w, r, auth.ErrNoCookie)
			return
		}

		if user.RoleID != roleId {
			handlers.HandleError(w, r, auth.ErrInvalidRole)
			return
		}

		next.ServeHTTP(w, r)
	})
}
