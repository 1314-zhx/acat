let selectedSlotId = null;

// 加载时段列表
async function loadSlots() {
    const container = document.getElementById('slotsContainer');
    try {
        const res = await fetch('/check');
        if (!res.ok) throw new Error('请求失败');
        const data = await res.json();

        if (data.status !== 200 || !Array.isArray(data.data)) {
            container.innerHTML = '<p class="error">时段数据格式错误</p>';
            return;
        }

        const rawSlots = data.data; // [[start, end, id, round, num, max_num], ...]
        if (rawSlots.length === 0) {
            container.innerHTML = `
                <div class="slot-card" style="text-align:center; opacity:0.8; cursor:default;">
                    暂无可用面试时段
                </div>
            `;
            return;
        }

        const now = Math.floor(Date.now() / 1000); // 当前时间戳（秒）
        const slots = rawSlots
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
            container.innerHTML = `
                <div class="slot-card" style="text-align:center; opacity:0.8; cursor:default;">
                    暂无可用面试时段
                </div>
            `;
            return;
        }

        // 格式化函数
        function formatTimestamp(ts) {
            const date = new Date(ts * 1000);
            return date.toLocaleString('zh-CN', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit'
            }).replace(/\//g, '-');
        }

        function formatTimeOnly(ts) {
            const date = new Date(ts * 1000);
            return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
        }

        container.innerHTML = slots.map(slot => {
            const isFull = slot.num >= slot.max_num;
            const start = formatTimestamp(slot.start_time);
            const end = formatTimeOnly(slot.end_time);
            return `
            <div class="slot-card ${isFull ? 'disabled' : ''}"
                 onclick="${isFull ? '' : `openModal(${slot.id})`}">
                <div class="slot-header">
                    <span>${start} - ${end}</span>
                    <span>${isFull ? '已满' : `${slot.num}/${slot.max_num}`}</span>
                </div>
                <div class="slot-footer">
                    轮次：第 ${slot.round} 轮
                </div>
            </div>
        `;
        }).join('');
    } catch (err) {
        console.error('加载时段失败:', err);
        container.innerHTML = '<p class="error">加载失败，请刷新重试</p>';
    }
}

// 打开报名弹窗
function openModal(slotId) {
    selectedSlotId = slotId;
    document.getElementById('modalName').value = '';
    document.getElementById('modalDirection').value = '0';
    document.getElementById('modalTitle').textContent = `报名时段 #${slotId}`;
    document.getElementById('bookingModal').style.display = 'block';
}

// 关闭弹窗
document.querySelector('.close').onclick = () => {
    document.getElementById('bookingModal').style.display = 'none';
};
window.onclick = (e) => {
    if (e.target === document.getElementById('bookingModal')) {
        document.getElementById('bookingModal').style.display = 'none';
    }
};

// 显示消息
function showMessage(el, text, isSuccess) {
    el.textContent = text;
    el.className = `message ${isSuccess ? 'success' : 'error'}`;
    el.style.display = 'block';
    setTimeout(() => {
        el.style.display = 'none';
    }, 4000);
}

// 提交报名
document.getElementById('submitBooking').addEventListener('click', async () => {
    const name = document.getElementById('modalName').value.trim();
    const direction = parseInt(document.getElementById('modalDirection').value);
    const btn = document.getElementById('submitBooking');
    const msgEl = document.getElementById('modalMessage');

    if (!name) {
        showMessage(msgEl, '姓名不能为空', false);
        return;
    }

    btn.disabled = true;
    btn.textContent = '提交中...';

    try {
        const response = await fetch('/user/auth/signup', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                name: name,
                direction: direction,
                slot_id: selectedSlotId
            })
        });

        const data = await response.json();
        let errorMsg = data.msg || '报名失败，请稍后重试';

        if (response.ok && data.status === 200 && data.msg === 'Success') {
            showMessage(msgEl, '报名成功！请留意通知。', true);
            setTimeout(() => {
                document.getElementById('bookingModal').style.display = 'none';
                loadSlots(); // 刷新人数
            }, 1500);
        } else {
            if (errorMsg.includes('截至')) {
                errorMsg = '面试时间已截至，无法报名。';
            } else if (errorMsg.includes('已报满') || errorMsg.includes('已满')) {
                errorMsg = '该时段名额已满，请选择其他时间。';
            } else if (errorMsg.includes('重复') || errorMsg.includes('已报名')) {
                errorMsg = '您已报名过该时段。';
            }
            showMessage(msgEl, errorMsg, false);
        }
    } catch (err) {
        console.error('网络错误:', err);
        showMessage(msgEl, '网络错误，请检查后重试。', false);
    } finally {
        btn.disabled = false;
        btn.textContent = '确认报名';
    }
});

// 初始化
loadSlots();