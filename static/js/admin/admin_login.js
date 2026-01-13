document.getElementById('adminLoginForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const phone = document.getElementById('phone').value.trim();
    const password = document.getElementById('password').value;
    const errorMsg = document.getElementById('errorMessage');
    const loginBtn = document.getElementById('loginBtn');

    // 简单校验
    const reg = /^[a-zA-Z0-9-_]{6,20}$/
    const num = /^[0-9]{11}$/
    if (!num.test(phone)) {
        errorMsg.textContent = '请输入有效的11位手机号';
        return;
    }else if (!reg.test(password)) {
        errorMsg.textContent = '请输入6-20位字符';
        return;
    }

    // 清除旧错误
    errorMsg.textContent = '';
    loginBtn.disabled = true;
    loginBtn.textContent = '登录中...';

    try {
        const response = await fetch('/admin_login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                phone: phone,
                password: password,
            }),
        });

        const data = await response.json();

        if (response.ok && (data.status === 200 || data.Status === 200)) {
            // ✅ 登录成功，跳转到管理员操作中心
            window.location.href = '/admin/center';
        } else {
            // 后端返回错误
            const msg = data.msg || data.Msg || data.error || '登录失败，请检查手机号或密码';
            errorMsg.textContent = msg;
        }
    } catch (err) {
        console.error('Login error:', err);
        errorMsg.textContent = '网络错误，请稍后重试';
    } finally {
        loginBtn.disabled = false;
        loginBtn.textContent = '登录';
    }
});