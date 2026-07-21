package middleware

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte(getSessionSecret()))

func getSessionSecret() string {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "lordis-super-secret-key-change-in-production"
	}
	return secret
}

// LoginRequired ensures that a user has an active session before accessing protected routes
func LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "lordis-session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		email, ok := session.Values["email"].(string)
		if !ok || email == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AdminRequired ensures that the logged-in user possesses admin privileges
func AdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "lordis-session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		role, ok := session.Values["role"].(string)
		if !ok || role != "admin" {
			http.Error(w, "Forbidden: Administrator access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}