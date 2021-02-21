package csrf

import (
	"errors"
	"net/http"
	"time"

	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/google/uuid"
)

type Module struct {
	lib.BlankModule
}

func (m *Module) ID() db.ID {
	return db.MakeID("738f74a2-6d1f-4633-9f60-b3f47914cd2c")
}

func (m *Module) Name() string {
	return "csrf"
}

func (m *Module) Package() string {
	return "awans.org/aft/internal/csrf"
}

func (m *Module) ProvideMiddleware() []lib.Middleware {
	return []lib.Middleware{CSRF}
}

var ttl = 30 * time.Minute

func getCSRFCookie(r *http.Request) (string, bool) {
	csrfCookie, err := r.Cookie("csrf")
	if err != nil || csrfCookie.Value == "" {
		return "", false
	}
	return csrfCookie.Value, true
}

func getCSRFHeader(r *http.Request) (string, bool) {
	val := r.Header.Get("X-CSRF")
	return val, val != ""
}

func setCSRFCookie(w http.ResponseWriter) {
	expires := time.Now().Add(ttl)
	val := uuid.New()

	c := &http.Cookie{
		Name:     "csrf",
		Value:    val.String(),
		Expires:  expires,
		Domain:   "",
		Path:     "/",
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, c)
}

func CSRF(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieVal, set := getCSRFCookie(r)
		if !set {
			setCSRFCookie(w)
		}
		if r.Method != "GET" && !set {
			lib.WriteError(w, errors.New("Invalid request; must include CSRF cookie"))
			return
		} else if r.Method != "GET" {
			headerVal, set := getCSRFHeader(r)
			if !set || (headerVal != cookieVal) {
				lib.WriteError(w, errors.New("Invalid request; CSRF cookie didn't match X-CSRF header"))
				return
			}
		}
		inner.ServeHTTP(w, r)
	})
}

func GetModule() lib.Module {
	return &Module{}
}
