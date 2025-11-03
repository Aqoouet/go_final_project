package handlers

import (
	"net/http"

	"go_final_project/internal/database"
	"go_final_project/internal/models"
	"go_final_project/internal/utils"
)

// TasksResponse represents the response structure for the tasks list
type TasksResponse struct {
	Tasks []*models.Task `json:"tasks"`
}

// TasksHandler processes GET requests to retrieve a list of tasks
// Query parameters:
//   - search: optional search string (searches by title, comment, or date)
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	// Retrieve tasks list (limited to 50 records)
	tasks, err := database.Tasks(50, search)
	if err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	// Return tasks list
	utils.WriteJSON(w, TasksResponse{
		Tasks: tasks,
	})
}
