// Package handlers exposes HTTP endpoints for task management and scheduling utilities.
package handlers

import (
	"net/http"

	"go_final_project/internal/database"
	"go_final_project/internal/utils"
)

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

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if err := database.DeleteTask(id); err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	utils.WriteEmptySuccess(w)
}
