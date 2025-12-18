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
    toast.classList.remove('show');
    void toast.offsetWidth;
    toast.classList.add('show');
    setTimeout(() => toast.classList.remove('show'), 1800);
}

// 发送验证码（调用 /register 接口）
async function sendRegisterCode(data) {
    try {
        const res = await fetch('/user/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                stu_id: data.stuId,
                name: data.name,
                phone: data.phone,
                email: data.email,
                gender: parseInt(data.gender),
                direction: parseInt(data.direction),
                // 注意：不传 password/rePassword/code
            })
        });
        const result = await res.json();
        if (result.status === 200) {
            showToast('验证码已发送，请查收邮箱', true);
            document.getElementById('codeGroup').style.display = 'block';
            document.getElementById('code')?.focus();
            return true;
        } else {
            showToast(result.msg || '发送失败', false);
            return false;
        }
    } catch (err) {
        console.error('发送验证码失败:', err);
        showToast('网络错误，请稍后再试', false);
        return false;
    }
}

// 完成注册（调用 /complete_register 接口）
async function completeRegister(data) {
    try {
        const res = await fetch('/user/complete_register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                stu_id: data.stuId,
                name: data.name,
                phone: data.phone,
                email: data.email,
                gender: parseInt(data.gender),
                direction: parseInt(data.direction),
                password: data.password,
                re_password: data.rePassword,
                code: data.code,
            })
        });
        const result = await res.json();
        if (result.status === 200) {
            showToast('注册成功，正在自动登录...', true);

            // 自动登录
            try {
                const loginRes = await fetch('/user/login', {
                    method: 'POST',
                    credentials: 'include',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        Phone: data.phone,
                        Password: data.password
                    })
                });
                const loginResult = await loginRes.json();
                if (loginResult.status === 200) {
                    showToast('跳转中...', true);
                    setTimeout(() => window.location.href = '/user/center', 1500);
                } else {
                    showToast('自动登录失败，请手动登录', false);
                    setTimeout(() => window.location.href = '/login.html', 2000);
                }
            } catch (e) {
                console.error('自动登录失败:', e);
                showToast('自动登录失败，请手动登录', false);
                setTimeout(() => window.location.href = '/login.html', 2000);
            }
        } else {
            showToast(result.msg || '注册失败', false);
        }
    } catch (err) {
        console.error('注册请求失败:', err);
        showToast('网络错误，请稍后再试', false);
    }
}

// 重新发送验证码
document.getElementById('resendBtn')?.addEventListener('click', async () => {
    const form = document.getElementById('registerForm');
    const formData = new FormData(form);
    const data = Object.fromEntries(formData.entries());

    // 基础校验
    if (!data.stuId || !data.name || !data.phone || !data.email || !data.gender || !data.direction) {
        showToast('请先填写完整信息', false);
        return;
    }
    if (data.phone.length !== 11 || isNaN(data.phone)) {
        showToast('手机号格式不正确', false);
        return;
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(data.email)) {
        showToast('邮箱格式不正确', false);
        return;
    }

    await sendRegisterCode(data);
});

// 表单提交主逻辑
document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const form = e.target;
    const formData = new FormData(form);
    const data = Object.fromEntries(formData.entries());

    // 必填校验
    if (!data.stuId || !data.name || !data.phone || !data.email ||
        !data.gender || !data.direction || !data.password || !data.rePassword) {
        showToast('请填写所有必填字段', false);
        return;
    }

    if (data.password !== data.rePassword) {
        showToast('两次密码不一致', false);
        return;
    }

    if (data.phone.length !== 11 || isNaN(data.phone)) {
        showToast('手机号必须为11位数字', false);
        return;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(data.email)) {
        showToast('邮箱格式不正确', false);
        return;
    }

    if (data.password.length < 6) {
        showToast('密码至少6位', false);
        return;
    }

    // 判断阶段
    if (!data.code) {
        // 第一阶段：发送验证码
        await sendRegisterCode(data);
    } else {
        // 第二阶段：完成注册
        await completeRegister(data);
    }
});