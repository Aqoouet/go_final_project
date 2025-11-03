package database

import (
	"database/sql"
	"errors"
	"fmt"
	"go_final_project/internal/models"
	"time"
)

// Common errors
var (
	ErrTaskNotFound = errors.New("task not found")
	ErrEmptyID      = errors.New("task ID is required")
)

// AddTask adds a task to the database and returns the ID of the inserted record
func AddTask(task *models.Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("failed to insert task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return id, nil
}

// Tasks retrieves a list of tasks from the database with optional search
// Parameters:
//   - limit: maximum number of records to return
//   - search: search string (searches by title, comment, or date in DD.MM.YYYY format)
func Tasks(limit int, search string) ([]*models.Task, error) {
	var rows *sql.Rows
	var err error

	// If search is specified
	if search != "" {
		// Check if search is a date in DD.MM.YYYY format
		if date, err := time.Parse("02.01.2006", search); err == nil {
			// Convert to YYYYMMDD format
			dateStr := date.Format("20060102")
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
			rows, err = db.Query(query, dateStr, limit)
		} else {
			// Search by title and comment
			searchPattern := "%" + search + "%"
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
			rows, err = db.Query(query, searchPattern, searchPattern, limit)
		}
	} else {
		// No search - return all tasks with limit
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		rows, err = db.Query(query, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	// Initialize empty slice (not nil, so JSON will be [] instead of null)
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

// GetTask retrieves a task by its ID
func GetTask(id string) (*models.Task, error) {
	if id == "" {
		return nil, ErrEmptyID
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	task := &models.Task{}
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err == sql.ErrNoRows {
		return nil, ErrTaskNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

// UpdateTask updates an existing task in the database
func UpdateTask(task *models.Task) error {
	if task.ID == "" {
		return ErrEmptyID
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	// Check if at least one row was updated
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if count == 0 {
		return ErrTaskNotFound
	}

	return nil
}

// DeleteTask deletes a task from the database by its ID
func DeleteTask(id string) error {
	if id == "" {
		return ErrEmptyID
	}

	query := `DELETE FROM scheduler WHERE id = ?`

	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	// Check if at least one row was deleted
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if count == 0 {
		return ErrTaskNotFound
	}

	return nil
}

// UpdateTaskDate updates only the date field of a task
func UpdateTaskDate(id string, date string) error {
	if id == "" {
		return ErrEmptyID
	}

	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := db.Exec(query, date, id)
	if err != nil {
		return fmt.Errorf("failed to update task date: %w", err)
	}

	// Check if at least one row was updated
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if count == 0 {
		return ErrTaskNotFound
	}

	return nil
}
