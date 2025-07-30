package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"gotodo/database"
	"gotodo/handlers"
	"gotodo/models"
)

func main() {
	// Initialize database
	dbPath := getEnv("DB_PATH", "./data/todos.db")
	db, err := database.Initialize(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize store
	store := models.NewTodoStore(db)
	todoHandler := handlers.NewTodoHandler(store)

	// API routes
	http.Handle("/api/todos", todoHandler)
	http.Handle("/api/todos/", todoHandler)

	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			log.Printf("Template error: %v", err)
			return
		}
		
		tmpl.Execute(w, nil)
	})

	port := getEnv("PORT", "8080")
	fmt.Printf("Server starting on http://localhost:%s\n", port)
	fmt.Printf("Database: %s\n", dbPath)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}