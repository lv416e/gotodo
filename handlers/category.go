package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gotodo/models"
)

type CategoryHandler struct {
	store *models.CategoryStore
}

func NewCategoryHandler(store *models.CategoryStore) *CategoryHandler {
	return &CategoryHandler{store: store}
}

func (h *CategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.Method {
	case http.MethodGet:
		h.getCategories(w, r)
	case http.MethodPost:
		h.createCategory(w, r)
	case http.MethodPut:
		h.updateCategory(w, r)
	case http.MethodDelete:
		h.deleteCategory(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) getCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.store.GetAll()
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) createCategory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	if req.Color == "" {
		req.Color = "#007bff" // Default color
	}
	
	category, err := h.store.Create(req.Name, req.Color)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Error(w, "Category already exists", http.StatusConflict)
			return
		}
		log.Printf("Error creating category: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) updateCategory(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	if req.Color == "" {
		req.Color = "#007bff"
	}
	
	category, err := h.store.Update(id, req.Name, req.Color)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Error(w, "Category name already exists", http.StatusConflict)
			return
		}
		log.Printf("Error updating category: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) deleteCategory(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	err = h.store.Delete(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "category is in use") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		log.Printf("Error deleting category: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}