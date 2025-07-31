class TodoApp {
    constructor() {
        this.todoForm = document.getElementById('todo-form');
        this.todoInput = document.getElementById('todo-input');
        this.todoList = document.getElementById('todo-list');
        this.totalCount = document.getElementById('total-count');
        this.completedCount = document.getElementById('completed-count');
        this.categorySelect = null; // Will be initialized after categories load
        this.categoryFilter = null; // Will be initialized after categories load
        this.prioritySelect = document.getElementById('priority-select');
        this.priorityFilter = document.getElementById('priority-filter');
        this.dueDateInput = document.getElementById('due-date-input');
        this.manageCategoriesBtn = document.getElementById('manage-categories');
        this.overdueWarning = document.getElementById('overdue-warning');
        this.overdueCount = document.getElementById('overdue-count');
        this.searchInput = document.getElementById('search-input');
        this.clearSearchBtn = document.getElementById('clear-search');
        this.statusFilter = document.getElementById('status-filter');
        this.dueDateFilter = document.getElementById('due-date-filter');
        this.categoryProgress = document.getElementById('category-progress');
        this.categoryProgressList = document.getElementById('category-progress-list');
        this.categories = [];
        this.allTodos = []; // Store all todos for filtering
        this.currentCategoryFilter = ''; // Current filter category ID
        this.currentPriorityFilter = ''; // Current filter priority
        this.currentSearchQuery = ''; // Current search query
        this.currentStatusFilter = ''; // Current status filter
        this.currentDueDateFilter = ''; // Current due date filter
        
        this.init();
    }
    
    init() {
        this.todoForm.addEventListener('submit', (e) => this.handleSubmit(e));
        this.manageCategoriesBtn.addEventListener('click', () => this.showCategoryManager());
        this.priorityFilter.addEventListener('change', (e) => {
            this.currentPriorityFilter = e.target.value;
            this.filterTodos();
        });
        
        this.statusFilter.addEventListener('change', (e) => {
            this.currentStatusFilter = e.target.value;
            this.filterTodos();
        });
        
        this.dueDateFilter.addEventListener('change', (e) => {
            this.currentDueDateFilter = e.target.value;
            this.filterTodos();
        });
        
        // Search functionality
        let searchTimeout;
        this.searchInput.addEventListener('input', (e) => {
            clearTimeout(searchTimeout);
            searchTimeout = setTimeout(() => {
                this.currentSearchQuery = e.target.value.trim();
                this.handleSearch();
            }, 300); // Debounce search for 300ms
        });
        
        this.clearSearchBtn.addEventListener('click', () => {
            this.clearSearch();
        });
        
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
        
        // Update filter dropdown
        this.updateCategoryFilter();
    }
    
    updateCategoryFilter() {
        this.categoryFilter = document.getElementById('category-filter');
        if (!this.categoryFilter) return;
        
        const filterOptions = this.categories.map(cat => 
            `<option value="${cat.id}" style="color: ${cat.color}">${this.escapeHtml(cat.name)}</option>`
        ).join('');
        
        this.categoryFilter.innerHTML = `
            <option value="">すべて</option>
            <option value="null">カテゴリなし</option>
            ${filterOptions}
        `;
        
        this.categoryFilter.addEventListener('change', (e) => {
            this.currentCategoryFilter = e.target.value;
            this.filterTodos();
        });
    }
    
    async handleSubmit(e) {
        e.preventDefault();
        
        const title = this.todoInput.value.trim();
        if (!title) return;
        
        const categoryId = this.categorySelect ? 
            (this.categorySelect.value ? parseInt(this.categorySelect.value) : null) : 
            null;
        const priority = this.prioritySelect.value ? parseInt(this.prioritySelect.value) : 1;
        const dueDate = this.dueDateInput.value || null;
        
        try {
            await this.createTodo(title, categoryId, priority, dueDate);
            this.todoInput.value = '';
            this.prioritySelect.value = '1'; // Reset to low priority
            this.dueDateInput.value = ''; // Reset due date
            this.loadTodos();
        } catch (error) {
            console.error('Failed to create todo:', error);
            alert('TODOの作成に失敗しました');
        }
    }
    
    async createTodo(title, categoryId, priority, dueDate) {
        const body = { 
            title,
            description: '' // Add description field for API compatibility
        };
        if (categoryId !== null) {
            body.category_id = categoryId;
        }
        if (priority !== null && priority !== undefined) {
            body.priority = priority;
        }
        if (dueDate) {
            body.due_date = dueDate;
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
            let url = '/api/todos';
            if (this.currentSearchQuery) {
                url += `?search=${encodeURIComponent(this.currentSearchQuery)}`;
            }
            
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const todos = await response.json();
            this.allTodos = todos; // Store all todos
            this.filterTodos(); // Apply current filter
        } catch (error) {
            console.error('Failed to load todos:', error);
            this.showError('TODOの読み込みに失敗しました');
        }
    }
    
    async handleSearch() {
        // Show/hide clear button
        if (this.currentSearchQuery) {
            this.clearSearchBtn.style.display = 'block';
        } else {
            this.clearSearchBtn.style.display = 'none';
        }
        
        // Reload todos with search query
        await this.loadTodos();
    }
    
    clearSearch() {
        this.currentSearchQuery = '';
        this.searchInput.value = '';
        this.clearSearchBtn.style.display = 'none';
        this.loadTodos();
    }
    
    filterTodos() {
        let filteredTodos = this.allTodos;
        
        // Apply category filter
        if (this.currentCategoryFilter === 'null') {
            // Show only todos without category
            filteredTodos = filteredTodos.filter(todo => !todo.category_id);
        } else if (this.currentCategoryFilter && this.currentCategoryFilter !== '') {
            // Show only todos with specific category
            const categoryId = parseInt(this.currentCategoryFilter);
            filteredTodos = filteredTodos.filter(todo => todo.category_id === categoryId);
        }
        
        // Apply priority filter
        if (this.currentPriorityFilter && this.currentPriorityFilter !== '') {
            const priority = parseInt(this.currentPriorityFilter);
            filteredTodos = filteredTodos.filter(todo => todo.priority === priority);
        }
        
        // Apply status filter
        if (this.currentStatusFilter && this.currentStatusFilter !== '') {
            if (this.currentStatusFilter === 'completed') {
                filteredTodos = filteredTodos.filter(todo => todo.completed);
            } else if (this.currentStatusFilter === 'incomplete') {
                filteredTodos = filteredTodos.filter(todo => !todo.completed);
            }
        }
        
        // Apply due date filter
        if (this.currentDueDateFilter && this.currentDueDateFilter !== '') {
            const now = new Date();
            const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
            const tomorrow = new Date(today);
            tomorrow.setDate(tomorrow.getDate() + 1);
            const weekEnd = new Date(today);
            weekEnd.setDate(weekEnd.getDate() + 7);
            
            filteredTodos = filteredTodos.filter(todo => {
                if (this.currentDueDateFilter === 'no_due_date') {
                    return !todo.due_date;
                }
                
                if (!todo.due_date) return false;
                
                const dueDate = new Date(todo.due_date);
                const dueDateOnly = new Date(dueDate.getFullYear(), dueDate.getMonth(), dueDate.getDate());
                
                switch (this.currentDueDateFilter) {
                    case 'overdue':
                        return dueDate < now;
                    case 'today':
                        return dueDateOnly.getTime() === today.getTime();
                    case 'tomorrow':
                        return dueDateOnly.getTime() === tomorrow.getTime();
                    case 'this_week':
                        return dueDateOnly >= today && dueDateOnly <= weekEnd;
                    default:
                        return true;
                }
            });
        }
        
        this.renderTodos(filteredTodos);
        this.updateStats(filteredTodos);
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
        
        // Add overdue class for visual emphasis
        const overdueClass = (!todo.completed && todo.due_date && this.isOverdue(todo.due_date)) ? 'overdue-item' : '';
        
        const categoryBadge = todo.category ? 
            `<span class="category-badge" style="background-color: ${todo.category.color}">${this.escapeHtml(todo.category.name)}</span>` : 
            '';
        
        // Priority badge with appropriate styling
        const priorityInfo = this.getPriorityInfo(todo.priority);
        const priorityBadge = `<span class="priority-badge priority-${priorityInfo.level}" title="優先度: ${priorityInfo.name}">${priorityInfo.icon}</span>`;
        
        return `
            <div class="todo-item ${completedClass} ${overdueClass}" data-id="${todo.id}">
                <input 
                    type="checkbox" 
                    class="todo-checkbox" 
                    ${todo.completed ? 'checked' : ''}
                    onchange="todoApp.toggleTodo(${todo.id})"
                >
                <div class="todo-content">
                    <div class="todo-header">
                        ${priorityBadge}
                        <div class="todo-title" onclick="todoApp.editTodo(${todo.id})">
                            ${this.escapeHtml(todo.title)}
                        </div>
                        ${categoryBadge}
                    </div>
                    ${todo.description ? `<div class="todo-description">${this.escapeHtml(todo.description)}</div>` : ''}
                    <div class="todo-dates">
                        ${todo.due_date ? `<span class="todo-due-date ${this.getDueDateClass(todo.due_date, todo.completed)}">期限: ${this.formatDueDate(todo.due_date)}</span>` : ''}
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
        
        // Count overdue items (only among incomplete todos)
        const overdue = todos ? todos.filter(todo => 
            !todo.completed && 
            todo.due_date && 
            this.isOverdue(todo.due_date)
        ).length : 0;
        
        this.totalCount.textContent = total;
        this.completedCount.textContent = completed;
        
        // Show/hide overdue warning
        if (overdue > 0) {
            this.overdueCount.textContent = overdue;
            this.overdueWarning.style.display = 'block';
        } else {
            this.overdueWarning.style.display = 'none';
        }
        
        // Update category progress
        this.updateCategoryProgress();
    }
    
    updateCategoryProgress() {
        if (!this.categories || this.categories.length === 0) {
            this.categoryProgress.style.display = 'none';
            return;
        }
        
        // Calculate progress for each category
        const categoryStats = this.categories.map(category => {
            const categoryTodos = this.allTodos.filter(todo => todo.category_id === category.id);
            const completedTodos = categoryTodos.filter(todo => todo.completed);
            const total = categoryTodos.length;
            const completed = completedTodos.length;
            const percentage = total > 0 ? Math.round((completed / total) * 100) : 0;
            
            return {
                ...category,
                total,
                completed,
                percentage
            };
        });
        
        // Add "カテゴリなし" category
        const uncategorizedTodos = this.allTodos.filter(todo => !todo.category_id);
        if (uncategorizedTodos.length > 0) {
            const completedUncategorized = uncategorizedTodos.filter(todo => todo.completed);
            categoryStats.push({
                id: null,
                name: 'カテゴリなし',
                color: '#6c757d',
                total: uncategorizedTodos.length,
                completed: completedUncategorized.length,
                percentage: Math.round((completedUncategorized.length / uncategorizedTodos.length) * 100)
            });
        }
        
        // Filter out categories with no todos
        const activeCategoryStats = categoryStats.filter(stat => stat.total > 0);
        
        if (activeCategoryStats.length === 0) {
            this.categoryProgress.style.display = 'none';
            return;
        }
        
        // Show category progress section
        this.categoryProgress.style.display = 'block';
        
        // Generate HTML for category progress
        const progressHTML = activeCategoryStats.map(stat => `
            <div class="category-progress-item">
                <div class="category-info">
                    <span class="category-badge" style="background-color: ${stat.color}">${this.escapeHtml(stat.name)}</span>
                </div>
                <div class="category-progress-bar">
                    <div class="category-progress-fill" style="width: ${stat.percentage}%; background-color: ${stat.color}"></div>
                </div>
                <div class="category-progress-text">
                    ${stat.completed}/${stat.total} (${stat.percentage}%)
                </div>
            </div>
        `).join('');
        
        this.categoryProgressList.innerHTML = progressHTML;
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
                        <div class="form-group">
                            <label for="edit-priority">優先度</label>
                            <select id="edit-priority" class="priority-select">
                                <option value="1" ${todo.priority === 1 ? 'selected' : ''}>低</option>
                                <option value="2" ${todo.priority === 2 ? 'selected' : ''}>中</option>
                                <option value="3" ${todo.priority === 3 ? 'selected' : ''}>高</option>
                            </select>
                        </div>
                        <div class="form-group">
                            <label for="edit-due-date">期限</label>
                            <input type="datetime-local" id="edit-due-date" value="${todo.due_date ? this.formatDateForInput(todo.due_date) : ''}" title="期限を設定（任意）">
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
        const prioritySelect = document.getElementById('edit-priority');
        const priority = prioritySelect.value ? parseInt(prioritySelect.value) : 1;
        const dueDateInput = document.getElementById('edit-due-date');
        const dueDate = dueDateInput.value || null;
        
        if (!title) {
            alert('タイトルは必須です');
            return;
        }
        
        const body = { title, description };
        if (categoryId !== null) {
            body.category_id = categoryId;
        }
        if (priority !== null && priority !== undefined) {
            body.priority = priority;
        }
        if (dueDate) {
            body.due_date = dueDate;
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
    
    showCategoryManager() {
        const modalHtml = `
            <div id="category-manager-modal" class="modal">
                <div class="modal-content category-manager">
                    <div class="modal-header">
                        <h3>カテゴリ管理</h3>
                        <button class="modal-close" onclick="todoApp.closeCategoryManager()">&times;</button>
                    </div>
                    <div class="category-manager-content">
                        <div class="add-category-section">
                            <h4>新しいカテゴリを追加</h4>
                            <form id="add-category-form" class="add-category-form">
                                <div class="form-row">
                                    <input type="text" id="new-category-name" placeholder="カテゴリ名" required>
                                    <input type="color" id="new-category-color" value="#007bff">
                                    <button type="submit">追加</button>
                                </div>
                            </form>
                        </div>
                        <div class="categories-list-section">
                            <h4>既存のカテゴリ</h4>
                            <div id="categories-list" class="categories-list">
                                <!-- Categories will be populated here -->
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
        
        document.body.insertAdjacentHTML('beforeend', modalHtml);
        this.loadCategoriesInManager();
        
        // Add form submit handler
        document.getElementById('add-category-form').addEventListener('submit', (e) => {
            e.preventDefault();
            this.addCategory();
        });
    }
    
    async loadCategoriesInManager() {
        const categoriesList = document.getElementById('categories-list');
        if (!categoriesList) return;
        
        const categoriesHtml = this.categories.map(cat => `
            <div class="category-item" data-id="${cat.id}">
                <div class="category-info">
                    <span class="category-badge" style="background-color: ${cat.color}">${this.escapeHtml(cat.name)}</span>
                    <span class="category-details">${cat.color}</span>
                </div>
                <div class="category-actions">
                    <button class="edit-category-btn" onclick="todoApp.editCategory(${cat.id})">編集</button>
                    <button class="delete-category-btn" onclick="todoApp.deleteCategory(${cat.id})">削除</button>
                </div>
            </div>
        `).join('');
        
        categoriesList.innerHTML = categoriesHtml;
    }
    
    async addCategory() {
        const name = document.getElementById('new-category-name').value.trim();
        const color = document.getElementById('new-category-color').value;
        
        if (!name) {
            alert('カテゴリ名を入力してください');
            return;
        }
        
        try {
            const response = await fetch('/api/categories', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ name, color }),
            });
            
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText);
            }
            
            // Reset form
            document.getElementById('new-category-name').value = '';
            document.getElementById('new-category-color').value = '#007bff';
            
            // Reload categories
            await this.loadCategories();
            this.loadCategoriesInManager();
            
        } catch (error) {
            console.error('Failed to add category:', error);
            alert('カテゴリの追加に失敗しました: ' + error.message);
        }
    }
    
    async deleteCategory(id) {
        if (!confirm('このカテゴリを削除しますか？使用中の場合は削除できません。')) {
            return;
        }
        
        try {
            const response = await fetch(`/api/categories/${id}`, {
                method: 'DELETE',
            });
            
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText);
            }
            
            // Reload categories
            await this.loadCategories();
            this.loadCategoriesInManager();
            this.loadTodos(); // Reload todos to update display
            
        } catch (error) {
            console.error('Failed to delete category:', error);
            alert('カテゴリの削除に失敗しました: ' + error.message);
        }
    }
    
    async editCategory(id) {
        const category = this.categories.find(cat => cat.id === id);
        if (!category) return;
        
        const newName = prompt('カテゴリ名を編集:', category.name);
        if (!newName || newName.trim() === '') return;
        
        const newColor = prompt('カテゴリの色を編集 (例: #ff0000):', category.color);
        if (!newColor) return;
        
        try {
            const response = await fetch(`/api/categories/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ name: newName.trim(), color: newColor }),
            });
            
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText);
            }
            
            // Reload categories
            await this.loadCategories();
            this.loadCategoriesInManager();
            this.loadTodos(); // Reload todos to update display
            
        } catch (error) {
            console.error('Failed to edit category:', error);
            alert('カテゴリの編集に失敗しました: ' + error.message);
        }
    }
    
    closeCategoryManager() {
        const modal = document.getElementById('category-manager-modal');
        if (modal) {
            modal.remove();
        }
    }
    
    getPriorityInfo(priority) {
        switch (priority) {
            case 3:
                return { level: 'high', name: '高', icon: '🔴' };
            case 2:
                return { level: 'medium', name: '中', icon: '🟡' };
            case 1:
            default:
                return { level: 'low', name: '低', icon: '⚪' };
        }
    }
    
    formatDueDate(dueDateString) {
        const dueDate = new Date(dueDateString);
        const now = new Date();
        const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
        const dueDateOnly = new Date(dueDate.getFullYear(), dueDate.getMonth(), dueDate.getDate());
        
        const diffTime = dueDateOnly - today;
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        
        const dateStr = dueDate.toLocaleDateString('ja-JP', {
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
        
        if (diffDays === 0) {
            return `${dateStr} (今日)`;
        } else if (diffDays === 1) {
            return `${dateStr} (明日)`;
        } else if (diffDays === -1) {
            return `${dateStr} (昨日)`;
        } else if (diffDays > 0) {
            return `${dateStr} (${diffDays}日後)`;
        } else {
            return `${dateStr} (${Math.abs(diffDays)}日前)`;
        }
    }
    
    getDueDateClass(dueDateString, completed) {
        if (completed) return 'completed';
        
        const dueDate = new Date(dueDateString);
        const now = new Date();
        const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
        const dueDateOnly = new Date(dueDate.getFullYear(), dueDate.getMonth(), dueDate.getDate());
        
        const diffTime = dueDateOnly - today;
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        
        if (diffDays < 0) {
            return 'overdue';
        } else if (diffDays === 0) {
            return 'today';
        } else if (diffDays === 1) {
            return 'tomorrow';
        } else if (diffDays <= 3) {
            return 'soon';
        } else {
            return 'future';
        }
    }
    
    formatDateForInput(dateString) {
        const date = new Date(dateString);
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        return `${year}-${month}-${day}T${hours}:${minutes}`;
    }
    
    isOverdue(dueDateString) {
        const dueDate = new Date(dueDateString);
        const now = new Date();
        return dueDate < now;
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