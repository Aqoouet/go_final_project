// Package handlers provides HTTP handlers for authentication endpoints.
package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"go_final_project/internal/auth"
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
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}
	
	// Get password from environment
	envPassword := os.Getenv("TODO_PASSWORD")
	if envPassword == "" {
		utils.RespondWithError(w, "Аутентификация не настроена", http.StatusInternalServerError)
		return
	}
	
	// Parse request body
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	
	// Validate password
	if req.Password != envPassword {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(SignInResponse{
			Error: "Неверный пароль",
		})
		return
	}
	
	// Generate JWT token
	token, err := auth.GenerateToken(req.Password)
	if err != nil {
		utils.RespondWithError(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}
	
	// Return token
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SignInResponse{
		Token: token,
	})
}

