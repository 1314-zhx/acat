// 显示顶部提示
function showToast(message, isSuccess = true) {
    let toast = document.getElementById('js-toast');
    if (!toast) {
        toast = document.createElement('div');
        toast.id = 'js-toast';
        toast.className = 'toast';
        document.body.appendChild(toast);
    }
    toast.textContent = message;
    toast.className = `toast ${isSuccess ? 'success' : 'error'}`;
    toast.classList.remove('show');
    void toast.offsetWidth;
    toast.classList.add('show');

    setTimeout(() => {
        toast.classList.remove('show');
    }, 2000);
}

// 表单提交：不再检查前端 token，交由后端验证
document.getElementById('resultForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const roundInput = document.getElementById('round');
    const submitBtn = document.getElementById('submitBtn');

    const round = parseInt(roundInput.value.trim(), 10);
    if (isNaN(round) || round < 1 || round > 2) {
        showToast('请输入有效的轮次（1-2）', false);
        return;
    }

    submitBtn.disabled = true;
    submitBtn.textContent = '查询中...';

    try {
        const res = await fetch('/user/auth/result', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                // 注意：如果 token 是通过 Cookie 自动携带的，就不需要手动加 header
                // 如果后端仍要求从 header 读 token，则需从 localStorage 取（但你说不用检查，所以这里先不加）
            },
            body: JSON.stringify({ round }),
            credentials: 'include' // 如果 token 在 Cookie 中，必须加这行！
        });

        const data = await res.json();

        if (data.status === 200) {
            showToast(data.data || '查询成功', true);
            alert(`面试结果：${data.data}`);
        } else {
            // 假设后端约定：status 401 或 msg 包含 "登录" 表示未认证
            if (res.status === 401 || data.status === 401 || /登录|认证|token/i.test(data.msg || '')) {
                showToast('请先登录', false);
                setTimeout(() => {
                    window.location.href = '/login.html';
                }, 1500);
            } else {
                showToast(data.msg || '查询失败', false);
            }
        }
    } catch (err) {
        console.error('请求失败:', err);
        showToast('网络错误，请重试', false);
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = '查询结果';
    }
});