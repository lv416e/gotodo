package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"gotodo/handlers"
	"gotodo/models"
)

func main() {
	store := models.NewTodoStore()
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

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}