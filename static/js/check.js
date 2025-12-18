// 格式化 Unix 时间戳（秒）为中文时间字符串
function formatTime(timestamp) {
    const date = new Date(timestamp * 1000); // 转为毫秒
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        hour12: false
    }).replace(/\//g, '-');
}

// 提取日期部分（如 "2025-01-01"）
function extractDate(timeStr) {
    const parts = timeStr.split(' ');
    return parts[0] || '—';
}

function renderLoading() {
    document.getElementById('content').innerHTML = '<div class="loading">加载中...</div>';
}

function renderError(message) {
    document.getElementById('content').innerHTML = `
        <div class="error">
            <div>请求失败</div>
            <div style="margin-top: 10px; font-size: 14px;">${message || '未知错误'}</div>
        </div>
    `;
}

function renderEmpty() {
    document.getElementById('content').innerHTML = `
        <div class="empty">
            <div>暂无面试安排</div>
            <div style="margin-top: 10px; font-size: 14px; color: #868e96;">
                请关注后续消息
            </div>
        </div>
    `;
}

function renderTable(slots) {
    if (!Array.isArray(slots) || slots.length === 0) {
        renderEmpty();
        return;
    }

    let rows = '';
    slots.forEach(slot => {
        if (!Array.isArray(slot) || slot.length < 2) return;

        const startTs = slot[0];
        const endTs = slot[1];

        const startStr = formatTime(startTs);
        const endStr = formatTime(endTs);
        const dateOnly = extractDate(startStr);

        rows += `
            <tr>
                <td>${startStr}</td>
                <td>${endStr}</td>
                <td>${dateOnly}</td>
            </tr>
        `;
    });

    document.getElementById('content').innerHTML = `
        <table>
            <thead>
                <tr>
                    <th>开始时间</th>
                    <th>结束时间</th>
                    <th>面试日期</th>
                </tr>
            </thead>
            <tbody>
                ${rows}
            </tbody>
        </table>
    `;
}

async function fetchAndRenderSchedule() {
    renderLoading();
    try {
        const response = await fetch('/check');
        if (!response.ok) {
            // HTTP 层错误（如 500、404）
            throw new Error(`HTTP ${response.status}`);
        }

        const data = await response.json();
        if (data.status !== 200) {
            renderError(data.msg || '服务返回错误');
            return;
        }

        // 业务成功，渲染 data（renderTable 会处理空数组）
        renderTable(data.data);

    } catch (err) {
        console.error('获取面试安排失败:', err);
        renderError('网络错误，请检查后端服务是否运行');
    }
}
document.addEventListener('DOMContentLoaded', () => {
    fetchAndRenderSchedule();
});