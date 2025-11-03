package main

import (
	"fmt"
	"log"
	"net/http"

	"go_final_project/internal/config"
	"go_final_project/internal/database"
	"go_final_project/internal/handlers"
)

func main() {
	// Load configuration from environment variables
	cfg := config.LoadConfig()

	// Initialize database connection
	if err := database.Init(cfg.DBFile); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Printf("Database initialized: %s", cfg.DBFile)

	// Register API handlers
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/task", handlers.TaskHandler)
	http.HandleFunc("/api/tasks", handlers.TasksHandler)
	http.HandleFunc("/api/task/done", handlers.DoneTaskHandler)

	// Serve static files from web directory
	http.Handle("/", http.FileServer(http.Dir(cfg.WebDir)))

	// Start HTTP server
	port := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting server at http://localhost%s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
