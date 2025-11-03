// Package handlers exposes HTTP endpoints for task management and scheduling utilities.
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

func handleAddTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	if task.Title == "" {
		utils.WriteError(w, "task title is required")
		return
	}

	if err := validateAndAdjustTaskDate(&task); err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	id, err := database.AddTask(&task)
	if err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	utils.WriteJSON(w, map[string]string{"id": fmt.Sprint(id)})
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	task, err := database.GetTask(id)
	if err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	utils.WriteJSON(w, task)
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	if task.Title == "" {
		utils.WriteError(w, "task title is required")
		return
	}

	if err := validateAndAdjustTaskDate(&task); err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	if err := database.UpdateTask(&task); err != nil {
		utils.WriteError(w, err.Error())
		return
	}

	utils.WriteEmptySuccess(w)
}

func validateAndAdjustTaskDate(task *models.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = utils.FormatDate(now)
		return nil
	}

	taskDate, err := utils.ParseDate(task.Date)
	if err != nil {
		return fmt.Errorf("date is in incorrect format")
	}

	if task.Repeat != "" {
		calculator := services.NewNextDateCalculator()
		_, err = calculator.Calculate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("repeat rule is in incorrect format")
		}
	}

	if utils.IsAfter(now, taskDate) {
		if task.Repeat == "" {
			task.Date = utils.FormatDate(now)
		} else {
			calculator := services.NewNextDateCalculator()
			currentDate := task.Date
			for {
				nextDate, err := calculator.Calculate(now, currentDate, task.Repeat)
				if err != nil {
					return fmt.Errorf("error calculating next date")
				}

				nextTime, err := utils.ParseDate(nextDate)
				if err != nil {
					return fmt.Errorf("error parsing date")
				}

				if !utils.IsAfter(now, nextTime) {
					task.Date = nextDate
					break
				}

				currentDate = nextDate
			}
		}
	}

	return nil
}
