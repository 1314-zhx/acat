// 创建并管理全局 Toast 元素
function showToast(message, isSuccess = true) {
    let toast = document.getElementById('js-toast');
    if (!toast) {
        toast = document.createElement('div');
        toast.id = 'js-toast';
        toast.className = 'toast';
        document.body.appendChild(toast);
    }

    toast.textContent = message;
    toast.style.background = isSuccess ? '#52c41a' : '#f5222d';

    // 触发动画
    toast.classList.remove('show');
    void toast.offsetWidth; // 强制重排以重置动画
    toast.classList.add('show');

    setTimeout(() => {
        toast.classList.remove('show');
    }, 1500);
}

// 登录逻辑
async function login() {
    const phoneInput = document.getElementById('phone');
    const passwordInput = document.getElementById('password');

    const phone = phoneInput.value.trim();
    const password = passwordInput.value;

    const reg = /^[a-zA-Z0-9_-]{6,20}$/
    const num = /^[0-9]{11}$/
    if (!num.test(phone)) {
       showToast('请输入11位手机号', false);
        return;
    }else if (!reg.test(password)) {
        showToast('请输入6-20位字符11', false);
        return;
    }

    const btn = document.getElementById('loginBtn');
    btn.disabled = true;
    btn.textContent = '登录中...';

    try {
        const res = await fetch('/user/login', {
            method: 'POST',
            credentials: 'include', // 重要：携带 Cookie（用于会话）
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ Phone: phone, Password: password })
        });

        let data;
        try {
            data = await res.json();
        } catch (e) {
            console.error('响应非JSON:', await res.text());
            showToast('服务器异常，请稍后再试', false);
            return;
        }

        if (res.ok && data.status === 200) {
            showToast('登录成功，跳转中...', true);
            setTimeout(() => {
                window.location.href = '/user/center';
            }, 1500);
        } else {
            showToast(data.msg || '登录失败', false);
        }
    } catch (err) {
        console.error('网络请求失败:', err);
        showToast('网络错误，请检查连接', false);
    } finally {
        btn.disabled = false;
        btn.textContent = '登录';
    }
}

// 绑定事件
document.addEventListener('DOMContentLoaded', () => {
    // 点击按钮登录
    document.getElementById('loginBtn')?.addEventListener('click', login);

    // 回车键触发登录
    document.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            login();
        }
    });
});