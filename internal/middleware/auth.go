package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// Initialize cookie store with a secure signing key
var Store = sessions.NewCookieStore([]byte("lordis-super-secret-key-12345"))

// LoginRequired blocks unauthenticated users from seeing staff routes
func LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "lordis-session")

		if _, ok := session.Values["name"]; !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AdminRequired blocks non-admin users from reaching sensitive views
func AdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "lordis-session")

		role, ok := session.Values["role"].(string)
		if !ok || role != "admin" {
			http.Error(w, "Unauthorized: Administration access required.", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
