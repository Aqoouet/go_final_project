// Package main starts the HTTP server, initializes configuration and database, and registers HTTP handlers.
package main

import (
	"fmt"
	"log"
	"net/http"

	"go_final_project/internal/auth"
	"go_final_project/internal/config"
	"go_final_project/internal/database"
	"go_final_project/internal/handlers"
)

func main() {
	cfg := config.LoadConfig()

	if err := database.Init(cfg.DBFile); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Printf("Database initialized: %s", cfg.DBFile)

	// Authentication endpoint (no auth required)
	http.HandleFunc("/api/signin", handlers.SignInHandler)
	
	// Public endpoint (no auth required)
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	
	// Protected endpoints (auth required if TODO_PASSWORD is set)
	http.HandleFunc("/api/task", auth.Middleware(handlers.TaskHandler))
	http.HandleFunc("/api/tasks", auth.Middleware(handlers.TasksHandler))
	http.HandleFunc("/api/task/done", auth.Middleware(handlers.DoneTaskHandler))

	// Static files
	http.Handle("/", http.FileServer(http.Dir(cfg.WebDir)))

	port := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting server at http://localhost%s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
