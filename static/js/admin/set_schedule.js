// 获取所有输入元素
const startTimeInput = document.getElementById('start_time');
const endTimeInput = document.getElementById('end_time');
const roundSelect = document.getElementById('round');
const maxNumInput = document.getElementById('max_num');
const jsonOutput = document.getElementById('jsonOutput');
const submitBtn = document.getElementById('submitBtn');
const clearBtn = document.getElementById('clearBtn');
const messageDiv = document.getElementById('message');

// 更新 JSON 显示
function updateJSON() {
    const data = {
        start_time: startTimeInput.value,
        end_time: endTimeInput.value,
        round: parseInt(roundSelect.value, 10),
        max_num: parseInt(maxNumInput.value, 10) || 10
    };
    jsonOutput.textContent = JSON.stringify(data, null, 2);
}

// 提交数据
async function handleSubmit() {
    const data = {
        start_time: startTimeInput.value,
        end_time: endTimeInput.value,
        round: parseInt(roundSelect.value, 10),
        max_num: parseInt(maxNumInput.value, 10)
    };

    // 简单校验
    if (!data.start_time || !data.end_time) {
        showMessage("请填写开始时间和结束时间", "error");
        return;
    }
    if (data.max_num < 1 || data.max_num > 50) {
        showMessage("最大人数需在 1~50 之间", "error");
        return;
    }

    try {
        const response = await fetch('/admin/set_schedule', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            showMessage("提交成功！", "success");
        } else {
            const err = await response.json().catch(() => ({}));
            showMessage(`提交失败: ${err.error || '未知错误'}`, "error");
        }
    } catch (err) {
        console.error(err);
        showMessage("网络错误，请确保后端正在运行", "error");
    }
}

// 清空表单
function handleClear() {
    startTimeInput.value = '';
    endTimeInput.value = '';
    roundSelect.value = '1';
    maxNumInput.value = '10';
    updateJSON();
    messageDiv.textContent = '';
}

// 显示消息
function showMessage(text, type) {
    messageDiv.textContent = text;
    messageDiv.style.color = type === 'success' ? 'green' : 'red';
}

// 事件绑定
startTimeInput.addEventListener('input', updateJSON);
endTimeInput.addEventListener('input', updateJSON);
roundSelect.addEventListener('change', updateJSON);
maxNumInput.addEventListener('input', updateJSON);

submitBtn.addEventListener('click', handleSubmit);
clearBtn.addEventListener('click', handleClear);

// 初始化 JSON 显示
updateJSON();