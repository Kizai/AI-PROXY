// js/dashboard.js

function loadDashboard() {
    // 加载核心数据
    fetch('/api/dashboard/summary', {
        headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    })
        .then(res => res.json())
        .then(data => {
            renderDashboardSummary(data.data || {});
        })
        .catch(err => alert('仪表板数据加载失败: ' + err));

    // 加载公告（可选）
    fetch('/api/dashboard/notice', {
        headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    })
        .then(res => res.json())
        .then(data => {
            renderDashboardNotice(data.data || '');
        })
        .catch(() => {}); // 公告不是必须的
}

function renderDashboardSummary(summary) {
    const el = document.getElementById('dashboardSummary');
    el.innerHTML = `
      <div class="row text-center">
        <div class="col">
          <div class="card border-primary mb-3">
            <div class="card-body">
              <h5 class="card-title">API总数</h5>
              <p class="card-text fs-3">${summary.api_count || 0}</p>
            </div>
          </div>
        </div>
        <div class="col">
          <div class="card border-success mb-3">
            <div class="card-body">
              <h5 class="card-title">今日请求数</h5>
              <p class="card-text fs-3 text-success">${summary.today_requests || 0}</p>
            </div>
          </div>
        </div>
        <div class="col">
          <div class="card border-info mb-3">
            <div class="card-body">
              <h5 class="card-title">在线状态</h5>
              <p class="card-text fs-3">${summary.online ? '在线' : '离线'}</p>
            </div>
          </div>
        </div>
      </div>
    `;
}

function renderDashboardNotice(notice) {
    const el = document.getElementById('dashboardNotice');
    if (notice) {
        el.innerText = notice;
        el.style.display = '';
    } else {
        el.style.display = 'none';
    }
}