// 只保留API配置相关逻辑，字段与api-config.html保持一致

// 加载API配置列表
window.loadAPIConfigs = function loadAPIConfigs() {
    fetch('/admin/api-config', {
        headers: { 'Authorization': 'Bearer ' + localStorage.getItem('admin_token') }
    })
        .then(res => res.json())
        .then(data => {
            renderApiConfigTable(data.data || []);
        })
        .catch(err => alert('加载失败: ' + err));
};

// 测试结果状态缓存
window.apiTestStatus = {};

// 渲染API配置表格
function renderApiConfigTable(list) {
    window.apiConfigList = list; // 缓存到全局变量
    const tbody = document.getElementById('apiConfigTableBody');
    tbody.innerHTML = '';
    list.forEach((item, idx) => {
        // 测试结果渲染
        let testStatus = window.apiTestStatus[item.name] || 'none';
        let testHtml = '';
        if (testStatus === 'pending') {
            testHtml = '<span class="badge bg-primary"><span class="spinner-border spinner-border-sm me-1"></span>测试中</span>';
        } else if (testStatus === 'success') {
            testHtml = '<span class="badge bg-success"><i class="bi bi-check-circle me-1"></i>可访问</span>';
        } else if (testStatus === 'fail') {
            testHtml = '<span class="badge bg-danger"><i class="bi bi-x-circle me-1"></i>不可访问</span>';
        } else {
            testHtml = '<span class="badge bg-secondary">未测试</span>';
        }
        const tr = document.createElement('tr');
        tr.innerHTML = `
          <td>${item.name}</td>
          <td>${item.base_url}</td>
          <td>${item.active ? '<span class="status-enabled">启用</span>' : '禁用'}</td>
          <td>${item.description || ''}</td>
          <td>${testHtml}</td>
          <td>
            <button class="btn btn-success btn-sm me-1" onclick="testAPIConfig('${item.name}')">测试</button>
            <button class="btn btn-info btn-sm me-1" onclick="editAPIConfig(${idx})">编辑</button>
            <button class="btn btn-danger btn-sm" onclick="deleteAPIConfig('${item.name}')">删除</button>
          </td>
        `;
        tbody.appendChild(tr);
    });
}

// 编辑API配置
window.editAPIConfig = function(index) {
    if (window.apiConfigList && window.apiConfigList[index]) {
        showAddAPIModal(window.apiConfigList[index]);
    }
};

// 显示添加/编辑API弹窗
window.showAddAPIModal = function showAddAPIModal(item) {
    item = item || null;
    document.getElementById('apiName').value = item ? item.name : '';
    document.getElementById('baseUrl').value = item ? item.base_url : '';
    document.getElementById('description').value = item ? item.description || '' : '';
    document.getElementById('enabled').checked = item ? !!item.active : true;
    document.getElementById('apiConfigModalTitle').textContent = item ? '编辑API配置' : '添加API配置';
    document.getElementById('apiConfigModal').setAttribute('data-edit-id', item ? item.name : '');
    new bootstrap.Modal(document.getElementById('apiConfigModal')).show();
};

// 保存API配置
window.saveAPIConfig = function saveAPIConfig() {
    const id = document.getElementById('apiConfigModal').getAttribute('data-edit-id');
    const isEdit = !!id;
    const name = document.getElementById('apiName').value.trim();
    const base_url = document.getElementById('baseUrl').value.trim();
    const description = document.getElementById('description').value.trim();
    const active = document.getElementById('enabled').checked;
    if (!name || !base_url) {
        alert('API名称和基址URL为必填项');
        return;
    }
    const url = isEdit ? `/admin/api-config/${encodeURIComponent(id)}` : '/admin/api-config';
    const method = isEdit ? 'PUT' : 'POST';
    const body = { name, base_url, description, active };
    fetch(url, {
        method,
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem('admin_token')
        },
        body: JSON.stringify(body)
    })
        .then(res => res.json())
        .then(data => {
            if (data.code === 200) {
                alert('保存成功');
                window.loadAPIConfigs();
                bootstrap.Modal.getInstance(document.getElementById('apiConfigModal')).hide();
            } else {
                alert('保存失败: ' + (data.message || '未知错误'));
            }
        })
        .catch(err => alert('请求失败: ' + err));
};

// 删除API配置
window.deleteAPIConfig = function deleteAPIConfig(name) {
    if (!confirm('确定要删除吗？')) return;
    fetch(`/admin/api-config/${encodeURIComponent(name)}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem('admin_token')
        }
    })
        .then(res => res.json())
        .then(data => {
            if (data.code === 200) {
                alert('删除成功');
                window.loadAPIConfigs();
            } else {
                alert('删除失败: ' + (data.message || '未知错误'));
            }
        })
        .catch(err => alert('请求失败: ' + err));
};

// 测试API可用性
window.testAPIConfig = function(name) {
    if (!name) return;
    window.apiTestStatus[name] = 'pending';
    renderApiConfigTable(window.apiConfigList);
    fetch('/admin/api-config/test', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem('admin_token')
        },
        body: JSON.stringify({ name })
    })
        .then(res => res.json())
        .then(data => {
            if (data && data.data) {
                if (data.data.success) {
                    window.apiTestStatus[name] = 'success';
                } else {
                    window.apiTestStatus[name] = 'fail';
                }
                renderApiConfigTable(window.apiConfigList);
                alert(data.data.message || (data.data.success ? 'API地址可访问' : 'API地址无法访问'));
            } else {
                window.apiTestStatus[name] = 'fail';
                renderApiConfigTable(window.apiConfigList);
                alert('测试结果未知');
            }
        })
        .catch(err => {
            window.apiTestStatus[name] = 'fail';
            renderApiConfigTable(window.apiConfigList);
            alert('测试请求失败: ' + err);
        });
}; 