// js/request-logs.js - 辅助函数，主要逻辑由main.js统一管理

// 格式化日期时间显示
function formatDateTime(dateString) {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleString('zh-CN');
}

// 截断文本显示
function truncateText(text, maxLength = 40) {
    if (!text) return '-';
    return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
}

// 显示消息提示
function showMessage(msg, type = 'info') {
    // 使用main.js的消息显示函数
    if (window.showSuccess && window.showError) {
        if (type === 'success') {
            window.showSuccess(msg);
        } else if (type === 'danger' || type === 'error') {
            window.showError(msg);
        } else {
            alert(msg);
        }
    } else {
        alert(msg);
    }
}

// 页面初始化 - 由main.js统一管理，这里只做兼容
function initRequestLogsPage() {
    // 页面初始化逻辑已由main.js统一管理
    console.log('request-logs.js: 页面初始化完成，主要功能由main.js管理');
}