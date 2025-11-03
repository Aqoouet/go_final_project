package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go_final_project/internal/database"
	"go_final_project/internal/models"
	"go_final_project/internal/services"
	"go_final_project/internal/utils"
)

// handleAddTask processes POST requests to add a new task
func handleAddTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	// Validate required field
	if task.Title == "" {
		utils.WriteError(w, "task title is required")
		return
	}

	// Validate and adjust task date
	if err := validateAndAdjustTaskDate(&task); err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	// Add task to database
	id, err := database.AddTask(&task)
	if err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	// Return the ID of the created task
	utils.WriteJSON(w, map[string]string{"id": fmt.Sprint(id)})
}

// handleGetTask processes GET requests to retrieve a task by ID
func handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// Retrieve task from database
	task, err := database.GetTask(id)
	if err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	// Return the task
	utils.WriteJSON(w, task)
}

// handleUpdateTask processes PUT requests to update an existing task
func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	// Validate required field
	if task.Title == "" {
		utils.WriteError(w, "task title is required")
		return
	}

	// Validate and adjust task date
	if err := validateAndAdjustTaskDate(&task); err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	// Update task in database
	if err := database.UpdateTask(&task); err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	// Return empty JSON on success
	utils.WriteEmptySuccess(w)
}

// validateAndAdjustTaskDate validates and adjusts the task date
func validateAndAdjustTaskDate(task *models.Task) error {
	now := time.Now()

	// If date is not specified, use today
	if task.Date == "" {
		task.Date = utils.FormatDate(now)
		return nil
	}

	// Validate date format
	taskDate, err := utils.ParseDate(task.Date)
	if err != nil {
		return fmt.Errorf("date is in incorrect format")
	}

	// If repetition rule is specified, validate it
	if task.Repeat != "" {
		calculator := services.NewNextDateCalculator()
		_, err = calculator.Calculate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("repeat rule is in incorrect format")
		}
	}

	// If task date is in the past
	if utils.IsAfter(now, taskDate) {
		if task.Repeat == "" {
			// No repetition rule - use today's date
			task.Date = utils.FormatDate(now)
		} else {
			// Has repetition rule - calculate next valid date
			calculator := services.NewNextDateCalculator()
			currentDate := task.Date
			for {
				nextDate, err := calculator.Calculate(now, currentDate, task.Repeat)
				if err != nil {
					return fmt.Errorf("error calculating next date")
				}

				// Parse the next date
				nextTime, err := utils.ParseDate(nextDate)
				if err != nil {
					return fmt.Errorf("error parsing date")
				}

				// If next date is >= now, use it
				if !utils.IsAfter(now, nextTime) {
					task.Date = nextDate
					break
				}

				// Otherwise continue with this date
				currentDate = nextDate
			}
		}
	}

	return nil
}
