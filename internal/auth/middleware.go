// Package auth provides authentication middleware for HTTP handlers.
package auth

import (
	"net/http"
	"os"
)

// Middleware wraps an HTTP handler with authentication check
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if password is set in environment
		password := os.Getenv("TODO_PASSWORD")
		
		// If no password is set, skip authentication
		if len(password) == 0 {
			next(w, r)
			return
		}
		
		// Get JWT token from cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			// No token found
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		
		jwtToken := cookie.Value
		
		// Validate the token
		valid, err := ValidateToken(jwtToken, password)
		if err != nil || !valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		
		// Token is valid, proceed to the next handler
		next(w, r)
	})
}

