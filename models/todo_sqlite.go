package models

import (
	"database/sql"
	"fmt"
	"time"

	"gotodo/database"
)

// Todo represents a TODO item with additional fields for SQLite
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TodoStore manages TODO items using SQLite database
type TodoStore struct {
	db *database.DB
}

// NewTodoStore creates a new TodoStore with SQLite backend
func NewTodoStore(db *database.DB) *TodoStore {
	return &TodoStore{
		db: db,
	}
}

// GetAll retrieves all TODO items from the database
func (ts *TodoStore) GetAll() ([]Todo, error) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at 
		FROM todos 
		ORDER BY created_at DESC
	`
	
	rows, err := ts.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return todos, nil
}

// Create adds a new TODO item to the database
func (ts *TodoStore) Create(title string) (*Todo, error) {
	query := `
		INSERT INTO todos (title, description, completed, created_at, updated_at)
		VALUES (?, '', FALSE, datetime('now'), datetime('now'))
	`
	
	result, err := ts.db.Exec(query, title)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// Retrieve the created todo
	return ts.GetByID(int(id))
}

// GetByID retrieves a specific TODO item by ID
func (ts *TodoStore) GetByID(id int) (*Todo, error) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at 
		FROM todos 
		WHERE id = ?
	`
	
	var todo Todo
	err := ts.db.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return &todo, nil
}

// Toggle switches the completion status of a TODO item
func (ts *TodoStore) Toggle(id int) (*Todo, error) {
	query := `
		UPDATE todos 
		SET completed = NOT completed, updated_at = datetime('now')
		WHERE id = ?
	`
	
	result, err := ts.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to toggle todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("todo not found")
	}

	// Return the updated todo
	return ts.GetByID(id)
}

// Update modifies a TODO item's title and description
func (ts *TodoStore) Update(id int, title, description string) (*Todo, error) {
	query := `
		UPDATE todos 
		SET title = ?, description = ?, updated_at = datetime('now')
		WHERE id = ?
	`
	
	result, err := ts.db.Exec(query, title, description, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("todo not found")
	}

	// Return the updated todo
	return ts.GetByID(id)
}

// Delete removes a TODO item from the database
func (ts *TodoStore) Delete(id int) error {
	query := `DELETE FROM todos WHERE id = ?`
	
	result, err := ts.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}