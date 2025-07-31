package models

import (
	"database/sql"
	"fmt"
	"time"

	"gotodo/database"
)

// Todo represents a TODO item with additional fields for SQLite
type Todo struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CategoryID  *int       `json:"category_id"`
	Category    *Category  `json:"category,omitempty"`
	Priority    int        `json:"priority"` // 1:低, 2:中, 3:高
	DueDate     *time.Time `json:"due_date"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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
		SELECT 
			t.id, t.title, t.description, t.category_id, t.priority, t.due_date,
			t.completed, t.created_at, t.updated_at,
			c.id, c.name, c.color
		FROM todos t
		LEFT JOIN categories c ON t.category_id = c.id
		ORDER BY 
			t.completed ASC,
			CASE 
				WHEN t.due_date IS NULL THEN 2
				WHEN t.due_date < datetime('now') THEN 0
				ELSE 1
			END ASC,
			t.due_date ASC,
			t.priority DESC,
			t.created_at DESC
	`
	
	rows, err := ts.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		var categoryID, categoryIDJoin sql.NullInt64
		var categoryName, categoryColor sql.NullString
		var dueDate sql.NullTime
		
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&categoryID,
			&todo.Priority,
			&dueDate,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&categoryIDJoin,
			&categoryName,
			&categoryColor,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		
		// Handle due date
		if dueDate.Valid {
			todo.DueDate = &dueDate.Time
		}
		
		// Handle category
		if categoryID.Valid {
			id := int(categoryID.Int64)
			todo.CategoryID = &id
			
			if categoryIDJoin.Valid && categoryName.Valid && categoryColor.Valid {
				todo.Category = &Category{
					ID:    int(categoryIDJoin.Int64),
					Name:  categoryName.String,
					Color: categoryColor.String,
				}
			}
		}
		
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return todos, nil
}

// Search retrieves TODO items that match the search query
func (ts *TodoStore) Search(query string) ([]Todo, error) {
	searchQuery := `
		SELECT 
			t.id, t.title, t.description, t.category_id, t.priority, t.due_date,
			t.completed, t.created_at, t.updated_at,
			c.id, c.name, c.color
		FROM todos t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.title LIKE ? OR t.description LIKE ?
		ORDER BY 
			t.completed ASC,
			CASE 
				WHEN t.due_date IS NULL THEN 2
				WHEN t.due_date < datetime('now') THEN 0
				ELSE 1
			END ASC,
			t.due_date ASC,
			t.priority DESC,
			t.created_at DESC
	`
	
	searchPattern := "%" + query + "%"
	rows, err := ts.db.Query(searchQuery, searchPattern, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search todos: %w", err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		var categoryID, categoryIDJoin sql.NullInt64
		var categoryName, categoryColor sql.NullString
		var dueDate sql.NullTime
		
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&categoryID,
			&todo.Priority,
			&dueDate,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&categoryIDJoin,
			&categoryName,
			&categoryColor,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		
		// Handle due date
		if dueDate.Valid {
			todo.DueDate = &dueDate.Time
		}
		
		// Handle category
		if categoryID.Valid {
			id := int(categoryID.Int64)
			todo.CategoryID = &id
			
			if categoryIDJoin.Valid && categoryName.Valid && categoryColor.Valid {
				todo.Category = &Category{
					ID:    int(categoryIDJoin.Int64),
					Name:  categoryName.String,
					Color: categoryColor.String,
				}
			}
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
		INSERT INTO todos (title, description, category_id, priority, completed, created_at, updated_at)
		VALUES (?, '', NULL, 1, FALSE, datetime('now'), datetime('now'))
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

// CreateWithCategory adds a new TODO item with a category
func (ts *TodoStore) CreateWithCategory(title string, categoryID *int) (*Todo, error) {
	query := `
		INSERT INTO todos (title, description, category_id, priority, completed, created_at, updated_at)
		VALUES (?, '', ?, 1, FALSE, datetime('now'), datetime('now'))
	`
	
	result, err := ts.db.Exec(query, title, categoryID)
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

// CreateWithCategoryAndPriority adds a new TODO item with category and priority
func (ts *TodoStore) CreateWithCategoryAndPriority(title string, categoryID *int, priority int) (*Todo, error) {
	// Validate priority range
	if priority < 1 || priority > 3 {
		priority = 1 // Default to low priority if invalid
	}
	
	query := `
		INSERT INTO todos (title, description, category_id, priority, due_date, completed, created_at, updated_at)
		VALUES (?, '', ?, ?, NULL, FALSE, datetime('now'), datetime('now'))
	`
	
	result, err := ts.db.Exec(query, title, categoryID, priority)
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

// CreateFull adds a new TODO item with all fields
func (ts *TodoStore) CreateFull(title, description string, categoryID *int, priority int, dueDate *time.Time) (*Todo, error) {
	// Validate priority range
	if priority < 1 || priority > 3 {
		priority = 1 // Default to low priority if invalid
	}
	
	query := `
		INSERT INTO todos (title, description, category_id, priority, due_date, completed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, FALSE, datetime('now'), datetime('now'))
	`
	
	result, err := ts.db.Exec(query, title, description, categoryID, priority, dueDate)
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
		SELECT 
			t.id, t.title, t.description, t.category_id, t.priority, t.due_date,
			t.completed, t.created_at, t.updated_at,
			c.id, c.name, c.color
		FROM todos t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.id = ?
	`
	
	var todo Todo
	var categoryID, categoryIDJoin sql.NullInt64
	var categoryName, categoryColor sql.NullString
	var dueDate sql.NullTime
	
	err := ts.db.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&categoryID,
		&todo.Priority,
		&dueDate,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&categoryIDJoin,
		&categoryName,
		&categoryColor,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}
	
	// Handle due date
	if dueDate.Valid {
		todo.DueDate = &dueDate.Time
	}
	
	// Handle category
	if categoryID.Valid {
		id := int(categoryID.Int64)
		todo.CategoryID = &id
		
		if categoryIDJoin.Valid && categoryName.Valid && categoryColor.Valid {
			todo.Category = &Category{
				ID:    int(categoryIDJoin.Int64),
				Name:  categoryName.String,
				Color: categoryColor.String,
			}
		}
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

// Update modifies a TODO item's title, description, category, priority, and due date
func (ts *TodoStore) Update(id int, title, description string, categoryID *int, priority *int, dueDate *time.Time) (*Todo, error) {
	// Validate priority if provided
	if priority != nil && (*priority < 1 || *priority > 3) {
		defaultPriority := 1
		priority = &defaultPriority
	}
	
	query := `
		UPDATE todos 
		SET title = ?, description = ?, category_id = ?, priority = COALESCE(?, priority), due_date = ?, updated_at = datetime('now')
		WHERE id = ?
	`
	
	result, err := ts.db.Exec(query, title, description, categoryID, priority, dueDate, id)
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