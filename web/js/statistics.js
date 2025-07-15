function loadStatistics() {
    fetch('/api/statistics/summary', {
        headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    })
        .then(res => res.json())
        .then(data => {
            renderStatisticsSummary(data.data || {});
        })
        .catch(err => alert('统计摘要加载失败: ' + err));

    fetch('/api/statistics/list', {
        headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    })
        .then(res => res.json())
        .then(data => {
            renderStatisticsTable(data.data || []);
        })
        .catch(err => alert('统计数据加载失败: ' + err));
}

function renderStatisticsSummary(summary) {
    const el = document.getElementById('statisticsSummary');
    el.innerHTML = `
      <div class="row text-center">
        <div class="col">
          <div class="card">
            <div class="card-body">
              <h5 class="card-title">总请求数</h5>
              <p class="card-text fs-3">${summary.total_requests || 0}</p>
            </div>
          </div>
        </div>
        <div class="col">
          <div class="card">
            <div class="card-body">
              <h5 class="card-title">成功数</h5>
              <p class="card-text fs-3 text-success">${summary.success_count || 0}</p>
            </div>
          </div>
        </div>
        <div class="col">
          <div class="card">
            <div class="card-body">
              <h5 class="card-title">失败数</h5>
              <p class="card-text fs-3 text-danger">${summary.fail_count || 0}</p>
            </div>
          </div>
        </div>
      </div>
    `;
}

function renderStatisticsTable(list) {
    const tbody = document.querySelector('#statisticsTable tbody');
    tbody.innerHTML = '';
    list.forEach(item => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
          <td>${item.api_name || ''}</td>
          <td>${item.total || 0}</td>
          <td>${item.success || 0}</td>
          <td>${item.fail || 0}</td>
          <td>${item.avg_duration || 0}</td>
        `;
        tbody.appendChild(tr);
    });
}