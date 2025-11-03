package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go_final_project/internal/database"
	"go_final_project/internal/models"
)

// AddTaskHandler обрабатывает POST-запросы для добавления новой задачи
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	// Десериализация JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Проверка обязательного поля title
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	// Проверка и обработка даты
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Добавление задачи в базу данных
	id, err := database.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Возврат идентификатора добавленной задачи
	writeJSON(w, map[string]string{"id": fmt.Sprint(id)})
}

// checkDate проверяет и корректирует дату задачи
func checkDate(task *models.Task) error {
	now := time.Now()

	// Если дата не указана, берём сегодняшнюю
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
		return nil
	}

	// Проверяем корректность формата даты
	taskDate, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("дата представлена в неправильном формате")
	}

	// Если правило повторения указано, проверяем его корректность
	if task.Repeat != "" {
		_, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("правило повторения указано в неправильном формате")
		}
	}

	// Если дата задачи меньше сегодняшней
	if afterNow(now, taskDate) {
		if task.Repeat == "" {
			// Нет правила повторения - берём сегодняшнюю дату
			task.Date = now.Format(DateFormat)
		} else {
			// Есть правило повторения - вычисляем следующую дату после now
			// Вызываем NextDate в цикле, пока не получим дату >= now
			currentDate := task.Date
			for {
				nextDate, err := NextDate(now, currentDate, task.Repeat)
				if err != nil {
					return fmt.Errorf("ошибка вычисления следующей даты")
				}
				
				// Парсим следующую дату
				nextTime, err := time.Parse(DateFormat, nextDate)
				if err != nil {
					return fmt.Errorf("ошибка парсинга даты")
				}
				
				// Если следующая дата >= now, используем её
				if !afterNow(now, nextTime) {
					task.Date = nextDate
					break
				}
				
				// Иначе продолжаем с этой даты
				currentDate = nextDate
			}
		}
	}

	return nil
}

// GetTaskHandler обрабатывает GET-запросы для получения задачи по ID
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из параметров запроса
	id := r.URL.Query().Get("id")
	
	// Получаем задачу из базы данных
	task, err := database.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	
	// Возвращаем задачу
	writeJSON(w, task)
}

// UpdateTaskHandler обрабатывает PUT-запросы для обновления существующей задачи
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	
	// Десериализация JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	
	// Проверка обязательного поля title
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}
	
	// Проверка и обработка даты
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	
	// Обновление задачи в базе данных
	if err := database.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	
	// Возврат пустого JSON-объекта при успехе
	writeJSON(w, map[string]string{})
}

// writeJSON сериализует данные в JSON и отправляет их клиенту
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

