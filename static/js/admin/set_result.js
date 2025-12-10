let currentRound = null;
let currentSlotId = null;

function formatTime(unix) {
    const date = new Date(unix * 1000);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${year}-${month}-${day} ${hours}:${minutes}`;
}

async function loadSlots() {
    try {
        const res = await fetch('/check');
        if (!res.ok) throw new Error('HTTP error ' + res.status);

        const response = await res.json();

        let data;
        if (typeof response === 'object' && response !== null) {
            if (Array.isArray(response.data)) {
                data = response.data;
            } else if (Array.isArray(response)) {
                data = response;
            } else {
                throw new Error('无法识别的数据格式');
            }
        } else {
            throw new Error('响应不是对象或数组');
        }

        const container = document.getElementById('slotsContainer');
        if (data.length === 0) {
            container.innerHTML = '<div class="loading">暂无面试场次</div>';
            return;
        }

        container.innerHTML = '';
        data.forEach(slot => {
            if (!Array.isArray(slot) || slot.length < 6) {
                console.warn('跳过无效 slot:', slot);
                return;
            }

            const [startTime, endTime, slotId, round, num, maxNum] = slot.map(Number);

            const card = document.createElement('div');
            card.className = 'slot-card';
            card.innerHTML = `
        <h3>第 ${round} 轮面试</h3>
        <div class="slot-info">
          <div>时间: ${formatTime(startTime)} - ${formatTime(endTime)}</div>
          <div>人数: ${num}/${maxNum}</div>
          <div>场次ID: ${slotId}</div>
        </div>
      `;
            card.addEventListener('click', () => {
                // 高亮当前卡片
                document.querySelectorAll('.slot-card').forEach(el => el.classList.remove('selected'));
                card.classList.add('selected');

                // 更新全局状态并加载用户
                currentRound = round;
                currentSlotId = slotId;

                fetchUsers(round, slotId);
            });
            container.appendChild(card);
        });
    } catch (err) {
        console.error('加载 slots 失败:', err);
        document.getElementById('slotsContainer').innerHTML =
            `<div class="loading">加载失败: ${err.message}</div>`;
    }
}

async function fetchUsers(round, slotId) {
    try {
        const res = await fetch('/admin/set_result', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ round, slot_id: slotId })
        });

        if (!res.ok) throw new Error('获取用户列表失败');
        const result = await res.json();

        if (result.status !== 200) {
            throw new Error(result.msg || '请求失败');
        }

        renderUsers(result.data);
        document.getElementById('usersSection').style.display = 'block';
    } catch (err) {
        alert('加载用户失败: ' + err.message);
        console.error(err);
    }
}

function renderUsers(users) {
    const list = document.getElementById('usersList');
    if (!users || users.length === 0) {
        list.innerHTML = '<div class="loading">该场次暂无用户</div>';
        return;
    }

    list.innerHTML = '';
    users.forEach(user => {
        const genderText = user.gender === 1 ? '男' : user.gender === 2 ? '女' : '未知';
        const item = document.createElement('div');
        item.className = 'user-item';
        item.innerHTML = `
      <div class="user-info">
        <div><strong>${user.name}</strong></div>
        <div>学号: ${user.stuId || '--'}</div>
        <div>电话: ${user.phone || '--'}</div>
        <div>性别: ${genderText}</div>
      </div>
      <div class="user-actions">
        <button class="btn-pass" onclick="handleAction(event, ${user.id}, 'pass')">通过</button>
        <button class="btn-reject" onclick="handleAction(event, ${user.id}, 'reject')">淘汰</button>
      </div>
    `;
        list.appendChild(item);
    });
}

// 全局函数（因 HTML 中使用了 onclick）
window.handleAction = async function(event, userId, action) {
    if (currentSlotId === null || currentRound === null) {
        alert('请先选择一个面试场次');
        return;
    }

    const pass = action === 'pass' ? 1 : 0;
    const actionText = pass === 1 ? '通过' : '淘汰';

    try {
        const response = await fetch('/admin/set_pass', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                user_id: userId,
                slot_id: currentSlotId,
                round: currentRound,
                pass: pass
            })
        });

        const result = await response.json();
        if (!response.ok || (result.status && result.status !== 200)) {
            throw new Error(result.msg || '操作失败');
        }

        const btn = event.target;
        btn.disabled = true;
        btn.textContent = actionText + '✓';
        btn.style.opacity = '0.75';
    } catch (err) {
        console.error('设置结果失败:', err);
        alert(`设置${actionText}失败: ${err.message}`);
    }
};

document.addEventListener('DOMContentLoaded', () => {
    loadSlots();
});