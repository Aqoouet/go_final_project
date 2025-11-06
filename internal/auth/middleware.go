// Package auth provides authentication middleware for HTTP handlers.
package auth

import (
	"net/http"
	
	"go_final_project/internal/config"
)

func Middleware(cfg *config.Config) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if len(cfg.Password) == 0 {
				next(w, r)
				return
			}
			
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
			
			jwtToken := cookie.Value
			
			valid, err := ValidateToken(jwtToken, cfg.Password, cfg.JWTSecretKey)
			if err != nil || !valid {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
			
			next(w, r)
		})
	}
}

