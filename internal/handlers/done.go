package handlers

import (
	"net/http"

	"go_final_project/internal/database"
	"go_final_project/internal/services"
	"go_final_project/internal/utils"
)

// DoneTaskHandler processes POST requests to mark a task as completed
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// Retrieve task from database
	task, err := database.GetTask(id)
	if err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	// If no repetition rule exists - delete the task
	if task.Repeat == "" {
		if err := database.DeleteTask(id); err != nil {
			utils.WriteError(w, err.Error())
			return
		}
		utils.WriteEmptySuccess(w)
		return
	}

	// For recurring tasks, calculate the next date
	taskDate, err := utils.ParseDate(task.Date)
	if err != nil {
		utils.WriteError(w, "invalid date format")
		return
	}

	// Use task date + 1 day as the base for calculating next date
	// This ensures NextDate computes the next interval after the current task date
	nextDateBase := taskDate.AddDate(0, 0, 1)
	calculator := services.NewNextDateCalculator()
	nextDate, err := calculator.Calculate(nextDateBase, task.Date, task.Repeat)
	if err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	// Update task date
	if err := database.UpdateTaskDate(id, nextDate); err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	// Return empty JSON on success
	utils.WriteEmptySuccess(w)
}
