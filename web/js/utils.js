const API_BASE_URL = 'http://localhost:8080';

window.request = async function request(url, options = {}) {
    const token = localStorage.getItem('admin_token');

    const defaultOptions = {
        headers: {
            'Content-Type': 'application/json',
            ...(token ? { 'Authorization': `Bearer ${token}` } : {})
      }
    };

    try {
        const response = await fetch(`${API_BASE_URL}${url}`, { ...defaultOptions, ...options });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error('请求失败', error);
        throw error;
    }
 }

window.showSuccess = function showSuccess(message) {
    alert('成功: ' + message);
 }

window.showError = function showError(message) {
    alert('错误: ' + message);
 }

window.showLoading = function showLoading() {
    const mainContent = document.getElementById('main-content');
    if (mainContent) {
        mainContent.innerHTML = `
       <div class="text-center">
            <div class="spinner-border" role="status">
            <span class="visually-hidden">加载中...</span>
            </div>
        </div>
    `;
 }
}

window.hideLoading = function hideLoading() {
    // 隐藏加载状态
    const mainContent = document.getElementById('main-content');
    if (mainContent && mainContent.innerHTML.includes('spinner-border')) {
        // 如果当前显示的是加载状态，清除它
        // 这里不直接清空，因为页面内容应该已经加载完成
    }
}

window.checkAuth = function checkAuth() {
    const token = localStorage.getItem('admin_token');
    if (!token) {
        console.log('未找到token，显示登录模态框');
        showTokenModal();
        return false;
    }
    return true;
}

// 页面加载时检查认证 - 移除重复的监听器，统一在main.js中处理

window.showTokenModal = function showTokenModal() {
    console.log('显示登录模态框');
    
    // 确保模态框存在
    const modalEl = document.getElementById('loginModal');
    if (!modalEl) {
        console.error('登录模态框不存在');
        return;
    }
    
    // 设置模态框内容
    const titleEl = modalEl.querySelector('.modal-title');
    if (titleEl) {
        titleEl.textContent = '管理员登录';
    }
    
    const tokenInput = document.getElementById('token');
    if (tokenInput) {
        tokenInput.placeholder = '请输入访问令牌 (默认: 123456)';
        tokenInput.value = '123456'; // 设置默认值为123456
    }
    
    // 显示模态框 - 使用正确的Bootstrap 5语法
    try {
        const modal = new bootstrap.Modal(modalEl, {
            backdrop: 'static',
            keyboard: false
        });
        modal.show();
        console.log('登录模态框已显示');
    } catch (error) {
        console.error('显示模态框失败:', error);
        // 备用方案：直接显示模态框
        modalEl.style.display = 'block';
        modalEl.classList.add('show');
        const backdrop = document.createElement('div');
        backdrop.className = 'modal-backdrop fade show';
        document.body.appendChild(backdrop);
    }
}

// 设置 API Token
window.login = async function login() {
    const token = document.getElementById('token').value;
    if (!token) {
        showError("请输入访问令牌");
        return;
    }

    try {
        console.log('用户登录，token:', token);
        
        // 保存token（不验证，直接保存）
        localStorage.setItem('admin_token', token);

        // 关闭模态框
        try {
            const loginModal = bootstrap.Modal.getInstance(document.getElementById('loginModal'));
            if (loginModal) {
                loginModal.hide();
            }
            // 无论如何都尝试移除所有遮罩
            document.querySelectorAll('.modal-backdrop').forEach(el => el.remove());
            // 彻底隐藏模态框
            const modalEl = document.getElementById('loginModal');
            if (modalEl) {
                modalEl.style.display = 'none';
                modalEl.classList.remove('show');
            }
        } catch (error) {
            console.error('关闭模态框失败:', error);
        }

        console.log('登录成功，开始加载仪表板');
        
        // 加载仪表板
        loadPage('dashboard');

        showSuccess('登录成功');
    } catch (error) {
        console.error('登录失败:', error);
        showError('登录失败: ' + error.message);
    }
}

// 清除token
window.logout = function logout() {
    localStorage.removeItem('admin_token');
    showTokenModal();
    // 清空主内容区
    const mainContent = document.getElementById('main-content');
    if (mainContent) {
        mainContent.innerHTML = '';
    }
 }


