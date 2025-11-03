// Package handlers exposes HTTP endpoints for task management and scheduling utilities.
package handlers

import (
	"net/http"

	"go_final_project/internal/database"
	"go_final_project/internal/models"
	"go_final_project/internal/utils"
)

type TasksResponse struct {
	Tasks []*models.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := database.Tasks(50, search)
	if err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	utils.WriteJSON(w, TasksResponse{
		Tasks: tasks,
	})
}
