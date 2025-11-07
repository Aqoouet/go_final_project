// Package handlers exposes HTTP endpoints for task management and scheduling utilities.
package handlers

import (
	"encoding/json"
	"errors"
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
		utils.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		utils.RespondWithError(w, "task title is required", http.StatusBadRequest)
		return
	}

	if err := validateAndAdjustTaskDate(&task); err != nil {
		utils.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := database.AddTask(&task)
	if err != nil {
		utils.RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"id": fmt.Sprint(id)})
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	task, err := database.GetTask(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, database.ErrTaskNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, database.ErrEmptyID) {
			statusCode = http.StatusBadRequest
		}
		utils.RespondWithError(w, err.Error(), statusCode)
		return
	}

	utils.WriteJSON(w, task)
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		utils.RespondWithError(w, "task id is required", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		utils.RespondWithError(w, "task title is required", http.StatusBadRequest)
		return
	}

	if err := validateAndAdjustTaskDate(&task); err != nil {
		utils.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.UpdateTask(&task); err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, database.ErrTaskNotFound) {
			statusCode = http.StatusNotFound
		}
		utils.RespondWithError(w, err.Error(), statusCode)
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
			const maxIterations = 1000
			
			for i := 0; i < maxIterations; i++ {
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
					return nil
				}

				currentDate = nextDate
			}
			
			return fmt.Errorf("unable to find future date within reasonable iterations")
		}
	}

	return nil
}
