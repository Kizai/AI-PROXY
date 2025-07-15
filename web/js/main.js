//主要页面逻辑
let currentPage = 'dashboard';

//页面路由映射
const pageRoutes = {
    'dashboard': 'pages/dashboard.html',
    'api-config': 'pages/api-config.html',
    'request-log': 'pages/request-logs.html',
    'statistics': 'pages/statistics.html'
};

//加载页面内容
async function loadPage(pageName) {
    // 检查认证状态
    if (!checkAuth()) {
        console.log('用户未登录，显示登录模态框');
        showTokenModal();
        return;
    }

    try {
        showLoading();
        updateNavbar(pageName);

        //加载页面内容
        const response = await fetch(pageRoutes[pageName]);
        if (!response.ok) {
            throw new Error(`页面加载失败: ${response.status} - ${response.statusText}`);
        }

        const html = await response.text();
        const mainContent = document.getElementById('main-content');
        if (mainContent) {
            mainContent.innerHTML = html;
        }

        //执行页面特定的初始化
        initializePage(pageName);

        currentPage = pageName;

    } catch (error) {
        console.error('页面加载失败', error);
        
        // 如果页面文件不存在，显示默认内容
        const mainContent = document.getElementById('main-content');
        if (mainContent) {
            mainContent.innerHTML = `
                <div class="row">
                    <div class="col-12">
                        <div class="card">
                            <div class="card-header">
                                <h5 class="card-title mb-0">${getPageTitle(pageName)}</h5>
                            </div>
                            <div class="card-body">
                                <div class="text-center">
                                    <i class="bi bi-tools fs-1 text-muted"></i>
                                    <p class="mt-3 text-muted">此功能正在开发中...</p>
                                    <p class="text-muted">页面文件: ${pageRoutes[pageName]}</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            `;
        }
        
        showError('页面加载失败: ' + error.message);
    } finally {
        // 确保加载状态被清除
        hideLoading();
    }
}

// 获取页面标题
function getPageTitle(pageName) {
    const titles = {
        'dashboard': '仪表板',
        'api-config': 'API配置',
        'request-log': '请求日志',
        'statistics': '统计数据'
    };
    return titles[pageName] || '页面';
}

//更新导航栏激活状态
function updateNavbar(pageName) {
    //移除所有激活状态
    document.querySelectorAll('#sidebar .components li').forEach(li => {
        li.classList.remove('active');
    });
    
    //设置当前页面的激活状态
    const currentLink = document.querySelector(`#sidebar .components li a[onclick="loadPage('${pageName}')"]`);
    if (currentLink) {
        currentLink.parentElement.classList.add('active');
    }
}

// 初始化页面特定功能
function initializePage(pageName) {
    switch (pageName) {
        case 'dashboard':
            loadDashboard();
            break;
        case 'api-config':
            loadApiConfig();
            break;
        case 'request-log':
            loadRequestLogs();
            break;
        case 'statistics':
            loadStatisticsPage();
            break;
    }
}

// 仪表板初始化
function loadDashboard() {
    console.log('加载仪表板');
    // 这里可以添加仪表板数据加载逻辑
    // 例如：加载统计数据、图表等
}

// API配置页面初始化
function loadApiConfig() {
    console.log('加载API配置页面');
    // 加载API配置数据
    window.loadAPIConfigs();
}

// 请求日志页面初始化
function loadRequestLogs() {
    console.log('加载请求日志页面');
    // 这里可以添加请求日志页面初始化逻辑
}

// 统计数据页面初始化
function loadStatisticsPage() {
    console.log('加载统计数据页面');
    // 加载统计数据页面
    window.loadStatisticsPage();
}

