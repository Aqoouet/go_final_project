package handlers

import (
	"net/http"
	"time"

	"go_final_project/internal/services"
	"go_final_project/internal/utils"
)

// NextDateHandler handles GET /api/nextdate
// Query parameters:
//   - now: current date reference (optional, defaults to current time)
//   - date: initial date in YYYYMMDD format (required)
//   - repeat: repetition rule (required)
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	// Parse "now" parameter or use current time
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(utils.DateFormat, nowStr)
		if err != nil {
			http.Error(w, "Invalid format for date \"now\"", http.StatusBadRequest)
			return
		}
	}

	// Calculate next date using the service
	calculator := services.NewNextDateCalculator()
	nextDate, err := calculator.Calculate(now, dateStr, repeat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(nextDate))
}
