package handlers

import (
	"net/http"

	"go_final_project/internal/database"
)

// TaskHandler главный обработчик для /api/task, который маршрутизирует запросы
// в зависимости от HTTP-метода
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetTaskHandler(w, r)
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// DeleteTaskHandler обрабатывает DELETE-запросы для удаления задачи
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из параметров запроса
	id := r.URL.Query().Get("id")
	
	// Удаляем задачу из базы данных
	if err := database.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	
	// Возвращаем пустой JSON при успехе
	writeJSON(w, map[string]string{})
}

