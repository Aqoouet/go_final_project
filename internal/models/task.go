package models

// Task represents a scheduler task in the system
type Task struct {
	ID      string `json:"id"`      // Task unique identifier
	Date    string `json:"date"`    // Task date in YYYYMMDD format
	Title   string `json:"title"`   // Task title (required)
	Comment string `json:"comment"` // Optional task description
	Repeat  string `json:"repeat"`  // Repetition rule (e.g., "y", "d 7", "w 1,3,5")
}
