package models

import (
	"database/sql"
	"fmt"
	"time"

	"gotodo/database"
)

// Category represents a TODO category
type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CategoryStore manages category items using SQLite database
type CategoryStore struct {
	db *database.DB
}

// NewCategoryStore creates a new CategoryStore with SQLite backend
func NewCategoryStore(db *database.DB) *CategoryStore {
	return &CategoryStore{
		db: db,
	}
}

// GetAll retrieves all categories from the database
func (cs *CategoryStore) GetAll() ([]Category, error) {
	query := `
		SELECT id, name, color, created_at, updated_at 
		FROM categories 
		ORDER BY name ASC
	`
	
	rows, err := cs.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Color,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return categories, nil
}

// GetByID retrieves a specific category by ID
func (cs *CategoryStore) GetByID(id int) (*Category, error) {
	query := `
		SELECT id, name, color, created_at, updated_at 
		FROM categories 
		WHERE id = ?
	`
	
	var category Category
	err := cs.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Color,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

// Create adds a new category to the database
func (cs *CategoryStore) Create(name, color string) (*Category, error) {
	query := `
		INSERT INTO categories (name, color, created_at, updated_at)
		VALUES (?, ?, datetime('now'), datetime('now'))
	`
	
	result, err := cs.db.Exec(query, name, color)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// Retrieve the created category
	return cs.GetByID(int(id))
}

// Update modifies a category's name and color
func (cs *CategoryStore) Update(id int, name, color string) (*Category, error) {
	query := `
		UPDATE categories 
		SET name = ?, color = ?, updated_at = datetime('now')
		WHERE id = ?
	`
	
	result, err := cs.db.Exec(query, name, color, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("category not found")
	}

	// Return the updated category
	return cs.GetByID(id)
}

// Delete removes a category from the database
func (cs *CategoryStore) Delete(id int) error {
	// Check if category is in use
	var count int
	err := cs.db.QueryRow(`SELECT COUNT(*) FROM todos WHERE category_id = ?`, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check category usage: %w", err)
	}
	
	if count > 0 {
		return fmt.Errorf("category is in use by %d todos", count)
	}

	query := `DELETE FROM categories WHERE id = ?`
	
	result, err := cs.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}