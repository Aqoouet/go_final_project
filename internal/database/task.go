package database

import (
	"database/sql"
	"fmt"
	"go_final_project/internal/models"
	"time"
)

// AddTask добавляет задачу в базу данных и возвращает идентификатор добавленной записи
func AddTask(task *models.Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	
	return id, nil
}

// Tasks возвращает список задач из базы данных с опциональным поиском
// limit - максимальное количество записей
// search - строка поиска (поиск по заголовку, комментарию или дате в формате DD.MM.YYYY)
func Tasks(limit int, search string) ([]*models.Task, error) {
	var rows *sql.Rows
	var err error

	// Если указан поиск
	if search != "" {
		// Проверяем, является ли search датой в формате DD.MM.YYYY
		if date, err := time.Parse("02.01.2006", search); err == nil {
			// Преобразуем в формат 20060102
			dateStr := date.Format("20060102")
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
			rows, err = db.Query(query, dateStr, limit)
		} else {
			// Поиск по заголовку и комментарию
			searchPattern := "%" + search + "%"
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
			rows, err = db.Query(query, searchPattern, searchPattern, limit)
		}
	} else {
		// Без поиска - просто выбираем все задачи с лимитом
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		rows, err = db.Query(query, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	// Инициализируем пустой слайс (не nil, чтобы в JSON получился [] вместо null)
	tasks := make([]*models.Task, 0)

	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// GetTask возвращает задачу по её идентификатору
func GetTask(id string) (*models.Task, error) {
	if id == "" {
		return nil, fmt.Errorf("не указан идентификатор")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	
	task := &models.Task{}
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("задача не найдена")
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	
	return task, nil
}

// UpdateTask обновляет существующую задачу в базе данных
func UpdateTask(task *models.Task) error {
	if task.ID == "" {
		return fmt.Errorf("не указан идентификатор")
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	
	// Проверяем, была ли обновлена хотя бы одна запись
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	
	return nil
}

// DeleteTask удаляет задачу из базы данных по её идентификатору
func DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("не указан идентификатор")
	}

	query := `DELETE FROM scheduler WHERE id = ?`
	
	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	
	// Проверяем, была ли удалена хотя бы одна запись
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	
	return nil
}

// UpdateTaskDate обновляет только дату задачи
func UpdateTaskDate(id string, date string) error {
	if id == "" {
		return fmt.Errorf("не указан идентификатор")
	}

	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	
	res, err := db.Exec(query, date, id)
	if err != nil {
		return fmt.Errorf("failed to update task date: %w", err)
	}
	
	// Проверяем, была ли обновлена хотя бы одна запись
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	
	return nil
}

