package handlers

import (
	"net/http"

	"go_final_project/internal/database"
	"go_final_project/internal/utils"
)

// TaskHandler routes /api/task requests based on HTTP method
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetTask(w, r)
	case http.MethodPost:
		handleAddTask(w, r)
	case http.MethodPut:
		handleUpdateTask(w, r)
	case http.MethodDelete:
		handleDeleteTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleDeleteTask processes DELETE requests to remove a task
func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// Delete task from database
	if err := database.DeleteTask(id); err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	// Return empty JSON on success
	utils.WriteEmptySuccess(w)
}
