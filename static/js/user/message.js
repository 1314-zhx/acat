// 方向映射：根据你的 AdminModel.Direction (1=Go, 2=Java, 3=前端)
const DIRECTION_MAP = {
    1: '后端Go',
    2: '后端Java',
    3: '前端'
};

// ✅ 从 /user/auth/show_admin 获取管理员列表
async function loadAdmins() {
    const container = document.getElementById('adminList');
    container.innerHTML = '<div class="loading">加载中...</div>';

    try {
        const response = await fetch('/user/auth/show_admin', {
            method: 'GET',
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`);
        }

        const result = await response.json();

        if (result.status === 200 && Array.isArray(result.data)) {
            renderAdmins(result.data);
        } else {
            throw new Error(result.msg || '未知错误');
        }
    } catch (err) {
        console.error('加载管理员失败:', err);
        container.innerHTML = `<div class="error">加载失败：${err.message}</div>`;
    }
}

function renderAdmins(admins) {
    const container = document.getElementById('adminList');
    if (!admins || admins.length === 0) {
        container.innerHTML = '<div style="text-align:center;color:#999;">暂无管理员可联系</div>';
        return;
    }

    container.innerHTML = admins.map(admin => `
        <div class="admin-card">
            <div class="admin-info">
                <h3>${admin.name}</h3>
                <p>📞 ${admin.phone}</p>
                <p>方向：${DIRECTION_MAP[admin.direction] || '未知'}</p>
            </div>
            <button class="btn-message" onclick="openMessageModal(${admin.aid})">
                私信该管理员
            </button>
        </div>
    `).join('');
}

function openMessageModal(adminId) {
    document.getElementById('receiveId').value = adminId;
    document.getElementById('title').value = "发给管理员的消息";
    document.getElementById('content').value = "";
    updateCharCount();
    document.getElementById('messageModal').style.display = 'block';
}

function closeModal() {
    document.getElementById('messageModal').style.display = 'none';
}

function updateCharCount() {
    const content = document.getElementById('content').value;
    document.getElementById('charCount').textContent = content.length;
}

document.getElementById('messageForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const receiveId = parseInt(document.getElementById('receiveId').value);
    const title = document.getElementById('title').value.trim();
    const content = document.getElementById('content').value.trim();

    if (!title) {
        alert("标题不能为空");
        return;
    }
    if (title.length > 50) {
        alert("标题不能超过50个字");
        return;
    }
    if (!content) {
        alert("正文不能为空");
        return;
    }
    if (content.length > 15) {
        alert("正文不能超过15个字");
        return;
    }

    try {
        // 修正：使用正确的后端 API 路径 /user/auth/conversation
        const response = await fetch('/user/auth/conversation', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({
                receive_id: receiveId,
                title: title,
                content: content
            })
        });

        // 🔍 防御性检查：确保返回的是 JSON
        const contentType = response.headers.get("content-type");
        if (!contentType || !contentType.includes("application/json")) {
            const text = await response.text();
            console.error("服务器返回非JSON内容:", text.substring(0, 200));
            throw new Error("服务器返回异常，请确认已登录");
        }

        const data = await response.json();
        if (response.ok && data.status === 200) {
            alert("消息发送成功！");
            closeModal();
        } else {
            alert("发送失败：" + (data.msg || data.error || "未知错误"));
        }
    } catch (err) {
        console.error("发送失败:", err);
        alert("网络错误，请稍后重试");
    }
});

document.getElementById('content').addEventListener('input', updateCharCount);

window.onclick = function(event) {
    const modal = document.getElementById('messageModal');
    if (event.target === modal) {
        closeModal();
    }
};

document.addEventListener('DOMContentLoaded', () => {
    loadAdmins();
});