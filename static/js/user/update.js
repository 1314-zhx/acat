let selectedSlotIdForUpdate = null;
let currentUserSlotId = null;

function formatDateTime(ts) {
    const d = new Date(ts * 1000);
    return d.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    }).replace(/\//g, '-');
}

function formatTimeOnly(ts) {
    return new Date(ts * 1000).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
}

// 1. 加载当前用户的预约
async function loadCurrentUserSlot() {
    const container = document.getElementById('currentSlotContainer');
    try {
        const res = await fetch('/user/auth/my_slot');

        // 情况1: 未登录
        if (res.status === 401) {
            container.innerHTML = '<p class="error">请先登录</p>';
            return;
        }

        // 情况2: 其他 HTTP 错误（500、404 等）
        if (!res.ok) {
            container.innerHTML = '<p class="error">加载失败，请重试</p>';
            return;
        }

        // 情况3: 成功响应，解析 JSON
        const data = await res.json();

        // 成功且有数据 → 显示预约
        if (data.status === 200 && data.data) {
            const slot = data.data;
            currentUserSlotId = slot.id;

            const startTimeTs = new Date(slot.start_time).getTime() / 1000;
            const endTimeTs = new Date(slot.end_time).getTime() / 1000;

            container.innerHTML = `
                <div class="current-card">
                    <div class="current-info">
                        <h2>${formatDateTime(startTimeTs)} - ${formatTimeOnly(endTimeTs)}</h2>
                        <p>轮次：第 ${slot.round} 轮</p>
                    </div>
                    <div class="actions">
                        <button class="btn-update" onclick="showAllSlots()">更新</button>
                        <button class="btn-delete" onclick="deleteCurrentAssignment()">删除</button>
                    </div>
                </div>
            `;
        }
        // 成功但无数据 → 用户确实没预约
        else {
            container.innerHTML = '<p class="info">您尚未预约面试</p>'; // 建议用 .info 样式（非 error）
        }
    } catch (err) {
        console.error('加载当前预约失败:', err);
        container.innerHTML = '<p class="error">网络错误，请检查连接</p>';
    }
}

// 2. 显示所有可选时段（用于更新）
async function showAllSlots() {
    const container = document.getElementById('allSlotsContainer');
    container.style.display = 'block';
    container.innerHTML = '<p class="loading">加载中...</p>';

    try {
        const res = await fetch('/check');
        if (!res.ok) throw new Error('加载失败');
        const data = await res.json();

        if (data.status !== 200 || !Array.isArray(data.data)) {
            container.innerHTML = '<p class="error">时段数据异常</p>';
            return;
        }

        const now = Math.floor(Date.now() / 1000);
        const slots = data.data
            .filter(arr => Array.isArray(arr) && arr.length >= 6)
            .map(([startTs, endTs, id, round, num, maxNum]) => ({
                id,
                round,
                start_time: startTs,
                end_time: endTs,
                num,
                max_num: maxNum
            }))
            .filter(slot => slot.start_time > now);

        if (slots.length === 0) {
            container.innerHTML = '<p>暂无可用面试时段</p>';
            return;
        }

        container.innerHTML = slots.map(slot => {
            const isFull = slot.num >= slot.max_num;
            return `
                <div class="slot-card ${isFull ? 'disabled' : ''}"
                     onclick="${isFull ? '' : `openUpdateModal(${slot.id})`}">
                    <div class="slot-header">
                        <span>${formatDateTime(slot.start_time)} - ${formatTimeOnly(slot.end_time)}</span>
                        <span>${isFull ? '已满' : `${slot.num}/${slot.max_num}`}</span>
                    </div>
                    <div>轮次：第 ${slot.round} 轮</div>
                </div>`;
        }).join('');
    } catch (err) {
        console.error(err);
        container.innerHTML = '<p class="error">加载时段失败</p>';
    }
}

// 3. 删除当前预约
async function deleteCurrentAssignment() {
    if (!confirm('确定要取消当前预约吗？')) return;
    if (!currentUserSlotId) {
        alert('未获取到当前预约信息');
        return;
    }

    try {
        const res = await fetch('/user/auth/update', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                slot_id: currentUserSlotId,
                is_delete: 1
            })
        });

        const data = await res.json();
        if (res.ok && data.status === 200) {
            alert('✅ 预约已取消');
        } else {
            alert(data.msg || '取消失败');
        }
    } catch (err) {
        alert('网络错误');
    }
}

// 4. 弹出更新表单
function openUpdateModal(slotId) {
    selectedSlotIdForUpdate = slotId;
    document.getElementById('modalName').value = '';
    document.getElementById('modalDirection').value = '0';
    document.getElementById('modalTitle').textContent = `更新至时段 #${slotId}`;
    document.getElementById('updateModal').style.display = 'block';
}

// 关闭弹窗
document.querySelector('.close').onclick = () => {
    document.getElementById('updateModal').style.display = 'none';
};
window.onclick = (e) => {
    if (e.target === document.getElementById('updateModal')) {
        document.getElementById('updateModal').style.display = 'none';
    }
};

// 5. 提交更新（更换时段或保留原时段）
document.getElementById('submitUpdate').addEventListener('click', async () => {
    const name = document.getElementById('modalName').value.trim();
    const direction = parseInt(document.getElementById('modalDirection').value);
    const msgEl = document.getElementById('modalMessage');

    if (!name) {
        showMessage(msgEl, '姓名不能为空', false);
        return;
    }

    const btn = document.getElementById('submitUpdate');
    btn.disabled = true;
    btn.textContent = '提交中...';

    try {
        const res = await fetch('/user/auth/update', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                slot_id: selectedSlotIdForUpdate,
                name: name,
                direction: direction,
                is_delete: 0
            })
        });

        const data = await res.json();
        if (res.ok && data.status === 200) {
            showMessage(msgEl, '更新成功！', true);
            setTimeout(() => {}, 1500);
        } else {
            let msg = data.msg || '操作失败';
            if (msg.includes('已满')) msg = '所选时段已满';
            showMessage(msgEl, msg, false);
        }
    } catch (err) {
        showMessage(msgEl, '网络错误', false);
    } finally {
        btn.disabled = false;
        btn.textContent = '确认提交';
    }
});

function showMessage(el, text, isSuccess) {
    el.textContent = text;
    el.className = `message ${isSuccess ? 'success' : 'error'}`;
    el.style.display = 'block';
    setTimeout(() => { el.style.display = 'none'; }, 3000);
}

// 初始化
loadCurrentUserSlot();