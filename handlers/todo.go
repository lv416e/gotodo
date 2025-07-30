package handlers

import (
	"encoding/json"
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
	todos := h.store.GetAll()
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
	
	todo := h.store.Create(req.Title)
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
	
	todo, found := h.store.Toggle(id)
	if !found {
		http.Error(w, "Todo not found", http.StatusNotFound)
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
	
	if !h.store.Delete(id) {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}