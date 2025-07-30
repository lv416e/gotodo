// This file contains the old in-memory implementation
// It's kept for reference and can be removed once SQLite migration is complete
package models

import (
	"sync"
	"time"
)

// OldTodo represents the original in-memory TODO structure
type OldTodo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// OldTodoStore manages TODO items in memory (deprecated)
type OldTodoStore struct {
	mu     sync.RWMutex
	todos  []OldTodo
	nextID int
}

// NewOldTodoStore creates a new in-memory TodoStore (deprecated)
func NewOldTodoStore() *OldTodoStore {
	return &OldTodoStore{
		todos:  make([]OldTodo, 0),
		nextID: 1,
	}
}

func (ts *OldTodoStore) GetAll() []OldTodo {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	
	result := make([]OldTodo, len(ts.todos))
	copy(result, ts.todos)
	return result
}

func (ts *OldTodoStore) Create(title string) OldTodo {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	
	todo := OldTodo{
		ID:        ts.nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	
	ts.todos = append(ts.todos, todo)
	ts.nextID++
	
	return todo
}

func (ts *OldTodoStore) Toggle(id int) (*OldTodo, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	
	for i := range ts.todos {
		if ts.todos[i].ID == id {
			ts.todos[i].Completed = !ts.todos[i].Completed
			return &ts.todos[i], true
		}
	}
	
	return nil, false
}

func (ts *OldTodoStore) Delete(id int) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	
	for i, todo := range ts.todos {
		if todo.ID == id {
			ts.todos = append(ts.todos[:i], ts.todos[i+1:]...)
			return true
		}
	}
	
	return false
}