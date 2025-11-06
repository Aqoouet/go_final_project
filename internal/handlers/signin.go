// Package handlers provides HTTP handlers for authentication endpoints.
package handlers

import (
	"encoding/json"
	"net/http"

	"go_final_project/internal/auth"
	"go_final_project/internal/config"
	"go_final_project/internal/utils"
)

// SignInRequest represents the signin request payload
type SignInRequest struct {
	Password string `json:"password"`
}

// SignInResponse represents the successful signin response
type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// SignInHandler handles POST /api/signin - authenticates user with password
func SignInHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.RespondWithError(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}
		
		if cfg.Password == "" {
			utils.RespondWithError(w, "Аутентификация не настроена", http.StatusInternalServerError)
			return
		}
		
		var req SignInRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}
		
		if req.Password != cfg.Password {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(SignInResponse{
				Error: "Неверный пароль",
			})
			return
		}
		
		token, err := auth.GenerateToken(req.Password, cfg.JWTSecretKey)
		if err != nil {
			utils.RespondWithError(w, "Ошибка генерации токена", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SignInResponse{
			Token: token,
		})
	}
}

