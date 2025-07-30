package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gotodo/models"
)

type TodoHandler struct {
	store *models.TodoStore
}

func NewTodoHandler(store *models.TodoStore) *TodoHandler {
	return &TodoHandler{store: store}
}

func (h *TodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.Method {
	case http.MethodGet:
		h.getTodos(w, r)
	case http.MethodPost:
		h.createTodo(w, r)
	case http.MethodPut:
		h.updateTodo(w, r)
	case http.MethodDelete:
		h.deleteTodo(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TodoHandler) getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.store.GetAll()
	if err != nil {
		log.Printf("Error getting todos: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	
	todo, err := h.store.Create(req.Title)
	if err != nil {
		log.Printf("Error creating todo: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) updateTodo(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	parts := strings.Split(path, "/")
	
	if len(parts) != 2 || parts[1] != "toggle" {
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}
	
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	todo, err := h.store.Toggle(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		log.Printf("Error toggling todo: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	err = h.store.Delete(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		log.Printf("Error deleting todo: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}