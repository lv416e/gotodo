package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

// Initialize creates and returns a new database connection
func Initialize(dbPath string) (*DB, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open SQLite database
	sqlDB, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{sqlDB}

	// Run migrations
	if err := db.migrate(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// migrate runs the database migrations
func (db *DB) migrate() error {
	// Create categories table
	categoriesSchema := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(50) NOT NULL UNIQUE,
		color VARCHAR(7) DEFAULT '#007bff',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);

	-- Insert default categories
	INSERT OR IGNORE INTO categories (name, color) VALUES 
		('仕事', '#ff6b6b'),
		('プライベート', '#4dabf7'),
		('勉強', '#51cf66'),
		('その他', '#868e96');
	`

	if _, err := db.Exec(categoriesSchema); err != nil {
		return fmt.Errorf("failed to create categories table: %w", err)
	}

	// Check if category_id column exists in todos table
	var columnExists bool
	err := db.QueryRow(`
		SELECT COUNT(*) FROM pragma_table_info('todos') 
		WHERE name='category_id'
	`).Scan(&columnExists)
	if err != nil {
		return fmt.Errorf("failed to check column existence: %w", err)
	}

	// Add category_id to todos table if it doesn't exist
	if !columnExists {
		alterTableSQL := `
		ALTER TABLE todos ADD COLUMN category_id INTEGER REFERENCES categories(id);
		CREATE INDEX IF NOT EXISTS idx_todos_category_id ON todos(category_id);
		`
		if _, err := db.Exec(alterTableSQL); err != nil {
			return fmt.Errorf("failed to add category_id column: %w", err)
		}
	}

	// Check if priority column exists in todos table
	var priorityColumnExists bool
	err = db.QueryRow(`
		SELECT COUNT(*) FROM pragma_table_info('todos') 
		WHERE name='priority'
	`).Scan(&priorityColumnExists)
	if err != nil {
		return fmt.Errorf("failed to check priority column existence: %w", err)
	}

	// Add priority to todos table if it doesn't exist
	if !priorityColumnExists {
		alterPrioritySQL := `
		ALTER TABLE todos ADD COLUMN priority INTEGER DEFAULT 1;
		CREATE INDEX IF NOT EXISTS idx_todos_priority ON todos(priority);
		`
		if _, err := db.Exec(alterPrioritySQL); err != nil {
			return fmt.Errorf("failed to add priority column: %w", err)
		}
	}

	// Check if due_date column exists in todos table
	var dueDateColumnExists bool
	err = db.QueryRow(`
		SELECT COUNT(*) FROM pragma_table_info('todos') 
		WHERE name='due_date'
	`).Scan(&dueDateColumnExists)
	if err != nil {
		return fmt.Errorf("failed to check due_date column existence: %w", err)
	}

	// Add due_date to todos table if it doesn't exist
	if !dueDateColumnExists {
		alterDueDateSQL := `
		ALTER TABLE todos ADD COLUMN due_date DATETIME;
		CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date);
		`
		if _, err := db.Exec(alterDueDateSQL); err != nil {
			return fmt.Errorf("failed to add due_date column: %w", err)
		}
	}

	// Original todos table schema for new installations
	todosSchema := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title VARCHAR(255) NOT NULL,
		description TEXT DEFAULT '',
		category_id INTEGER REFERENCES categories(id),
		priority INTEGER DEFAULT 1,
		due_date DATETIME,
		completed BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_todos_completed ON todos(completed);
	CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at);
	CREATE INDEX IF NOT EXISTS idx_todos_category_id ON todos(category_id);
	CREATE INDEX IF NOT EXISTS idx_todos_priority ON todos(priority);
	CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date);
	`

	if _, err := db.Exec(todosSchema); err != nil {
		return fmt.Errorf("failed to execute todos schema: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}