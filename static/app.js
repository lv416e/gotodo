class TodoApp {
    constructor() {
        this.todoForm = document.getElementById('todo-form');
        this.todoInput = document.getElementById('todo-input');
        this.todoList = document.getElementById('todo-list');
        this.totalCount = document.getElementById('total-count');
        this.completedCount = document.getElementById('completed-count');
        this.categorySelect = null; // Will be initialized after categories load
        this.categories = [];
        
        this.init();
    }
    
    init() {
        this.todoForm.addEventListener('submit', (e) => this.handleSubmit(e));
        this.loadCategories();
        this.loadTodos();
    }
    
    async loadCategories() {
        try {
            const response = await fetch('/api/categories');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            this.categories = await response.json();
            this.updateCategorySelect();
        } catch (error) {
            console.error('Failed to load categories:', error);
        }
    }
    
    updateCategorySelect() {
        const categorySelectContainer = document.getElementById('category-select-container');
        if (!categorySelectContainer) return;
        
        const selectHtml = `
            <select id="category-select" class="category-select">
                <option value="">カテゴリなし</option>
                ${this.categories.map(cat => 
                    `<option value="${cat.id}" style="color: ${cat.color}">${this.escapeHtml(cat.name)}</option>`
                ).join('')}
            </select>
        `;
        
        categorySelectContainer.innerHTML = selectHtml;
        this.categorySelect = document.getElementById('category-select');
    }
    
    async handleSubmit(e) {
        e.preventDefault();
        
        const title = this.todoInput.value.trim();
        if (!title) return;
        
        const categoryId = this.categorySelect ? 
            (this.categorySelect.value ? parseInt(this.categorySelect.value) : null) : 
            null;
        
        try {
            await this.createTodo(title, categoryId);
            this.todoInput.value = '';
            this.loadTodos();
        } catch (error) {
            console.error('Failed to create todo:', error);
            alert('TODOの作成に失敗しました');
        }
    }
    
    async createTodo(title, categoryId) {
        const body = { title };
        if (categoryId !== null) {
            body.category_id = categoryId;
        }
        
        const response = await fetch('/api/todos', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(body),
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
        const updatedAt = new Date(todo.updated_at).toLocaleDateString('ja-JP');
        const completedClass = todo.completed ? 'completed' : '';
        
        const categoryBadge = todo.category ? 
            `<span class="category-badge" style="background-color: ${todo.category.color}">${this.escapeHtml(todo.category.name)}</span>` : 
            '';
        
        return `
            <div class="todo-item ${completedClass}" data-id="${todo.id}">
                <input 
                    type="checkbox" 
                    class="todo-checkbox" 
                    ${todo.completed ? 'checked' : ''}
                    onchange="todoApp.toggleTodo(${todo.id})"
                >
                <div class="todo-content">
                    <div class="todo-header">
                        <div class="todo-title" onclick="todoApp.editTodo(${todo.id})">
                            ${this.escapeHtml(todo.title)}
                        </div>
                        ${categoryBadge}
                    </div>
                    ${todo.description ? `<div class="todo-description">${this.escapeHtml(todo.description)}</div>` : ''}
                    <div class="todo-dates">
                        <span class="todo-date">作成: ${createdAt}</span>
                        ${createdAt !== updatedAt ? `<span class="todo-date">更新: ${updatedAt}</span>` : ''}
                    </div>
                </div>
                <div class="todo-actions">
                    <button 
                        class="todo-edit" 
                        onclick="todoApp.editTodo(${todo.id})"
                        title="編集"
                    >
                        ✏️
                    </button>
                    <button 
                        class="todo-delete" 
                        onclick="todoApp.deleteTodo(${todo.id})"
                        title="削除"
                    >
                        🗑️
                    </button>
                </div>
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
    
    async editTodo(id) {
        // Find the todo item
        const todos = await this.getTodos();
        const todo = todos.find(t => t.id === id);
        
        if (!todo) {
            alert('TODOが見つかりません');
            return;
        }
        
        // Show edit modal
        this.showEditModal(todo);
    }
    
    async getTodos() {
        try {
            const response = await fetch('/api/todos');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        } catch (error) {
            console.error('Failed to get todos:', error);
            return [];
        }
    }
    
    showEditModal(todo) {
        const categoryOptions = this.categories.map(cat => 
            `<option value="${cat.id}" ${todo.category_id === cat.id ? 'selected' : ''} style="color: ${cat.color}">${this.escapeHtml(cat.name)}</option>`
        ).join('');
        
        // Create modal HTML
        const modalHtml = `
            <div id="edit-modal" class="modal">
                <div class="modal-content">
                    <div class="modal-header">
                        <h3>TODOを編集</h3>
                        <button class="modal-close" onclick="todoApp.closeEditModal()">&times;</button>
                    </div>
                    <form id="edit-form">
                        <div class="form-group">
                            <label for="edit-title">タイトル</label>
                            <input type="text" id="edit-title" value="${this.escapeHtml(todo.title)}" required>
                        </div>
                        <div class="form-group">
                            <label for="edit-description">説明</label>
                            <textarea id="edit-description" rows="3" placeholder="説明を入力（任意）">${this.escapeHtml(todo.description)}</textarea>
                        </div>
                        <div class="form-group">
                            <label for="edit-category">カテゴリ</label>
                            <select id="edit-category" class="category-select">
                                <option value="" ${!todo.category_id ? 'selected' : ''}>カテゴリなし</option>
                                ${categoryOptions}
                            </select>
                        </div>
                        <div class="form-actions">
                            <button type="button" onclick="todoApp.closeEditModal()">キャンセル</button>
                            <button type="submit">保存</button>
                        </div>
                    </form>
                </div>
            </div>
        `;
        
        // Add modal to body
        document.body.insertAdjacentHTML('beforeend', modalHtml);
        
        // Add event listener for form submission
        document.getElementById('edit-form').addEventListener('submit', (e) => {
            e.preventDefault();
            this.saveEdit(todo.id);
        });
        
        // Focus on title input
        document.getElementById('edit-title').focus();
    }
    
    async saveEdit(id) {
        const title = document.getElementById('edit-title').value.trim();
        const description = document.getElementById('edit-description').value.trim();
        const categorySelect = document.getElementById('edit-category');
        const categoryId = categorySelect.value ? parseInt(categorySelect.value) : null;
        
        if (!title) {
            alert('タイトルは必須です');
            return;
        }
        
        const body = { title, description };
        if (categoryId !== null) {
            body.category_id = categoryId;
        }
        
        try {
            const response = await fetch(`/api/todos/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(body),
            });
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            this.closeEditModal();
            this.loadTodos();
        } catch (error) {
            console.error('Failed to update todo:', error);
            alert('TODOの更新に失敗しました');
        }
    }
    
    closeEditModal() {
        const modal = document.getElementById('edit-modal');
        if (modal) {
            modal.remove();
        }
    }
    
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text || '';
        return div.innerHTML;
    }
}

// Initialize the app when the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.todoApp = new TodoApp();
});