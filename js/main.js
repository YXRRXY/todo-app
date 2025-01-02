const API_BASE_URL = 'http://localhost:8888';
let currentPage = 1;
let pageSize = 10;
let totalItems = 0;

// 添加状态过滤功能
let currentStatus = null;

function addStatusFilter() {
    const filterDiv = document.createElement('div');
    filterDiv.className = 'status-filter';
    filterDiv.innerHTML = `
        <select id="statusFilter" onchange="handleStatusFilter()">
            <option value="">全部</option>
            <option value="0">未完成</option>
            <option value="1">已完成</option>
        </select>
    `;
    document.querySelector('.search-bar').appendChild(filterDiv);
}

function handleStatusFilter() {
    const status = document.getElementById('statusFilter').value;
    currentStatus = status === '' ? null : parseInt(status);
    currentPage = 1;
    fetchTodos();
}

// 添加批量操作功能
let selectedTodos = new Set();

// 页面切换函数
function showLogin() {
    document.getElementById('loginForm').style.display = 'block';
    document.getElementById('registerForm').style.display = 'none';
    document.getElementById('todoApp').style.display = 'none';
}

function showRegister() {
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('registerForm').style.display = 'block';
    document.getElementById('todoApp').style.display = 'none';
}

function showTodoApp() {
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('registerForm').style.display = 'none';
    document.getElementById('todoApp').style.display = 'block';
    fetchTodos();
}

// API 请求函数
async function handleLogin(event) {
    event.preventDefault();
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const response = await fetch(`${API_BASE_URL}/user/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });

        const data = await response.json();
        if (response.ok) {
            localStorage.setItem('token', data.token);
            showTodoApp();
        } else {
            alert('登录失败：' + data.message);
        }
    } catch (error) {
        alert('登录失败：' + error.message);
    }
}

async function handleRegister(event) {
    event.preventDefault();
    const username = document.getElementById('registerUsername').value;
    const password = document.getElementById('registerPassword').value;

    try {
        const response = await fetch(`${API_BASE_URL}/user/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });

        if (response.ok) {
            alert('注册成功！');
            showLogin();
        } else {
            const data = await response.json();
            alert('注册失败：' + data.message);
        }
    } catch (error) {
        alert('注册失败：' + error.message);
    }
}

// 添加待办事项
async function handleAddTodo(event) {
    event.preventDefault();
    const title = document.getElementById('todoTitle').value;
    const content = document.getElementById('todoContent').value;
    const startTime = new Date(document.getElementById('startTime').value).getTime() / 1000;
    const endTime = new Date(document.getElementById('endTime').value).getTime() / 1000;

    try {
        const response = await fetch(`${API_BASE_URL}/todo/add`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
            body: JSON.stringify({
                title,
                content,
                start_time: startTime,
                end_time: endTime
            }),
        });

        if (response.ok) {
            document.getElementById('todoTitle').value = '';
            document.getElementById('todoContent').value = '';
            document.getElementById('startTime').value = '';
            document.getElementById('endTime').value = '';
            fetchTodos();
        } else {
            const data = await response.json();
            alert('添加失败：' + data.error);
        }
    } catch (error) {
        alert('添加失败：' + error.message);
    }
}

// 获取待办事项列表
async function fetchTodos() {
    try {
        let url = `${API_BASE_URL}/todo/list?page=${currentPage}&page_size=${pageSize}`;
        if (currentStatus !== null) {
            url += `&status=${currentStatus}`;
        }

        const response = await fetch(url, {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            }
        });

        if (response.ok) {
            const data = await response.json();
            totalItems = data.data.total;
            renderTodos(data.data.items);
            updatePagination();
        } else {
            throw new Error('获取待办事项失败');
        }
    } catch (error) {
        alert(error.message);
    }
}

// 修改日期格式化函数
function formatDateTime(timestamp) {
    // 检查时间戳是否有效
    if (!timestamp) {
        return '未设置';
    }
    
    // 如果时间戳是字符串格式，需要先转换为数字
    const date = new Date(typeof timestamp === 'string' ? timestamp : timestamp * 1000);
    
    // 检查日期是否有效
    if (isNaN(date.getTime())) {
        return '无效日期';
    }
    
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    
    return `${year}-${month}-${day} ${hours}:${minutes}`;
}

