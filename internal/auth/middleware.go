// Package auth provides authentication middleware for HTTP handlers.
package auth

import (
	"net/http"
	"os"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		password := os.Getenv("TODO_PASSWORD")
		
		if len(password) == 0 {
			next(w, r)
			return
		}
		
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		
		jwtToken := cookie.Value
		
		valid, err := ValidateToken(jwtToken, password)
		if err != nil || !valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		
		next(w, r)
	})
}

