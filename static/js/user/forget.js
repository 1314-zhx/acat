// 工具函数：校验邮箱
function isValidAccount(value) {
    const emailRegex = /^[a-zA-Z0-9]+(?:[._-][a-zA-Z0-9]+)*@[a-zA-Z0-9]+(?:[.-][a-zA-Z0-9]+)*\.[a-zA-Z]{2,}$/
    return emailRegex.test(value);
}

// 工具函数：校验密码长度
function isValidPassword(pwd) {
    return pwd.length >= 6 && pwd.length <= 20;
}

// 获取 DOM 元素
const step1 = document.getElementById('step1');
const step2 = document.getElementById('step2');
const displayAccount = document.getElementById('displayAccount');
const storedAccountInput = document.getElementById('storedAccount');

// ========== 第一步：发送验证码 ==========
document.getElementById('accountForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const accountInput = document.getElementById('account');
    const account = accountInput.value.trim();
    const accountGroup = document.getElementById('accountGroup');

    // 校验邮箱
    if (!account || !isValidAccount(account)) {
        accountGroup.classList.add('error');
        return;
    }
    accountGroup.classList.remove('error');

    const btn = document.getElementById('submitBtn');
    btn.disabled = true;
    btn.textContent = '发送中...';

    try {
        const res = await fetch('/user/forget', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ param: account })
        });
        const data = await res.json();

        if (res.ok && data.status === 200) {
            // 存储邮箱到 hidden 字段 + 显示
            displayAccount.textContent = account;
            storedAccountInput.value = account;
            // 切换步骤
            step1.classList.remove('active');
            step2.classList.add('active');
        } else {
            alert(data.msg || '发送失败，请稍后重试');
        }
    } catch (err) {
        console.error('发送验证码出错:', err);
        alert('网络错误，请检查连接');
    } finally {
        btn.disabled = false;
        btn.textContent = '发送验证码';
    }
});

// ========== 第二步：重置密码 ==========
document.getElementById('verifyForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();

    const account = storedAccountInput.value.trim();
    const code = document.getElementById('code').value.trim();
    const pwd = document.getElementById('newPassword').value;
    const confirm = document.getElementById('confirmPassword').value;

    const codeGroup = document.getElementById('codeGroup');
    const passwordGroup = document.getElementById('passwordGroup');
    const confirmGroup = document.getElementById('confirmGroup');

    // 清除旧错误
    codeGroup.classList.remove('error');
    passwordGroup.classList.remove('error');
    confirmGroup.classList.remove('error');

    // 必填校验
    if (!account) {
        alert("邮箱信息丢失，请重新操作");
        return;
    }
    if (!code) {
        codeGroup.classList.add('error');
        return;
    }
    if (!isValidPassword(pwd)) {
        passwordGroup.classList.add('error');
        return;
    }
    if (pwd !== confirm) {
        confirmGroup.classList.add('error');
        return;
    }

    const btn = document.getElementById('verifyBtn');
    btn.disabled = true;
    btn.textContent = '处理中...';

    try {
        const res = await fetch('/user/reset-password', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                account: account,
                code: code,
                new_password: pwd
            })
        });
        const data = await res.json();

        if (res.ok && data.status === 200) {
            document.getElementById('successMessage').textContent = '密码重置成功！正在跳转...';
            setTimeout(() => window.location.href = '/user/login', 2000);
        } else {
            alert(data.msg || '操作失败，请重试');
            // 若提示含“验证码”，高亮验证码框
            if (data.msg && (data.msg.includes('验证码') || data.msg.includes('code'))) {
                codeGroup.classList.add('error');
            }
        }
    } catch (err) {
        console.error('重置密码出错:', err);
        alert('网络错误，请检查连接');
    } finally {
        btn.disabled = false;
        btn.textContent = '重置密码';
    }
});

// 返回上一步
function goBack() {
    step2.classList.remove('active');
    step1.classList.add('active');
}