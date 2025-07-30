class TodoApp {
    constructor() {
        this.todoForm = document.getElementById('todo-form');
        this.todoInput = document.getElementById('todo-input');
        this.todoList = document.getElementById('todo-list');
        this.totalCount = document.getElementById('total-count');
        this.completedCount = document.getElementById('completed-count');
        
        this.init();
    }
    
    init() {
        this.todoForm.addEventListener('submit', (e) => this.handleSubmit(e));
        this.loadTodos();
    }
    
    async handleSubmit(e) {
        e.preventDefault();
        
        const title = this.todoInput.value.trim();
        if (!title) return;
        
        try {
            await this.createTodo(title);
            this.todoInput.value = '';
            this.loadTodos();
        } catch (error) {
            console.error('Failed to create todo:', error);
            alert('TODOの作成に失敗しました');
        }
    }
    
    async createTodo(title) {
        const response = await fetch('/api/todos', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ title }),
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        return response.json();
    }
    
    async loadTodos() {
        try {
            const response = await fetch('/api/todos');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const todos = await response.json();
            this.renderTodos(todos);
            this.updateStats(todos);
        } catch (error) {
            console.error('Failed to load todos:', error);
            this.showError('TODOの読み込みに失敗しました');
        }
    }
    
    renderTodos(todos) {
        if (!todos || todos.length === 0) {
            this.todoList.innerHTML = `
                <div class="empty-state">
                    <p>TODOがありません</p>
                    <p>上のフォームから新しいTODOを追加してください</p>
                </div>
            `;
            return;
        }
        
        this.todoList.innerHTML = todos.map(todo => this.createTodoElement(todo)).join('');
    }
    
    createTodoElement(todo) {
        const createdAt = new Date(todo.created_at).toLocaleDateString('ja-JP');
        const completedClass = todo.completed ? 'completed' : '';
        
        return `
            <div class="todo-item ${completedClass}" data-id="${todo.id}">
                <input 
                    type="checkbox" 
                    class="todo-checkbox" 
                    ${todo.completed ? 'checked' : ''}
                    onchange="todoApp.toggleTodo(${todo.id})"
                >
                <span class="todo-text">${this.escapeHtml(todo.title)}</span>
                <span class="todo-date">${createdAt}</span>
                <button 
                    class="todo-delete" 
                    onclick="todoApp.deleteTodo(${todo.id})"
                >
                    削除
                </button>
            </div>
        `;
    }
    
    async toggleTodo(id) {
        try {
            const response = await fetch(`/api/todos/${id}/toggle`, {
                method: 'PUT',
            });
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            this.loadTodos();
        } catch (error) {
            console.error('Failed to toggle todo:', error);
            alert('TODOの更新に失敗しました');
        }
    }
    
    async deleteTodo(id) {
        if (!confirm('このTODOを削除しますか？')) {
            return;
        }
        
        try {
            const response = await fetch(`/api/todos/${id}`, {
                method: 'DELETE',
            });
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            this.loadTodos();
        } catch (error) {
            console.error('Failed to delete todo:', error);
            alert('TODOの削除に失敗しました');
        }
    }
    
    updateStats(todos) {
        const total = todos ? todos.length : 0;
        const completed = todos ? todos.filter(todo => todo.completed).length : 0;
        
        this.totalCount.textContent = total;
        this.completedCount.textContent = completed;
    }
    
    showError(message) {
        this.todoList.innerHTML = `
            <div class="empty-state">
                <p style="color: #ff6b6b;">エラー: ${message}</p>
                <p>ページを再読み込みしてください</p>
            </div>
        `;
    }
    
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// Initialize the app when the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.todoApp = new TodoApp();
});