package auth

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"fmt"
	"net/http"
)

func makeAuthMiddleware(appDB db.DB) lib.Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := w.Header().Get("Authorization")
			if token != "" {
				user, err := UserForToken(appDB, token)
				fmt.Printf("AUTH: user: %v err: %v\n", user, err)
			}

			inner.ServeHTTP(w, r)
		})
	}
}
