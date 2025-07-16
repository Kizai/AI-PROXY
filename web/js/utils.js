const API_BASE_URL = '';

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

        if (response.status === 401) {
            // Token无效，清除本地token并显示登录模态框
            localStorage.removeItem('admin_token');
            showError("访问令牌已失效，请重新登录");
            showTokenModal();
            throw new Error("访问令牌已失效");
        }
        
        if (response.status === 403) {
            showError("权限不足，无法访问此功能");
            throw new Error("权限不足");
        }
        
        if (!response.ok) {
            const errorText = await response.text();
            let errorMessage = `请求失败: ${response.status}`;
            try {
                const errorJson = JSON.parse(errorText);
                if (errorJson.error) {
                    errorMessage = errorJson.error;
                }
            } catch (e) {
                // 如果解析JSON失败，使用原始错误文本
                if (errorText) {
                    errorMessage = errorText;
                }
            }
            showError(errorMessage);
            throw new Error(errorMessage);
        }
        
        return await response.json();
    } catch (error) {
        console.error('请求失败', error);
        // 如果是网络错误或其他非HTTP错误，显示通用错误信息
        if (!error.message.includes('访问令牌已失效') && !error.message.includes('权限不足')) {
            showError('网络请求失败，请检查网络连接');
        }
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
        tokenInput.placeholder = '请输入访问令牌';
        tokenInput.value = ''; // 不自动填充，保持为空
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
        console.log('用户登录，验证token:', token);
        
        // 先验证token是否正确
        const testResponse = await fetch('/admin/api-config', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (testResponse.status === 401) {
            showError("访问令牌错误，请检查后重试");
            return;
        }
        
        if (!testResponse.ok) {
            showError(`验证失败: ${testResponse.status} - ${testResponse.statusText}`);
            return;
        }
        
        // 验证成功，保存token
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

        console.log('登录成功，开始加载API配置页面');
        // 加载API配置页面
        loadPage('api-config');
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


