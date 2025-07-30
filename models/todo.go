package models

import (
	"sync"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type TodoStore struct {
	mu    sync.RWMutex
	todos []Todo
	nextID int
}

func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos:  make([]Todo, 0),
		nextID: 1,
	}
}

func (ts *TodoStore) GetAll() []Todo {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	
	result := make([]Todo, len(ts.todos))
	copy(result, ts.todos)
	return result
}

func (ts *TodoStore) Create(title string) Todo {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	
	todo := Todo{
		ID:        ts.nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	
	ts.todos = append(ts.todos, todo)
	ts.nextID++
	
	return todo
}

func (ts *TodoStore) Toggle(id int) (*Todo, bool) {
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

func (ts *TodoStore) Delete(id int) bool {
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