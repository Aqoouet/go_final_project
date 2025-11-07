// Package database provides CRUD operations for scheduler tasks.
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"go_final_project/internal/models"
	"time"
)

const (
	DateFormat = "20060102"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrEmptyID      = errors.New("task ID is required")
)

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

func Tasks(limit int, search string) ([]*models.Task, error) {
	var rows *sql.Rows
	var err error

	if search != "" {
		if date, err := time.Parse("02.01.2006", search); err == nil {
			dateStr := date.Format(DateFormat)
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
			rows, err = db.Query(query, dateStr, limit)
		} else {
			searchPattern := "%" + search + "%"
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
			rows, err = db.Query(query, searchPattern, searchPattern, limit)
		}
	} else {
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		rows, err = db.Query(query, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

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

func UpdateTask(task *models.Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if count == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func DeleteTask(id string) error {
	if id == "" {
		return ErrEmptyID
	}

	query := `DELETE FROM scheduler WHERE id = ?`

	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if count == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func UpdateTaskDate(id string, date string) error {
	if id == "" {
		return ErrEmptyID
	}

	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := db.Exec(query, date, id)
	if err != nil {
		return fmt.Errorf("failed to update task date: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if count == 0 {
		return ErrTaskNotFound
	}

	return nil
}
