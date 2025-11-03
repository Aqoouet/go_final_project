package handlers

import (
	"net/http"

	"go_final_project/internal/database"
	"go_final_project/internal/models"
)

// TasksResponse представляет ответ со списком задач
type TasksResponse struct {
	Tasks []*models.Task `json:"tasks"`
}

// TasksHandler обрабатывает GET-запросы для получения списка задач
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр search из URL (если есть)
	search := r.URL.Query().Get("search")

	// Получаем список задач (ограничиваем 50 записями)
	tasks, err := database.Tasks(50, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Возвращаем список задач
	writeJSON(w, TasksResponse{
		Tasks: tasks,
	})
}

