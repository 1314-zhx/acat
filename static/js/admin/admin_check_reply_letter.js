// 全局变量
let currentLetter = null;
const lettersContainer = document.getElementById('letters-container');
const replyModal = document.getElementById('reply-modal');

// 工具函数：显示提示
function showAlert(message) {
    alert(message);
}

// 更新字数统计
function updateCharCount() {
    const title = document.getElementById('reply-title').value;
    const content = document.getElementById('reply-content').value;
    document.getElementById('title-count').textContent = title.length;
    document.getElementById('content-count').textContent = content.length;
}

// 获取全部信件
async function fetchLetters() {
    try {
        const token = localStorage.getItem('admin_token');
        const headers = { 'Content-Type': 'application/json' };
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }

        const res = await fetch('/admin/mailbox', {
            method: 'POST',
            headers,
            body: JSON.stringify({})
        });

        const data = await res.json();
        if (res.ok && data.status === 200) {
            renderLetters(data.data || []);
        } else {
            showAlert('获取信件失败：' + (data.msg || '未知错误'));
        }
    } catch (err) {
        console.error('网络错误:', err);
        showAlert('网络错误，请检查后端服务是否运行，并确认已登录');
    }
}

// 渲染信件列表
function renderLetters(letters) {
    if (letters.length === 0) {
        lettersContainer.innerHTML = '<div class="no-letters">暂无信件。</div>';
        return;
    }

    lettersContainer.innerHTML = letters.map(letter => `
    <div class="letter-card"
         data-id="${letter.id}"
         data-send-id="${letter.send_id}"
         data-receive-id="${letter.receive_id}">
        <div class="letter-header">发件人：${escapeHtml(letter.send_name || '未知用户')}</div>
        <div class="letter-title">标题：${escapeHtml(letter.title || '无题')}</div>
        <div class="letter-content">内容：${escapeHtml(letter.content || '')}</div>
        <div class="letter-actions">
            <button class="btn-mark-read ${letter.is_read ? 'read' : 'unread'}"
                    data-id="${letter.id}"
                    data-is-read="${letter.is_read}">
                ${letter.is_read ? '取消已读' : '标记已读'}
            </button>
            <button class="btn-reply"
                    data-id="${letter.id}"
                    data-send-id="${letter.send_id}"
                    data-receive-id="${letter.receive_id}">
                回信
            </button>
        </div>
    </div>
`).join('');
}

// 防 XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// 【关键】统一调用 /admin/reply 接口
async function sendToReplyAPI(payload) {
    try {
        const token = localStorage.getItem('admin_token');
        const res = await fetch('/admin/reply', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token
            },
            body: JSON.stringify(payload)
        });

        const result = await res.json().catch(() => ({}));
        if (res.ok && result.status === 200) {
            return { success: true };
        } else {
            return { success: false, msg: result.msg || result.error || '操作失败' };
        }
    } catch (err) {
        return { success: false, msg: '网络错误' };
    }
}

// 切换已读状态（走 /admin/reply）
async function toggleReadStatus(letterId, currentIsRead, receiveId) {
    const newIsRead = !currentIsRead;

    // 构造 payload：content 为空表示“仅标记已读”
    const payload = {
        letter_id: letterId,
        is_read: newIsRead,
        admin_id: receiveId,   // ✅ 改为使用 receive_id
        user_id: 0,
        title: "",
        content: ""
    };

    const result = await sendToReplyAPI(payload);
    if (result.success) {
        const btn = document.querySelector(`.btn-mark-read[data-id="${letterId}"]`);
        if (btn) {
            btn.classList.toggle('read', newIsRead);
            btn.classList.toggle('unread', !newIsRead);
            btn.textContent = newIsRead ? '取消已读' : '标记已读';
            btn.dataset.isRead = newIsRead;
        }
    } else {
        showAlert('操作失败：' + result.msg);
    }
}

// 打开回信弹窗
function openReplyModal(letterId, sendId, receiveId) {
    currentLetter = {
        id: letterId,
        send_id: sendId,
        receive_id: receiveId
    };
    document.getElementById('reply-title').value = '管理员回信';
    document.getElementById('reply-content').value = '';
    updateCharCount();
    replyModal.style.display = 'flex';
}

// 关闭弹窗
function closeReplyModal() {
    replyModal.style.display = 'none';
    currentLetter = null;
}

// 发送回信（走 /admin/reply）
async function sendReply() {
    if (!currentLetter) return;

    const title = document.getElementById('reply-title').value.trim();
    const content = document.getElementById('reply-content').value.trim();

    if (!title || !content) {
        showAlert('标题和正文不能为空');
        return;
    }

    const payload = {
        letter_id: currentLetter.id,
        admin_id: currentLetter.receive_id,   // ✅ 使用 receive_id
        user_id: currentLetter.send_id,
        title: title,
        content: content,
        is_read: false
    };

    const result = await sendToReplyAPI(payload);
    if (result.success) {
        showAlert('回信发送成功！');
        closeReplyModal();
    } else {
        showAlert('发送失败：' + result.msg);
    }
}

// 事件委托
lettersContainer.addEventListener('click', (e) => {
    const card = e.target.closest('.letter-card');
    if (!card) return;

    const letterId = parseInt(card.dataset.id);
    const sendId = parseInt(card.dataset.sendId);
    const receiveId = parseInt(card.dataset.receiveId);

    if (e.target.classList.contains('btn-mark-read')) {
        const isRead = e.target.dataset.isRead === 'true';
        toggleReadStatus(letterId, isRead, receiveId);
    } else if (e.target.classList.contains('btn-reply')) {
        openReplyModal(letterId, sendId, receiveId);
    }
});

// 弹窗事件
document.getElementById('cancel-reply').addEventListener('click', closeReplyModal);
document.getElementById('send-reply').addEventListener('click', sendReply);
document.getElementById('reply-title').addEventListener('input', updateCharCount);
document.getElementById('reply-content').addEventListener('input', updateCharCount);

// 页面加载
document.addEventListener('DOMContentLoaded', () => {
    fetchLetters();
});