// js/api-config.js

function initApiConfigPage() {
    // 绑定新增按钮
    document.getElementById('btnAddApiConfig').onclick = function () {
        showApiConfigModal();
    };

    // 绑定表单提交
    document.getElementById('apiConfigForm').onsubmit = function (e) {
        e.preventDefault();
        saveApiConfig();
    };

    // 加载API配置列表
    loadApiConfigList();
}

function loadApiConfigList() {
    fetch('/api/config/list', {
        headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    })
        .then(res => res.json())
        .then(data => {
            renderApiConfigTable(data.data || []);
        })
        .catch(err => alert('加载失败: ' + err));
}

function renderApiConfigTable(list) {
    const tbody = document.querySelector('#apiConfigTable tbody');
    tbody.innerHTML = '';
    list.forEach(item => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
          <td>${item.id}</td>
          <td>${item.name}</td>
          <td>${item.target_url}</td>
          <td>${item.desc || ''}</td>
          <td>
            <button class="btn btn-sm btn-info" onclick="showApiConfigModal(${item.id}, '${item.name}', '${item.target_url}', '${item.desc || ''}')">编辑</button>
            <button class="btn btn-sm btn-danger" onclick="deleteApiConfig(${item.id})">删除</button>
          </td>
        `;
        tbody.appendChild(tr);
    });
}

function showApiConfigModal(id = '', name = '', target = '', desc = '') {
    document.getElementById('apiConfigId').value = id;
    document.getElementById('apiConfigName').value = name;
    document.getElementById('apiConfigTarget').value = target;
    document.getElementById('apiConfigDesc').value = desc;
    document.getElementById('apiConfigModalTitle').innerText = id ? '编辑 API 配置' : '新增 API 配置';
    new bootstrap.Modal(document.getElementById('apiConfigModal')).show();
}

function saveApiConfig() {
    const id = document.getElementById('apiConfigId').value;
    const name = document.getElementById('apiConfigName').value;
    const target = document.getElementById('apiConfigTarget').value;
    const desc = document.getElementById('apiConfigDesc').value;
    const url = id ? '/api/config/update' : '/api/config/add';
    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem('token')
        },
        body: JSON.stringify({ id, name, target_url: target, desc })
    })
        .then(res => res.json())
        .then(data => {
            if (data.code === 0) {
                alert('保存成功');
                loadApiConfigList();
                bootstrap.Modal.getInstance(document.getElementById('apiConfigModal')).hide();
            } else {
                alert('保存失败: ' + data.msg);
            }
        })
        .catch(err => alert('请求失败: ' + err));
}

function deleteApiConfig(id) {
    if (!confirm('确定要删除吗？')) return;
    fetch('/api/config/delete', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem('token')
        },
        body: JSON.stringify({ id })
    })
        .then(res => res.json())
        .then(data => {
            if (data.code === 0) {
                alert('删除成功');
                loadApiConfigList();
            } else {
                alert('删除失败: ' + data.msg);
            }
        })
        .catch(err => alert('请求失败: ' + err));
} 