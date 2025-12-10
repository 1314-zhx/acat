// 🎯 顶部 Toast 提示函数
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
    void toast.offsetWidth; // 强制重排
    toast.classList.add('show');

    // 1.8 秒后自动隐藏
    setTimeout(() => {
        toast.classList.remove('show');
    }, 1800);
}

// 表单提交处理
document.getElementById('registerForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const formData = new FormData(this);
    const data = Object.fromEntries(formData.entries());

    // 前端校验
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

    try {
        const response = await fetch('/user/register', {
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
            }),
        });

        const result = await response.json();

        if (result.status === 200) {
            showToast('注册成功，正在自动登录...', true);

            // 自动登录
            try {
                const loginRes = await fetch('/user/login', {
                    method: 'POST',
                    credentials: 'include', // ⚠️ 关键！带上 Cookie
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        Phone: data.phone,      // 注意字段名是否匹配后端
                        Password: data.password
                    }),
                });

                const loginResult = await loginRes.json();

                if (loginResult.status === 200) {
                    showToast('注册成功，跳转中...', true);
                    setTimeout(() => {
                        window.location.href = '/user/center'; // 或你的真实路径
                    }, 1500);
                } else {
                    showToast('自动登录失败: ' + (loginResult.msg || '未知错误'), false);
                    setTimeout(() => {
                        window.location.href = '/login.html';
                    }, 2000);
                }
            } catch (loginErr) {
                console.error('自动登录失败:', loginErr);
                showToast('自动登录失败，请手动登录', false);
                setTimeout(() => {
                    window.location.href = '/login.html';
                }, 2000);
            }
        } else {
            showToast(result.msg || '注册失败', false);
        }
    } catch (err) {
        console.error('网络请求失败:', err);
        showToast('网络错误，请稍后再试', false);
    }
});