// ========== 仪表板相关功能 ==========
window.loadDashboard = async function loadDashboard() {
    window.loadStatistics();
    window.loadRequestTrend();
    window.loadAPIStatus();
    window.loadRecentRequests();
    if (window.dashboardInterval) clearInterval(window.dashboardInterval);
    window.dashboardInterval = setInterval(() => {
        window.loadStatistics();
        window.loadRecentRequests();
    }, 30000);
};
window.loadStatistics = async function loadStatistics() {
    try {
        const response = await request('/admin/stats');
        if (response.success || response.code === 200) {
            const stats = response.data;
            const el = id => document.getElementById(id);
            if (el('totalRequests')) el('totalRequests').textContent = stats.total_requests || 0;
            if (el('successRequests')) el('successRequests').textContent = stats.success_requests || 0;
            if (el('errorRequests')) el('errorRequests').textContent = stats.error_requests || 0;
            if (el('avgResponseTime')) el('avgResponseTime').textContent = ((stats.avg_response_time || 0) / 1000).toFixed(1) + 's';
            if (el('totalRequestsStats')) el('totalRequestsStats').textContent = stats.total_requests || 0;
            if (el('successRateStats')) el('successRateStats').textContent = Number(stats.success_rate || 0).toFixed(2) + '%';
            if (el('avgResponseTimeStats')) el('avgResponseTimeStats').textContent = ((stats.avg_response_time || 0) / 1000).toFixed(1) + 's';
            if (el('activeAPIsStats')) el('activeAPIsStats').textContent = stats.active_apis || 0;
        }
    } catch (error) {
        console.error('加载统计数据失败:', error);
    }
};
window.loadRequestTrend = async function loadRequestTrend() {
    try {
        let url = '/admin/stats/realtime';
        const trendPeriod = document.getElementById('trendPeriod');
        if (trendPeriod) url += `?period=${trendPeriod.value}`;
        const response = await request(url);
        if (response.success || response.code === 200) {
            window.renderRequestTrendChart(response.data);
        }
    } catch (error) {
        console.error('加载请求趋势失败:', error);
    }
};
window.renderRequestTrendChart = function renderRequestTrendChart(data) {
    const chartEl = document.getElementById('requestTrendChart');
    if (!chartEl) return;
    const ctx = chartEl.getContext('2d');
    if (window.requestTrendChart) window.requestTrendChart.destroy();
    window.requestTrendChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: data.labels || [],
            datasets: [{
                label: '请求数',
                data: data.values || [],
                borderColor: '#667eea',
                backgroundColor: 'rgba(102, 126, 234, 0.1)',
                tension: 0.4,
                fill: true
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: { legend: { display: false } },
            scales: { y: { beginAtZero: true } }
        }
    });
};
window.loadAPIStatus = async function loadAPIStatus() {
    try {
        const response = await request('/admin/api-config');
        if (response.success || response.code === 200) {
            window.renderAPIStatus(response.data);
        }
    } catch (error) {
        console.error('加载API状态失败:', error);
    }
};
window.renderAPIStatus = function renderAPIStatus(apis) {
    const container = document.getElementById('apiStatusList');
    if (!container) return;
    container.innerHTML = '';
    if (!apis || apis.length === 0) {
        container.innerHTML = '<p class="text-muted">暂无API配置</p>';
        return;
    }
    apis.forEach((api, index) => {
        let statusClass, statusIcon;
        if (api.last_test_status === 'success') {
            statusClass = 'text-success';
            statusIcon = 'bi-check-circle';
        } else {
            statusClass = 'text-danger';
            statusIcon = 'bi-x-circle';
        }
        const apiElement = document.createElement('div');
        apiElement.className = 'd-flex justify-content-between align-items-center mb-2 pb-2';
        apiElement.style.borderBottom = '1px solid #e9ecef';
        apiElement.innerHTML = `
            <div>
                <span style="font-size:0.9rem;color:#2563eb;font-weight:600;">${api.name}</span>
            </div>
            <div class="${statusClass}">
                <i class="bi ${statusIcon}"></i>
            </div>
        `;
        container.appendChild(apiElement);
        
        // 最后一个元素不需要底部边框
        if (index === apis.length - 1) {
            apiElement.style.borderBottom = 'none';
        }
    });
};
window.loadRecentRequests = async function loadRecentRequests() {
    try {
        const response = await request('/admin/logs?page=1&size=10');
        if (response.success || response.code === 200) {
            window.renderRecentRequests(response.data.logs || []);
        }
    } catch (error) {
        console.error('加载最近请求失败:', error);
    }
};
window.renderRecentRequests = function renderRecentRequests(logs) {
    const tbody = document.getElementById('recentRequestsTable');
    if (!tbody) return;
    tbody.innerHTML = '';
    if (logs.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6" class="text-center">暂无数据</td></tr>';
        return;
    }
    logs.forEach(log => {
        const statusClass = log.response_status >= 400 ? 'text-danger' : 'text-success';
        const row = document.createElement('tr');
        row.innerHTML = `
            <td style="font-weight:500;">${window.formatDateTime(log.created_at)}</td>
            <td style="font-weight:500;">${log.api_name || '-'} </td>
            <td style="font-weight:500;"><span class="badge bg-secondary">${log.request_method}</span></td>
            <td style="font-weight:500;">${window.truncateText(log.request_path, 30)}</td>
            <td style="font-weight:500;"><span class="${statusClass}">${log.response_status}</span></td>
            <td style="font-weight:500;">${((log.response_time || 0) / 1000).toFixed(1)}s</td>
        `;
        tbody.appendChild(row);
    });
};
window.refreshRecentRequests = function refreshRecentRequests() {
    window.loadRecentRequests();
};
window.viewLogDetail = function viewLogDetail(logId) {
    // 这里可以实现查看日志详情的功能
    alert('查看日志详情功能待实现: ' + logId);
};
window.formatDateTime = function formatDateTime(dateString) {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleString('zh-CN');
};
window.truncateText = function truncateText(text, maxLength) {
    if (!text) return '-';
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
};

