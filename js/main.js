const API_BASE_URL = 'http://localhost:8888';

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

async function handleAddTodo(event) {
    event.preventDefault();
    const title = document.getElementById('todoTitle').value;
    const dueDate = document.getElementById('todoDueDate').value;
    const priority = document.getElementById('todoPriority').value;

    try {
        const response = await fetch(`${API_BASE_URL}/todo/add`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
            body: JSON.stringify({ title, dueDate, priority }),
        });

        if (response.ok) {
            document.getElementById('todoTitle').value = '';
            document.getElementById('todoDueDate').value = '';
            document.getElementById('todoPriority').value = 'medium';
            fetchTodos();
        } else {
            const data = await response.json();
            alert('添加失败：' + data.message);
        }
    } catch (error) {
        alert('添加失败：' + error.message);
    }
}

async function fetchTodos() {
    try {
        const response = await fetch(`${API_BASE_URL}/todo/list`, {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
        });

        if (response.ok) {
            const data = await response.json();
            renderTodos(data.todos);
        } else {
            throw new Error('获取待办事项失败');
        }
    } catch (error) {
        alert(error.message);
    }
}

function renderTodos(todos) {
    const todoList = document.getElementById('todoList');
    todoList.innerHTML = '';

    todos.forEach(todo => {
        const todoItem = document.createElement('div');
        todoItem.className = `todo-item ${todo.status === 1 ? 'completed' : ''}`;
        todoItem.innerHTML = `
            <span>${todo.title}</span>
            <div>
                <button onclick="updateTodoStatus(${todo.id}, ${todo.status === 0 ? 1 : 0})">
                    ${todo.status === 0 ? '完成' : '取消完成'}
                </button>
            </div>
        `;
        todoList.appendChild(todoItem);
    });
}

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

function handleLogout() {
    localStorage.removeItem('token');
    showLogin();
}

// 初始化
document.addEventListener('DOMContentLoaded', () => {
    if (localStorage.getItem('token')) {
        showTodoApp();
    } else {
        showLogin();
    }
}); 