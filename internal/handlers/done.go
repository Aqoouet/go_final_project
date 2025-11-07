// Package handlers exposes HTTP endpoints for task management and scheduling utilities.
package handlers

import (
	"net/http"

	"go_final_project/internal/database"
	"go_final_project/internal/services"
	"go_final_project/internal/utils"
)

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		utils.WriteError(w, "task id is required")
		return
	}

	task, err := database.GetTask(id)
	if err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	if task.Repeat == "" {
		if err := database.DeleteTask(id); err != nil {
			utils.WriteError(w, err.Error())
			return
		}
		utils.WriteEmptySuccess(w)
		return
	}

	taskDate, err := utils.ParseDate(task.Date)
	if err != nil {
		utils.WriteError(w, "invalid date format")
		return
	}

	nextDateBase := taskDate.AddDate(0, 0, 1)
	calculator := services.NewNextDateCalculator()
	nextDate, err := calculator.Calculate(nextDateBase, task.Date, task.Repeat)
	if err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	if err := database.UpdateTaskDate(id, nextDate); err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	utils.WriteEmptySuccess(w)
}