// ========== API配置相关功能 ==========
let currentEditingAPI = null;
window.loadApiConfig = function loadApiConfig() {
    window.loadAPIConfigs();
};
window.loadAPIConfigs = async function loadAPIConfigs() {
    try {
        const response = await request('/admin/api-config');
        console.log('API配置响应:', response); // 调试信息
        if (response.success || response.code === 200) {
            window.renderAPIConfigs(response.data);
        } else {
            console.error('API配置加载失败:', response);
            showError('加载API配置失败: ' + (response.message || '未知错误'));
        }
    } catch (error) {
        console.error('加载API配置失败:', error);
        showError('加载API配置失败: ' + (error.message || '网络错误'));
    }
};
window.renderAPIConfigs = function renderAPIConfigs(apis) {
    const tbody = document.getElementById('apiConfigTable');
    if (!tbody) return;
    tbody.innerHTML = '';
    if (!apis || apis.length === 0) {
        tbody.innerHTML = '<tr><td colspan="8" class="text-center">暂无API配置</td></tr>';
        return;
    }
    apis.forEach(api => {
        const statusBadge = api.active ? '<span class="badge bg-success">启用</span>' : '<span class="badge bg-secondary">禁用</span>';
        const row = document.createElement('tr');
        row.innerHTML = `
            <td><strong class="text-primary">${api.name}</strong></td>
            <td><span class="text-dark">${window.truncateText(api.base_url, 35)}</span></td>
            <td><span class="text-muted">${api.auth_type || 'none'}</span></td>
            <td><span class="text-dark">${api.timeout}s</span></td>
            <td>${statusBadge}</td>
            <td><small class="text-muted">${api.description || '-'}</small></td>
            <td>
                <div class="btn-group btn-group-sm" role="group">
                    <button class="btn btn-outline-success" data-action="showApiTestModal" data-api-name="${api.name}" title="测试API">
                        <i class="bi bi-play-circle"></i> 测试
                    </button>
                    <button class="btn btn-outline-primary" data-action="editAPIConfig" data-api-name="${api.name}" title="编辑">
                        <i class="bi bi-pencil"></i>
                    </button>
                    <button class="btn btn-outline-danger" data-action="deleteAPIConfig" data-api-name="${api.name}" title="删除">
                        <i class="bi bi-trash"></i>
                    </button>
                </div>
            </td>
        `;
        tbody.appendChild(row);
    });
};
window.showAddAPIModal = function showAddAPIModal() {
    currentEditingAPI = null;
    document.getElementById('apiModalTitle').textContent = '添加API配置';
    document.getElementById('apiForm').reset();
    document.getElementById('apiName').readOnly = false;
    document.getElementById('apiName').value = '';
    document.getElementById('apiBaseURL').value = '';
    document.getElementById('authType').value = 'bearer';
    document.getElementById('authValue').value = '';
    document.getElementById('timeout').value = 30;
    document.getElementById('rateLimit').value = 0;
    document.getElementById('headers').value = '';
    document.getElementById('description').value = '';
    document.getElementById('isActive').checked = true;
    // 只创建一次Modal实例，后续复用
    const modalEl = document.getElementById('apiModal');
    let modal = bootstrap.Modal.getInstance(modalEl);
    if (!modal) {
        modal = new bootstrap.Modal(modalEl);
    }
    modal.show();
};
window.editAPIConfig = async function editAPIConfig(apiName) {
    console.log('editAPIConfig被调用，API名称:', apiName); // 调试信息
    try {
        console.log('开始请求API配置...'); // 调试信息
        const response = await request(`/admin/api-config/${apiName}`);
        console.log('API配置响应:', response); // 调试信息
        if (response.success || response.code === 200) {
            currentEditingAPI = response.data;
            console.log('准备显示编辑模态框...'); // 调试信息
            window.showEditAPIModal(currentEditingAPI);
        } else {
            console.error('API配置请求失败:', response); // 调试信息
            showError('获取API配置失败: ' + (response.message || '未知错误'));
        }
    } catch (error) {
        console.error('获取API配置失败:', error);
        showError('获取API配置失败: ' + error.message);
    }
};
window.showEditAPIModal = function showEditAPIModal(api) {
    console.log('showEditAPIModal被调用，API数据:', api); // 调试信息
    try {
        document.getElementById('apiModalTitle').textContent = '编辑API配置';
        document.getElementById('apiName').value = api.name;
        document.getElementById('apiName').readOnly = true;
        document.getElementById('apiBaseURL').value = api.base_url;
        document.getElementById('authType').value = api.auth_type || 'none';
        document.getElementById('authValue').value = api.auth_value || '';
        document.getElementById('timeout').value = api.timeout || 30;
        document.getElementById('rateLimit').value = api.rate_limit || 0;
        document.getElementById('headers').value = api.headers || '';
        document.getElementById('description').value = api.description || '';
        document.getElementById('isActive').checked = api.active !== false;
        
        console.log('准备显示模态框...'); // 调试信息
        const modalEl = document.getElementById('apiModal');
        if (!modalEl) {
            console.error('模态框元素不存在!'); // 调试信息
            return;
        }
        const modal = new bootstrap.Modal(modalEl);
        modal.show();
        console.log('模态框已显示'); // 调试信息
    } catch (error) {
        console.error('显示编辑模态框失败:', error); // 调试信息
        showError('显示编辑模态框失败: ' + error.message);
    }
};
window.saveAPIConfig = async function saveAPIConfig() {
    const form = document.getElementById('apiForm');
    if (!form.checkValidity()) {
        form.reportValidity();
        return;
    }
    // headers校验，只做合法JSON校验，始终作为字符串传递
    let headersValue = document.getElementById('headers').value;
    let headersStr = '';
    if (headersValue.trim()) {
        try {
            JSON.parse(headersValue); // 只校验格式
            headersStr = headersValue.trim();
        } catch (e) {
            showError('自定义请求头必须是合法的JSON格式！如：{"Content-Type": "application/json"}');
            document.getElementById('headers').focus();
            return;
        }
    }
    const apiData = {
        name: document.getElementById('apiName').value,
        base_url: document.getElementById('apiBaseURL').value,
        auth_type: document.getElementById('authType').value,
        auth_value: document.getElementById('authValue').value,
        timeout: parseInt(document.getElementById('timeout').value),
        rate_limit: parseInt(document.getElementById('rateLimit').value),
        headers: headersStr,
        description: document.getElementById('description').value,
        active: document.getElementById('isActive').checked
    };
    try {
        let response;
        if (currentEditingAPI) {
            response = await request(`/admin/api-config/${apiData.name}`, {
                method: 'PUT',
                body: JSON.stringify(apiData)
            });
        } else {
            response = await request('/admin/api-config', {
                method: 'POST',
                body: JSON.stringify(apiData)
            });
        }
        if (response.success || response.code === 200) {
            showSuccess('保存成功');
            const modal = bootstrap.Modal.getInstance(document.getElementById('apiModal'));
            if (modal) modal.hide();
            window.loadAPIConfigs();
        } else {
            showError(response.message || '保存失败');
        }
    } catch (error) {
        // 优先显示后端返回的message
        if (error && error.message) {
            showError('保存失败: ' + error.message);
        } else {
            showError('保存失败');
        }
        // 即使保存"失败"，也重新加载列表，以防数据实际已保存成功
        window.loadAPIConfigs();
    }
};
window.deleteAPIConfig = async function deleteAPIConfig(apiName) {
    if (!confirm('确定要删除该API配置吗？')) return;
    try {
        const response = await request(`/admin/api-config/${apiName}`, { method: 'DELETE' });
        if (response.success || response.code === 200) {
            showSuccess('删除成功');
            window.loadAPIConfigs();
        } else {
            showError(response.message || '删除失败');
        }
    } catch (error) {
        showError('删除失败: ' + error.message);
    }
};

