package utils

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a standard error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteJSON serializes data to JSON and writes it to the response
// Sets appropriate Content-Type header
func WriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// WriteError writes an error response with the given error message
func WriteError(w http.ResponseWriter, errorMsg string) {
	WriteJSON(w, ErrorResponse{Error: errorMsg})
}

// WriteEmptySuccess writes an empty JSON object to indicate success
func WriteEmptySuccess(w http.ResponseWriter) {
	WriteJSON(w, map[string]string{})
}

