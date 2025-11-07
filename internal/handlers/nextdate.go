// Package handlers exposes HTTP endpoints for task management and scheduling utilities.
package handlers

import (
	"net/http"
	"time"

	"go_final_project/internal/services"
	"go_final_project/internal/utils"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(utils.DateFormat, nowStr)
		if err != nil {
			utils.RespondWithError(w, "Invalid format for date \"now\"", http.StatusBadRequest)
			return
		}
	}

	calculator := services.NewNextDateCalculator()
	nextDate, err := calculator.Calculate(now, dateStr, repeat)

	if err != nil {
		utils.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte(nextDate)); err != nil {
		utils.RespondWithError(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