// ========== 请求日志相关功能 ==========
window.loadRequestLogs = function loadRequestLogs() {
    window.loadRequestLogsData();
};
window.loadRequestLogsData = async function loadRequestLogsData(page = 1) {
    page = Number(page) || 1;
    try {
        const params = new URLSearchParams({ page: page, size: 20 });
        // 筛选条件 - 参数名与后端结构体一致
        const apiName = document.getElementById('filterApiName')?.value;
        const method = document.getElementById('filterMethod')?.value;
        const status = document.getElementById('filterStatus')?.value;
        const hasError = document.getElementById('filterHasError')?.value;
        const startDate = document.getElementById('filterStartDate')?.value;
        const endDate = document.getElementById('filterEndDate')?.value;
        
        if (apiName) params.append('api_name', apiName);
        if (method) params.append('request_method', method);
        if (status) params.append('status_code', status); // 保持参数名
        if (hasError !== '') params.append('has_error', hasError);
        if (startDate) params.append('start_time', new Date(startDate).toISOString());
        if (endDate) params.append('end_time', new Date(endDate).toISOString());
        
        const url = `/admin/logs?${params.toString()}`;
        console.log('[调试] 请求日志URL:', url);
        const response = await request(url);
        console.log('[调试] 后端原始返回:', response);
        if (response.success || response.code === 200) {
            window.requestLogsCurrentPage = page; // 提前到这里
            console.log('[调试] 渲染日志数据:', response.data);
            window.renderRequestLogs(response.data);
        } else {
            console.warn('[调试] 日志接口返回非成功:', response);
        }
    } catch (error) {
        console.error('加载请求日志失败:', error);
        showError('加载请求日志失败: ' + error.message);
    }
};
window.renderRequestLogs = function renderRequestLogs(data) {
    const tbody = document.getElementById('requestLogsTable');
    if (!tbody) return;
    tbody.innerHTML = '';
    if (!data.logs || data.logs.length === 0) {
        tbody.innerHTML = '<tr><td colspan="8" class="text-center">暂无数据</td></tr>';
        return;
    }
    data.logs.forEach(log => {
        const statusClass = log.response_status >= 400 ? 'text-danger' : 'text-success';
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${window.formatDateTime(log.created_at)}</td>
            <td>${log.api_name || '-'}</td>
            <td><span class="badge bg-secondary">${log.request_method}</span></td>
            <td>${window.truncateText(log.request_path, 40)}</td>
            <td><span class="${statusClass}">${log.response_status}</span></td>
            <td>${log.response_time}ms</td>
            <td>${log.user_ip || '-'}</td>
            <td>
                <button class="btn btn-sm btn-outline-secondary" data-action="viewLogDetail" data-log-id="${log.id}">
                    <i class="bi bi-eye"></i>
                </button>
            </td>
        `;
        tbody.appendChild(row);
    });
    window.renderPagination(data.total, window.requestLogsCurrentPage, data.size);
};
window.renderPagination = function renderPagination(total, currentPage, pageSize) {
    currentPage = Number(currentPage) || 1;
    pageSize = Number(pageSize) || 10;
    const totalPages = Math.ceil(total / pageSize);
    const pagination = document.getElementById('pagination');
    if (!pagination) return;
    pagination.innerHTML = '';
    if (totalPages <= 1) return;
    // 上一页
    const prevPage = currentPage > 1 ? currentPage - 1 : 1;
    const prevLi = document.createElement('li');
    prevLi.className = `page-item ${currentPage <= 1 ? 'disabled' : ''}`;
    prevLi.innerHTML = `<a class="page-link" href="#" data-action="loadRequestLogsData" data-page="${prevPage}">上一页</a>`;
    pagination.appendChild(prevLi);
    // 页码
    for (let i = 1; i <= totalPages; i++) {
        const li = document.createElement('li');
        li.className = `page-item ${i === currentPage ? 'active' : ''}`;
        li.innerHTML = `<a class="page-link" href="#" data-action="loadRequestLogsData" data-page="${i}">${i}</a>`;
        pagination.appendChild(li);
    }
    // 下一页
    const nextPage = currentPage < totalPages ? currentPage + 1 : totalPages;
    const nextLi = document.createElement('li');
    nextLi.className = `page-item ${currentPage >= totalPages ? 'disabled' : ''}`;
    nextLi.innerHTML = `<a class="page-link" href="#" data-action="loadRequestLogsData" data-page="${nextPage}">下一页</a>`;
    pagination.appendChild(nextLi);
};
window.applyFilter = function applyFilter() {
    window.loadRequestLogsData(1);
};
window.resetFilter = function resetFilter() {
    document.getElementById('filterForm').reset();
    window.loadRequestLogsData(1);
};
window.clearLogs = async function clearLogs() {
    if (!confirm('确定要清空所有日志吗？')) return;
    try {
        const response = await request('/admin/logs/clear', { method: 'POST' }); // 恢复接口路径
        if (response.success || response.code === 200) {
            showSuccess('日志已清空');
            window.loadRequestLogsData(1);
        } else {
            showError(response.message || '清空失败');
        }
    } catch (error) {
        showError('清空失败: ' + error.message);
    }
};
window.exportLogs = async function exportLogs() {
    try {
        const params = new URLSearchParams();
        // 添加筛选条件 - 参数名与后端结构体一致
        const apiName = document.getElementById('filterApiName')?.value;
        const method = document.getElementById('filterMethod')?.value;
        const status = document.getElementById('filterStatus')?.value;
        const hasError = document.getElementById('filterHasError')?.value;
        const startDate = document.getElementById('filterStartDate')?.value;
        const endDate = document.getElementById('filterEndDate')?.value;
        
        if (apiName) params.append('api_name', apiName);
        if (method) params.append('request_method', method);
        if (status) params.append('status_code', status);
        if (hasError !== '') params.append('has_error', hasError);
        if (startDate) params.append('start_time', new Date(startDate).toISOString());
        if (endDate) params.append('end_time', new Date(endDate).toISOString());
        
        const response = await fetch(`/admin/logs/export?${params.toString()}`, { // 恢复接口路径
            headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` } // 修正token键名
        });
        if (!response.ok) throw new Error('导出失败');
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'request-logs.csv';
        document.body.appendChild(a);
        a.click();
        a.remove();
        window.URL.revokeObjectURL(url);
        showSuccess('导出成功');
    } catch (error) {
        showError('导出失败: ' + error.message);
    }
};

