document.addEventListener('DOMContentLoaded', () => {
    const toast = document.getElementById('toast');
    const downloadBtn = document.getElementById('downloadBtn');

    // 显示 Toast 提示
    function showToast(message, type = 'success', duration = 6000) {
        toast.textContent = message;
        toast.className = `toast ${type} show`;

        if (window.toastTimer) clearTimeout(window.toastTimer);
        window.toastTimer = setTimeout(() => {
            toast.classList.remove('show');
        }, duration);
    }

    // 下载按钮点击事件
    downloadBtn?.addEventListener('click', async function () {
        const btn = this;
        btn.disabled = true;
        btn.textContent = '检测中...';

        try {
            // 先用 HEAD 请求检查文件是否存在
            const checkRes = await fetch('/user/download/file', { method: 'HEAD' });

            if (checkRes.ok) {
                // 触发下载（浏览器会处理 Content-Disposition）
                window.location.href = '/user/download/file';
                showToast('下载已开始！请查看浏览器底部下载栏', 'success', 6000);
            } else {
                showToast('暂无面试题文件，请联系管理员上传', 'error', 6000);
            }
        } catch (err) {
            console.error('检测或下载失败:', err);
            showToast('网络错误，请稍后重试', 'error', 6000);
        } finally {
            // 恢复按钮状态（延迟 1 秒避免闪现）
            setTimeout(() => {
                if (downloadBtn) {
                    downloadBtn.disabled = false;
                    downloadBtn.textContent = '立即下载';
                }
            }, 1000);
        }
    });
});