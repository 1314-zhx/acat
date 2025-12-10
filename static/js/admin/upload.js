const fileInput = document.getElementById('fileInput');
const uploadBtn = document.getElementById('uploadBtn');
const fileInfo = document.getElementById('fileInfo');
const messageEl = document.getElementById('message');

let selectedFile = null;

// 监听文件选择
fileInput.addEventListener('change', (e) => {
    const file = e.target.files[0];
    selectedFile = file;

    if (file) {
        // 检查文件类型
        const allowedExts = ['.pdf', '.doc', '.docx', '.txt', '.md'];
        const ext = '.' + file.name.split('.').pop().toLowerCase();
        if (!allowedExts.includes(ext)) {
            showMessage(`不支持的文件类型：${ext}`, 'error');
            uploadBtn.disabled = true;
            selectedFile = null;
            return;
        }

        // 显示文件信息
        const sizeMB = (file.size / (1024 * 1024)).toFixed(2);
        fileInfo.innerHTML = `
      <strong>文件名：</strong>${file.name}<br>
      <strong>大小：</strong>${sizeMB} MB
    `;
        fileInfo.classList.add('show');
        uploadBtn.disabled = false;
        hideMessage();
    } else {
        fileInfo.classList.remove('show');
        uploadBtn.disabled = true;
    }
});

// 上传按钮点击
uploadBtn.addEventListener('click', () => {
    if (!selectedFile) return;

    const formData = new FormData();
    formData.append('file', selectedFile);

    showMessage('上传中...', 'loading');
    uploadBtn.disabled = true;

    fetch('/admin/upload', {
        method: 'POST',
        body: formData,
        credentials: 'same-origin' // 保留认证 cookie
    })
        .then(response => response.json())
        .then(data => {
            if (data.status === 200) {
                showMessage('上传成功！', 'success');
                resetForm();
            } else {
                throw new Error(data.error || data.msg || '上传失败');
            }
        })
        .catch(err => {
            console.error('Upload error:', err);
            showMessage(err.message, 'error');
            uploadBtn.disabled = false;
        });
});

function showMessage(text, type) {
    messageEl.textContent = text;
    messageEl.className = 'message ' + type;
    messageEl.style.display = 'block';
}

function hideMessage() {
    messageEl.style.display = 'none';
}

function resetForm() {
    fileInput.value = '';
    fileInfo.classList.remove('show');
    uploadBtn.disabled = true;
    selectedFile = null;
}