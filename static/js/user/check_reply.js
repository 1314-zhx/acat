// 防 XSS：转义 HTML 特殊字符
function escapeHtml(text) {
    if (typeof text !== 'string') return String(text);
    return text
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}

// 页面加载完成后立即请求数据
document.addEventListener('DOMContentLoaded', () => {
    const messageEl = document.getElementById('message');
    const lettersEl = document.getElementById('letters');

    // 显示加载状态
    messageEl.innerHTML = '<div class="loading">加载中...</div>';

    // 请求当前路径的数据（即 /user/auth/check_reply）
    fetch(window.location.pathname, {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
        },
        credentials: 'same-origin' // 携带 cookie（用于身份验证）
    })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            // 假设成功响应结构为：{ status: 200, data: [...], msg: "..." }
            if (data.status === 200 && Array.isArray(data.data)) {
                const letters = data.data;

                if (letters.length === 0) {
                    messageEl.innerHTML = '<div class="empty">暂无回信</div>';
                    lettersEl.innerHTML = '';
                } else {
                    messageEl.innerHTML = '';
                    lettersEl.innerHTML = letters.map(letter => `
                        <div class="letter">
                            <div class="letter-sender">来自：${escapeHtml(letter.send_name || '系统')}</div>
                            <div class="letter-title">${escapeHtml(letter.title)}</div>
                            <div class="letter-content">${escapeHtml(letter.content)}</div>
                        </div>
                    `).join('');
                }
            } else {
                throw new Error(data.msg || '获取回信失败');
            }
        })
        .catch(error => {
            console.error('加载回信失败:', error);
            messageEl.innerHTML = `<div class="error">加载失败：${escapeHtml(error.message)}</div>`;
            lettersEl.innerHTML = '';
        });
});