// 渲染待办事项列表
function renderTodos(todos) {
    const todoList = document.getElementById('todoList');
    todoList.innerHTML = '';

    // 添加批量操作按钮
    const batchActions = document.createElement('div');
    batchActions.className = 'batch-actions';
    batchActions.innerHTML = `
        <button onclick="batchUpdateStatus(1)">批量完成</button>
        <button onclick="batchUpdateStatus(0)">批量取消完成</button>
        <button onclick="batchDelete()" class="delete-btn">批量删除</button>
    `;
    todoList.appendChild(batchActions);

    todos.forEach(todo => {
        const todoItem = document.createElement('div');
        todoItem.className = `todo-item ${todo.status === 1 ? 'status-complete' : ''}`;
        
        todoItem.innerHTML = `
            <input type="checkbox" class="todo-checkbox" 
                   onchange="toggleTodoSelection(${todo.id})" 
                   ${selectedTodos.has(todo.id) ? 'checked' : ''}>
            <div class="todo-header">
                <span class="todo-title">${todo.title}</span>
                <div class="todo-actions">
                    <button onclick="updateTodoStatus(${todo.id}, ${todo.status === 0 ? 1 : 0})">
                        ${todo.status === 0 ? '完成' : '取消完成'}
                    </button>
                    <button onclick="deleteTodo(${todo.id})" class="delete-btn">删除</button>
                </div>
            </div>
            <div class="todo-content">${todo.content}</div>
            <div class="todo-time">
                开始时间: ${formatDateTime(todo.start_time)}
                <br>
                结束时间: ${formatDateTime(todo.end_time)}
            </div>
        `;
        
        todoList.appendChild(todoItem);
    });
}

function toggleTodoSelection(id) {
    if (selectedTodos.has(id)) {
        selectedTodos.delete(id);
    } else {
        selectedTodos.add(id);
    }
}

async function batchUpdateStatus(status) {
    if (selectedTodos.size === 0) {
        alert('请先选择待办事项');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/todo/status/batch`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
            body: JSON.stringify({
                status: status,
                ids: Array.from(selectedTodos)
            }),
        });

        if (response.ok) {
            selectedTodos.clear();
            fetchTodos();
        } else {
            throw new Error('批量更新失败');
        }
    } catch (error) {
        alert(error.message);
    }
}

async function batchDelete() {
    if (selectedTodos.size === 0) {
        alert('请先选择待办事项');
        return;
    }

    if (!confirm('确定要删除选中的待办事项吗？')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/todo/batch`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
            body: JSON.stringify({
                ids: Array.from(selectedTodos)
            }),
        });

        if (response.ok) {
            selectedTodos.clear();
            fetchTodos();
        } else {
            throw new Error('批量删除失败');
        }
    } catch (error) {
        alert(error.message);
    }
}

// 更新待办事项状态
async function updateTodoStatus(id, status) {
    try {
        const response = await fetch(`${API_BASE_URL}/todo/${id}/status/${status}`, {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
        });

        if (response.ok) {
            fetchTodos();
        } else {
            throw new Error('更新状态失败');
        }
    } catch (error) {
        alert(error.message);
    }
}

// 搜索待办事项
async function handleSearch() {
    const keyword = document.getElementById('searchInput').value;
    try {
        let url = `${API_BASE_URL}/todo/search?keyword=${encodeURIComponent(keyword)}&page=${currentPage}&page_size=${pageSize}`;
        if (currentStatus !== null) {
            url += `&status=${currentStatus}`;
        }

        const response = await fetch(url, {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            }
        });

        if (response.ok) {
            const data = await response.json();
            totalItems = data.data.total;
            renderTodos(data.data.items);
            updatePagination();
        } else {
            throw new Error('搜索失败');
        }
    } catch (error) {
        alert(error.message);
    }
}

// 分页相关函数
function handlePageSizeChange() {
    pageSize = parseInt(document.getElementById('pageSize').value);
    currentPage = 1;
    fetchTodos();
}

function changePage(delta) {
    const maxPage = Math.ceil(totalItems / pageSize);
    const newPage = currentPage + delta;
    
    if (newPage >= 1 && newPage <= maxPage) {
        currentPage = newPage;
        fetchTodos();
    }
}

function updatePagination() {
    const maxPage = Math.ceil(totalItems / pageSize);
    document.getElementById('currentPage').textContent = `第 ${currentPage} 页 / 共 ${maxPage} 页`;
}

function handleLogout() {
    localStorage.removeItem('token');
    showLogin();
}

// 添加单个删除函数
async function deleteTodo(id) {
    if (!confirm('确定要删除这个待办事项吗？')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/todo/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            }
        });

        if (response.ok) {
            fetchTodos();
        } else {
            throw new Error('删除失败');
        }
    } catch (error) {
        alert(error.message);
    }
}

// 修改日期时间输入框的本地化显示
function initializeDateTimeInputs() {
    const startTimeInput = document.getElementById('startTime');
    const endTimeInput = document.getElementById('endTime');
    
    // 设置输入框的语言为 en-GB 以使用统一的格式
    startTimeInput.lang = 'en-GB';
    endTimeInput.lang = 'en-GB';
    
    // 可选：设置默认值为当前时间
    const now = new Date();
    const localDateTime = now.toISOString().slice(0, 16); // 格式：YYYY-MM-DDTHH:mm
    startTimeInput.value = localDateTime;
    endTimeInput.value = localDateTime;
}

// 初始化
document.addEventListener('DOMContentLoaded', () => {
    if (localStorage.getItem('token')) {
        showTodoApp();
        addStatusFilter();
        initializeDateTimeInputs();
    } else {
        showLogin();
    }
}); 