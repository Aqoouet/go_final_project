package handlers

import (
	"net/http"
	"time"

	"go_final_project/internal/database"
)

// DoneTaskHandler обрабатывает POST-запросы для отметки задачи как выполненной
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из параметров запроса
	id := r.URL.Query().Get("id")
	
	// Получаем задачу из базы данных
	task, err := database.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	
	// Если правило повторения отсутствует - удаляем задачу
	if task.Repeat == "" {
		if err := database.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		// Возвращаем пустой JSON при успехе
		writeJSON(w, map[string]string{})
		return
	}
	
	// Для периодической задачи вычисляем следующую дату
	// Парсим текущую дату задачи
	taskDate, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		writeJSON(w, map[string]string{"error": "неверный формат даты"})
		return
	}
	
	// Используем дату задачи + 1 день как базу для расчета следующей даты
	// Это гарантирует, что NextDate вычислит следующий интервал после текущей даты задачи
	nextDateBase := taskDate.AddDate(0, 0, 1)
	nextDate, err := NextDate(nextDateBase, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	
	// Обновляем дату задачи
	if err := database.UpdateTaskDate(id, nextDate); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	
	// Возвращаем пустой JSON при успехе
	writeJSON(w, map[string]string{})
}