// ========== 统计数据相关功能 ==========
window.loadStatisticsPage = function loadStatisticsPage() {
    window.loadStatisticsSummary();
    window.loadStatisticsSummaryCards(); // 加载统计卡片
    // 暂时禁用不存在的图表接口
    // window.loadRequestTrend();
    // window.loadAPIUsage();
    // window.loadStatusCodeDistribution();
    // window.loadResponseTimeDistribution();
    window.loadAPIStatsTable();
    if (window.statisticsInterval) clearInterval(window.statisticsInterval);
    window.statisticsInterval = setInterval(() => {
        window.loadStatisticsSummary();
        window.loadStatisticsSummaryCards(); // 定期更新统计卡片
        window.loadAPIStatsTable();
    }, 60000);
};
window.loadStatisticsSummary = async function loadStatisticsSummary() {
    try {
        const response = await request('/admin/stats');
        if (response.success || response.code === 200) {
            const stats = response.data;
            
            // 更新统计数据页面的元素
            const el = id => document.getElementById(id);
            if (el('statTotalRequests')) {
                el('statTotalRequests').textContent = stats.total_requests || 0;
            }
            if (el('statSuccessRate')) {
                el('statSuccessRate').textContent = Number(stats.success_rate || 0).toFixed(2) + '%';
            }
            if (el('statAvgResponseTime')) {
                el('statAvgResponseTime').textContent = ((stats.avg_response_time || 0) / 1000).toFixed(1) + 's';
            }
            if (el('statApiCount')) {
                el('statApiCount').textContent = stats.active_apis || 0;
            }
            
            // 同时更新仪表板的元素（如果存在）
            if (el('totalRequestsStats')) el('totalRequestsStats').textContent = stats.total_requests || 0;
            if (el('successRateStats')) el('successRateStats').textContent = Number(stats.success_rate || 0).toFixed(2) + '%';
            if (el('avgResponseTimeStats')) el('avgResponseTimeStats').textContent = ((stats.avg_response_time || 0) / 1000).toFixed(1) + 's';
            if (el('activeAPIsStats')) el('activeAPIsStats').textContent = stats.active_apis || 0;
        }
    } catch (error) {
        console.error('加载统计摘要失败:', error);
    }
};
window.loadAPIUsage = async function loadAPIUsage() {
    try {
        const response = await request('/admin/stats/api-usage');
        if (response.success || response.code === 200) {
            window.renderAPIUsageChart(response.data);
        }
    } catch (error) {
        console.error('加载API使用分布失败:', error);
    }
};
window.renderAPIUsageChart = function renderAPIUsageChart(data) {
    const ctx = document.getElementById('apiUsageChart').getContext('2d');
    if (window.apiUsageChart) window.apiUsageChart.destroy();
    window.apiUsageChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: data.labels || [],
            datasets: [{
                data: data.values || [],
                backgroundColor: ['#667eea', '#48bb78', '#f6ad55', '#63b3ed', '#f56565', '#ed64a6', '#ecc94b']
            }]
        },
        options: {
            responsive: true,
            plugins: { legend: { position: 'bottom' } }
        }
    });
};
window.loadStatusCodeDistribution = async function loadStatusCodeDistribution() {
    try {
        const response = await request('/admin/stats/status-code');
        if (response.success || response.code === 200) {
            window.renderStatusCodeChart(response.data);
        }
    } catch (error) {
        console.error('加载状态码分布失败:', error);
    }
};
window.renderStatusCodeChart = function renderStatusCodeChart(data) {
    const ctx = document.getElementById('statusCodeChart').getContext('2d');
    if (window.statusCodeChart) window.statusCodeChart.destroy();
    window.statusCodeChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: data.labels || [],
            datasets: [{
                label: '请求数',
                data: data.values || [],
                backgroundColor: '#667eea'
            }]
        },
        options: {
            responsive: true,
            plugins: { legend: { display: false } },
            scales: { y: { beginAtZero: true } }
        }
    });
};
window.loadResponseTimeDistribution = async function loadResponseTimeDistribution() {
    try {
        const response = await request('/admin/stats/response-time');
        if (response.success || response.code === 200) {
            window.renderResponseTimeChart(response.data);
        }
    } catch (error) {
        console.error('加载响应时间分布失败:', error);
    }
};
window.renderResponseTimeChart = function renderResponseTimeChart(data) {
    const ctx = document.getElementById('responseTimeChart').getContext('2d');
    if (window.responseTimeChart) window.responseTimeChart.destroy();
    window.responseTimeChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: data.labels || [],
            datasets: [{
                label: '请求数',
                data: data.values || [],
                backgroundColor: '#48bb78'
            }]
        },
        options: {
            responsive: true,
            plugins: { legend: { display: false } },
            scales: { y: { beginAtZero: true } }
        }
    });
};
window.loadAPIStatsTable = async function loadAPIStatsTable() {
    try {
        const response = await request('/admin/stats/api-table');
        if (response.success || response.code === 200) {
            window.renderAPIStatsTable(response.data);
        }
    } catch (error) {
        console.error('加载API统计表失败:', error);
    }
};
window.renderAPIStatsTable = function renderAPIStatsTable(data) {
    const tbody = document.getElementById('apiStatsTable');
    if (!tbody) return;
    tbody.innerHTML = '';
    if (!data || data.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" class="text-center">暂无数据</td></tr>';
        return;
    }
    data.forEach(api => {
        const statusBadge = api.active ? '<span class="badge bg-success">启用</span>' : '<span class="badge bg-secondary">禁用</span>';
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${api.name}</td>
            <td>${api.total_requests}</td>
            <td>${api.success_requests}</td>
            <td>${api.error_requests}</td>
            <td>${api.success_rate}%</td>
            <td>${((api.avg_response_time || 0) / 1000).toFixed(1)}s</td>
            <td>${statusBadge}</td>
        `;
        tbody.appendChild(row);
    });
};
window.refreshStatistics = function refreshStatistics() {
    window.loadStatisticsPage();
};

// ========== API配置测试功能 ==========

window.showApiTestModal = async function(apiName) {
    console.log('showApiTestModal被调用，API名称:', apiName);
    
    // 获取API配置详情
    let apiConfig = null;
    try {
        const response = await request(`/admin/api-config/${apiName}`);
        console.log('API配置响应:', response);
        if (response.success || response.code === 200) {
            apiConfig = response.data;
        } else {
            showError('获取API配置失败: ' + (response.message || '未知错误'));
            return;
        }
    } catch (err) {
        console.error('获取API配置错误:', err);
        showError('获取API配置失败: ' + (err.message || err));
        return;
    }
    
    // 隐藏测试结果
    document.getElementById('apiTestResult').style.display = 'none';
    
    // 显示模态框
    const modalEl = document.getElementById('apiTestModal');
    let modal = bootstrap.Modal.getInstance(modalEl);
    if (!modal) modal = new bootstrap.Modal(modalEl);
    modal.show();
    
    // 绑定发送按钮
    document.getElementById('sendApiTestBtn').onclick = async function() {
        await window.sendApiTest(apiConfig);
    };
};

window.sendApiTest = async function(apiConfig) {
    // 显示测试结果区域
    document.getElementById('apiTestResult').style.display = 'block';
    document.getElementById('apiTestResultAlert').className = 'alert alert-info';
    document.getElementById('apiTestResultAlert').textContent = '测试中...';
    
    try {
        const res = await request('/admin/api-config/test', {
            method: 'POST',
            body: JSON.stringify({
                name: apiConfig.name
            })
        });
        console.log('API测试响应:', res);
        const testData = res.data || res;
        console.log('测试数据:', testData);
        
        // 更新状态码和响应时间
        document.getElementById('testStatusCode').textContent = testData.status || '未知';
        document.getElementById('testResponseTime').textContent = ((testData.response_time || 0) / 1000).toFixed(1) + 's';
        
        // 检查是否成功
        if (testData.success === true) {
            document.getElementById('apiTestResultAlert').className = 'alert alert-success';
            document.getElementById('apiTestResultAlert').textContent = testData.message || 'API配置测试成功';
        } else {
            document.getElementById('apiTestResultAlert').className = 'alert alert-danger';
            document.getElementById('apiTestResultAlert').textContent = testData.message || testData.error || 'API配置测试失败';
        }
        // 测试后自动刷新API状态
        window.loadAPIStatus && window.loadAPIStatus();
    } catch (err) {
        document.getElementById('apiTestResultAlert').className = 'alert alert-danger';
        document.getElementById('apiTestResultAlert').textContent = '测试请求异常: ' + (err.message || err);
        document.getElementById('testStatusCode').textContent = '错误';
        document.getElementById('testResponseTime').textContent = '0ms';
    }
};

// ========== 统计数据页面卡片渲染 ==========
window.loadStatisticsSummaryCards = async function loadStatisticsSummaryCards() {
    try {
        // 获取统计数据
        const statsRes = await request('/admin/stats');
        if (statsRes.success || statsRes.code === 200) {
            const data = statsRes.data || {};
            
            const statTotalRequestsEl = document.getElementById('statTotalRequests');
            const statSuccessRateEl = document.getElementById('statSuccessRate');
            const statAvgResponseTimeEl = document.getElementById('statAvgResponseTime');
            const statApiCountEl = document.getElementById('statApiCount');
            
            if (statTotalRequestsEl) {
                statTotalRequestsEl.textContent = data.total_requests || 0;
            }
            if (statSuccessRateEl) {
                statSuccessRateEl.textContent = Number(data.success_rate || 0).toFixed(2) + '%';
            }
            if (statAvgResponseTimeEl) {
                statAvgResponseTimeEl.textContent = ((data.avg_response_time || 0) / 1000).toFixed(1) + 's';
            }
            
            // 获取API配置总数
            const apiRes = await request('/admin/api-config');
            if (apiRes.success || apiRes.code === 200) {
                const apiCount = (apiRes.data && apiRes.data.length) || 0;
                if (statApiCountEl) {
                    statApiCountEl.textContent = apiCount;
                }
            }
        }
    } catch (err) {
        console.error('加载统计摘要卡片失败:', err);
        // 失败时显示0
        const statTotalRequestsEl = document.getElementById('statTotalRequests');
        const statSuccessRateEl = document.getElementById('statSuccessRate');
        const statAvgResponseTimeEl = document.getElementById('statAvgResponseTime');
        const statApiCountEl = document.getElementById('statApiCount');
        
        if (statTotalRequestsEl) statTotalRequestsEl.textContent = 0;
        if (statSuccessRateEl) statSuccessRateEl.textContent = '0%';
        if (statAvgResponseTimeEl) statAvgResponseTimeEl.textContent = '0s';
        if (statApiCountEl) statApiCountEl.textContent = 0;
    }
};

// 页面加载时自动渲染统计卡片
if (window.location.pathname.includes('statistics')) {
    window.addEventListener('DOMContentLoaded', function() {
        window.loadStatisticsSummaryCards();
    });
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
    console.log('页面加载完成');
    
    // 检查认证状态
    if (!checkAuth()) {
        return;
    }

    // 如果已经设置了token，加载仪表板
    loadPage('dashboard');
});

// 全局事件委托，保证所有data-action按钮都能响应
document.addEventListener('click', function(e) {
    let target = e.target;
    
    // 处理图标点击
    if (target.tagName === 'I' && target.parentElement) {
        target = target.parentElement;
    }
    
    const action = target.getAttribute('data-action');
    console.log('点击事件:', action, target); // 调试信息
    
    if (action && typeof window[action] === 'function') {
        e.preventDefault();
        console.log('执行函数:', action); // 调试信息
        
        if (target.hasAttribute('data-api-name')) {
            const apiName = target.getAttribute('data-api-name');
            console.log('API名称:', apiName); // 调试信息
            window[action](apiName);
        } else if (target.hasAttribute('data-log-id')) {
            window[action](target.getAttribute('data-log-id'));
        } else if (target.hasAttribute('data-page')) {
            window[action](target.getAttribute('data-page'));
        } else {
            window[action](e);
        }
    }
});
