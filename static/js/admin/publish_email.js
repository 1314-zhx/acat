let currentRound = null;
let currentUserForModal = null;

// 通用 POST JSON 请求函数
async function postJSON(url, data) {
    const response = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });
    return await response.json();
}

// 获取通过用户列表
async function fetchUsers(round) {
    try {
        const res = await postJSON('/admin/get_pass_user', { round });
        if (res.status === 200) {
            renderUserList(res.data || [], round);
        } else {
            alert('获取用户失败: ' + res.msg);
        }
    } catch (err) {
        console.error(err);
        alert('请求失败，请检查网络或后端服务');
    }
}

// 渲染用户列表（使用事件委托，避免 onclick 字符串问题）
function renderUserList(users, round) {
    currentRound = round;
    document.getElementById('user-count').textContent = users.length;
    const container = document.getElementById('users-container');
    container.innerHTML = '';

    // 转义函数（防 XSS）
    function escapeHtml(str) {
        return String(str).replace(/[<>"'&]/g, match =>
            ({ '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;', '&': '&amp;' }[match])
        );
    }

    users.forEach(user => {
        const card = document.createElement('div');
        card.className = 'user-card';
        card.innerHTML = `
      <div class="user-info">
        <strong>ID:</strong> ${user.id}<br/>
        <strong>姓名:</strong> ${escapeHtml(user.name)}<br/>
        <strong>邮箱:</strong> ${escapeHtml(user.email)}
      </div>
      <div class="actions">
        <button class="btn btn-warning"
                data-action="default"
                data-id="${user.id}"
                data-name="${escapeHtml(user.name)}"
                data-email="${escapeHtml(user.email)}">
          默认发件格式
        </button>
        <button class="btn"
                data-action="custom"
                data-id="${user.id}"
                data-name="${escapeHtml(user.name)}"
                data-email="${escapeHtml(user.email)}">
          自定义发件格式
        </button>
      </div>
    `;
        container.appendChild(card);
    });

    // 事件委托：统一处理按钮点击
    container.addEventListener('click', function(e) {
        const btn = e.target.closest('button[data-action]');
        if (!btn) return;

        const action = btn.dataset.action;
        const id = parseInt(btn.dataset.id);
        const name = btn.dataset.name;
        const email = btn.dataset.email;

        if (action === 'default') {
            sendDefault(id, name, email);
        } else if (action === 'custom') {
            openCustom(id, name, email);
        }
    });

    // 切换视图
    document.getElementById('round-select').style.display = 'none';
    document.getElementById('user-list').style.display = 'block';
}

// 返回轮次选择
function backToRoundSelect() {
    document.getElementById('round-select').style.display = 'block';
    document.getElementById('user-list').style.display = 'none';
}

// 发送默认邮件
async function sendDefault(userId, name, email) {
    const payload = {
        user_id: userId,
        name: name,
        email: email,
        round: currentRound,
        customize: false,
        content: ""
    };
    try {
        const res = await postJSON('/admin/post_email', payload);
        if (res.status === 200) {
            alert(`默认邮件已发送给 ${name}`);
        } else {
            alert('发送失败: ' + res.msg);
        }
    } catch (err) {
        console.error(err);
        alert('网络错误');
    }
}

// 打开自定义邮件弹窗
function openCustom(userId, name, email) {
    currentUserForModal = { userId, name, email };
    document.getElementById('modal-username').textContent = name;
    document.getElementById('custom-content').value = '';
    document.getElementById('custom-modal').style.display = 'flex';
}

// 关闭弹窗
function closeModal() {
    document.getElementById('custom-modal').style.display = 'none';
    currentUserForModal = null;
}

// 点击背景关闭弹窗
function closeModalIfBackground(e) {
    if (e.target.id === 'custom-modal') {
        closeModal();
    }
}

// 发送自定义邮件
async function sendCustomEmail() {
    if (!currentUserForModal) {
        alert('用户信息异常');
        return;
    }
    const content = document.getElementById('custom-content').value.trim();
    if (!content) {
        alert('邮件内容不能为空');
        return;
    }

    const { userId, name, email } = currentUserForModal;
    const payload = {
        user_id: userId,
        name: name,
        email: email,
        round: currentRound,
        customize: true,
        content: content
    };

    try {
        const res = await postJSON('/admin/post_email', payload);
        if (res.status === 200) {
            alert(`自定义邮件已发送给 ${name}`);
            closeModal();
        } else {
            alert('发送失败: ' + res.msg);
        }
    } catch (err) {
        console.error(err);
        alert('网络错误');
    }
}

// 绑定弹窗点击背景关闭（冗余但安全）
document.addEventListener('DOMContentLoaded', () => {
    const modal = document.getElementById('custom-modal');
    if (modal) {
        modal.addEventListener('click', function(e) {
            if (e.target === this) closeModal();
        });
    }
});