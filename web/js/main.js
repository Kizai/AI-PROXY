// 只保留API配置SPA逻辑
const pageRoutes = {
    'api-config': 'pages/api-config.html'
};

async function loadPage(pageName) {
    if (!checkAuth()) {
        showTokenModal();
        return;
    }
    try {
        showLoading();
        updateNavbar(pageName);
        const response = await fetch(pageRoutes[pageName]);
        if (!response.ok) throw new Error(`页面加载失败: ${response.status} - ${response.statusText}`);
        const html = await response.text();
        const mainContent = document.getElementById('main-content');
        if (mainContent) mainContent.innerHTML = html;
        initializePage(pageName);
    } catch (error) {
        const mainContent = document.getElementById('main-content');
        if (mainContent) mainContent.innerHTML = `<div class='alert alert-danger'>页面加载失败: ${error.message}</div>`;
    } finally {
        hideLoading();
    }
}

function updateNavbar(pageName) {
    document.querySelectorAll('#sidebar .components li').forEach(li => li.classList.remove('active'));
    const currentLink = document.querySelector(`#sidebar .components li a[onclick=\"loadPage('${pageName}')\"]`);
    if (currentLink) currentLink.parentElement.classList.add('active');
}

function initializePage(pageName) {
    if (pageName === 'api-config') {
        if (window.loadAPIConfigs) window.loadAPIConfigs();
    }
}

// 认证、加载动画、错误提示等工具函数可保留在utils.js
// 其余无关功能全部移除
