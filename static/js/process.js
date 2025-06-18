// 全局变量
let currentImageFile = null;
let currentImageData = null;

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
    initializeUpload();
    initializeAttackLevelSlider();
});

// 初始化上传功能
function initializeUpload() {
    const imageFile = document.getElementById('imageFile');
    const uploadArea = document.getElementById('uploadArea');
    
    // 文件选择事件
    imageFile.addEventListener('change', handleFileSelect);
    
    // 拖拽功能
    uploadArea.addEventListener('dragover', function(e) {
        e.preventDefault();
        uploadArea.style.backgroundColor = 'rgba(0,113,227,0.1)';
    });
    
    uploadArea.addEventListener('dragleave', function(e) {
        e.preventDefault();
        uploadArea.style.backgroundColor = '';
    });
    
    uploadArea.addEventListener('drop', function(e) {
        e.preventDefault();
        uploadArea.style.backgroundColor = '';
        
        const files = e.dataTransfer.files;
        if (files.length > 0) {
            imageFile.files = files;
            handleFileSelect();
        }
    });
}

// 处理文件选择
function handleFileSelect() {
    const imageFile = document.getElementById('imageFile');
    const file = imageFile.files[0];
    
    if (!file) return;
    
    // 验证文件类型
    if (!file.type.startsWith('image/')) {
        alert('请选择图片文件！');
        return;
    }
    
    currentImageFile = file;
    
    // 显示图片预览
    const reader = new FileReader();
    reader.onload = function(e) {
        currentImageData = e.target.result;
        showImagePreview(file, e.target.result);
        switchToProcessSection();
    };
    reader.readAsDataURL(file);
}

// 显示图片预览
function showImagePreview(file, dataUrl) {
    const previewImg = document.getElementById('previewImg');
    const fileName = document.getElementById('fileName');
    const imageDimensions = document.getElementById('imageDimensions');
    
    previewImg.src = dataUrl;
    fileName.textContent = file.name;
    
    // 获取图片尺寸
    previewImg.onload = function() {
        imageDimensions.textContent = `${this.naturalWidth} × ${this.naturalHeight}`;
    };
}

// 切换到处理界面
function switchToProcessSection() {
    document.getElementById('uploadSection').style.display = 'none';
    document.getElementById('processSection').style.display = 'block';
}

// 重置上传
function resetUpload() {
    document.getElementById('uploadSection').style.display = 'block';
    document.getElementById('processSection').style.display = 'none';
    document.getElementById('imageFile').value = '';
    currentImageFile = null;
    currentImageData = null;
    
    // 重置处理动画状态
    resetProcessingAnimation();
    
    // 清理处理后的图片blob
    if (window.currentProcessedBlob) {
        window.currentProcessedBlob = null;
    }
    
    // 隐藏攻击强度信息
    const attackLevelInfo = document.getElementById('attackLevelInfo');
    if (attackLevelInfo) {
        attackLevelInfo.style.display = 'none';
    }
}

// 初始化攻击强度滑块
function initializeAttackLevelSlider() {
    const attackLevel = document.getElementById('attackLevel');
    const strengthValue = document.getElementById('strengthValue');
    
    if (attackLevel && strengthValue) {
        // 更新数值显示
        attackLevel.addEventListener('input', function() {
            const value = parseFloat(this.value);
            strengthValue.textContent = value.toFixed(2);
            updateValueColor(strengthValue, value);
        });
        
        // 初始化显示
        updateValueColor(strengthValue, 0.5);
    }
}

// 更新数值颜色
function updateValueColor(element, value) {
    let color;
    if (value <= 0.3) {
        color = '#10b981'; // 绿色
    } else if (value <= 0.7) {
        color = '#f59e0b'; // 橙色
    } else {
        color = '#ef4444'; // 红色
    }
    
    element.style.background = color;
    element.style.color = 'white';
}

// 处理图片 - 带实时预览功能
function processImageWithPreview() {
    if (!currentImageFile) {
        alert('请先选择图片！');
        return;
    }
    
    const attackLevel = document.getElementById('attackLevel').value;
    const processBtn = document.getElementById('processBtn');
    const downloadBtn = document.getElementById('downloadBtn');
    
    // 更新攻击强度显示
    const currentAttackValue = document.getElementById('currentAttackValue');
    const attackLevelInfo = document.getElementById('attackLevelInfo');
    if (currentAttackValue && attackLevelInfo) {
        currentAttackValue.textContent = parseFloat(attackLevel).toFixed(2);
        attackLevelInfo.style.display = 'block';
    }
    
    // 禁用处理按钮
    processBtn.disabled = true;
    processBtn.innerHTML = '<span data-i18n="processing">处理中...</span>';
    
    // 隐藏下载按钮
    downloadBtn.style.display = 'none';
    
    // 开始动画：原图左移，显示处理区域
    startProcessingAnimation();
    
    // 创建FormData
    const formData = new FormData();
    formData.append('image', currentImageFile);
    formData.append('attackLevel', attackLevel);
    
    // 发送请求
    fetch('/admin/process', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('处理失败');
        }
        return response.blob();
    })
    .then(blob => {
        // 创建处理后图片的URL
        const processedImageUrl = window.URL.createObjectURL(blob);
        
        // 显示处理结果
        showProcessedResult(processedImageUrl, blob);
        
        // 恢复按钮状态
        processBtn.disabled = false;
        processBtn.innerHTML = '<span data-i18n="processImage">处理图片</span>';
        
        // 显示下载按钮
        downloadBtn.style.display = 'inline-block';
    })
    .catch(error => {
        console.error('Error:', error);
        alert('图片处理失败，请重试！');
        
        // 恢复按钮状态和布局
        processBtn.disabled = false;
        processBtn.innerHTML = '<span data-i18n="processImage">处理图片</span>';
        resetProcessingAnimation();
    });
}

// 开始处理动画
function startProcessingAnimation() {
    const originalContainer = document.getElementById('originalImageContainer');
    const processedContainer = document.getElementById('processedImageContainer');
    const loadingAnimation = document.getElementById('loadingAnimation');
    const processedImg = document.getElementById('processedImg');
    const previewImg = document.getElementById('previewImg');
    
    // 等待图片完全加载后获取准确尺寸
    const getImageDimensions = () => {
        return new Promise((resolve) => {
            if (previewImg.complete && previewImg.naturalWidth > 0) {
                resolve();
            } else {
                previewImg.onload = () => resolve();
            }
        });
    };
    
    getImageDimensions().then(() => {
        // 获取原图的实际显示尺寸（考虑max-width和max-height约束）
        const imgRect = previewImg.getBoundingClientRect();
        const originalDisplayWidth = Math.round(imgRect.width);
        const originalDisplayHeight = Math.round(imgRect.height);
        
        console.log(`原图显示尺寸: ${originalDisplayWidth}x${originalDisplayHeight}`);
        console.log(`原图自然尺寸: ${previewImg.naturalWidth}x${previewImg.naturalHeight}`);
        
        // 原图左移到48%宽度，保持居中对齐
        originalContainer.style.width = '48%';
        originalContainer.style.display = 'inline-block';
        originalContainer.style.marginRight = '4%';
        originalContainer.style.textAlign = 'center';
        
        // 显示处理区域，保持居中对齐
        processedContainer.style.display = 'inline-block';
        processedContainer.style.position = 'relative';
        processedContainer.style.width = '48%';
        processedContainer.style.right = 'auto';
        processedContainer.style.top = 'auto';
        processedContainer.style.opacity = '1';
        processedContainer.style.transform = 'translateX(0)';
        processedContainer.style.textAlign = 'center';
        
        // 精确设置加载动画区域的尺寸与原图显示尺寸一致
        loadingAnimation.style.width = originalDisplayWidth + 'px';
        loadingAnimation.style.height = originalDisplayHeight + 'px';
        loadingAnimation.style.minWidth = originalDisplayWidth + 'px';
        loadingAnimation.style.minHeight = originalDisplayHeight + 'px';
        loadingAnimation.style.maxWidth = originalDisplayWidth + 'px';
        loadingAnimation.style.maxHeight = originalDisplayHeight + 'px';
        loadingAnimation.style.display = 'flex';
        loadingAnimation.style.margin = '0 auto'; // 确保加载动画在容器中居中
        loadingAnimation.classList.add('loading-size-matched');
        
        // 隐藏处理后的图片
        processedImg.style.display = 'none';
        
        // 启动进度条动画
        startProgressAnimation();
    });
}

// 启动进度条动画
function startProgressAnimation() {
    const progressBar = document.getElementById('progressBar');
    let progress = 0;
    
    const progressInterval = setInterval(() => {
        progress += Math.random() * 15 + 5; // 随机增长5-20%
        if (progress > 95) {
            progress = 95; // 最多到95%，等待实际完成
        }
        progressBar.style.width = progress + '%';
        
        if (progress >= 95) {
            clearInterval(progressInterval);
        }
    }, 200);
    
    // 保存interval ID以便后续清理
    window.currentProgressInterval = progressInterval;
}

// 显示处理结果
function showProcessedResult(imageUrl, blob) {
    const loadingAnimation = document.getElementById('loadingAnimation');
    const processedImg = document.getElementById('processedImg');
    const progressBar = document.getElementById('progressBar');
    const previewImg = document.getElementById('previewImg');
    
    // 清理进度条动画
    if (window.currentProgressInterval) {
        clearInterval(window.currentProgressInterval);
    }
    
    // 完成进度条
    progressBar.style.width = '100%';
    
    // 设置处理后的图片
    processedImg.src = imageUrl;
    processedImg.onload = function() {
        // 获取原图的实际显示尺寸（与加载动画时保持一致）
        const imgRect = previewImg.getBoundingClientRect();
        const originalDisplayWidth = Math.round(imgRect.width);
        const originalDisplayHeight = Math.round(imgRect.height);
        
        console.log(`设置处理后图片尺寸: ${originalDisplayWidth}x${originalDisplayHeight}`);
        
        // 精确设置处理后图片的尺寸与原图显示尺寸完全一致
        processedImg.style.width = originalDisplayWidth + 'px';
        processedImg.style.height = originalDisplayHeight + 'px';
        processedImg.style.maxWidth = 'none'; // 重要：取消max-width限制
        processedImg.style.maxHeight = 'none'; // 重要：取消max-height限制
        processedImg.style.minWidth = originalDisplayWidth + 'px';
        processedImg.style.minHeight = originalDisplayHeight + 'px';
        processedImg.style.objectFit = 'contain'; // 保持图片比例，填充容器
        processedImg.style.margin = '0 auto'; // 确保图片在容器中居中
        
        // 隐藏加载动画，显示处理后的图片
        setTimeout(() => {
            loadingAnimation.style.display = 'none';
            processedImg.style.display = 'block';
            
            // 添加淡入效果
            processedImg.style.opacity = '0';
            processedImg.style.transition = 'opacity 0.5s ease';
            setTimeout(() => {
                processedImg.style.opacity = '1';
            }, 50);
        }, 500); // 稍微延迟以显示100%进度
    };
    
    // 保存blob用于下载
    window.currentProcessedBlob = blob;
}

// 重置处理动画
function resetProcessingAnimation() {
    const originalContainer = document.getElementById('originalImageContainer');
    const processedContainer = document.getElementById('processedImageContainer');
    const loadingAnimation = document.getElementById('loadingAnimation');
    const processedImg = document.getElementById('processedImg');
    const progressBar = document.getElementById('progressBar');
    
    // 清理进度条动画
    if (window.currentProgressInterval) {
        clearInterval(window.currentProgressInterval);
    }
    
    // 重置原图容器
    originalContainer.style.width = '100%';
    originalContainer.style.marginRight = '0';
    
    // 隐藏处理区域
    processedContainer.style.display = 'none';
    processedContainer.style.position = '';
    processedContainer.style.width = '';
    processedContainer.style.right = '';
    processedContainer.style.top = '';
    processedContainer.style.opacity = '0';
    processedContainer.style.transform = 'translateX(20px)';
    
    // 完全重置加载动画区域尺寸
    loadingAnimation.style.width = '';
    loadingAnimation.style.height = '';
    loadingAnimation.style.minWidth = '';
    loadingAnimation.style.minHeight = '';
    loadingAnimation.style.maxWidth = '';
    loadingAnimation.style.maxHeight = '';
    loadingAnimation.style.margin = '';
    loadingAnimation.style.display = 'flex';
    loadingAnimation.classList.remove('loading-size-matched');
    
    // 完全重置处理后图片尺寸和样式
    processedImg.style.width = '';
    processedImg.style.height = '';
    processedImg.style.minWidth = '';
    processedImg.style.minHeight = '';
    processedImg.style.maxWidth = '';
    processedImg.style.maxHeight = '';
    processedImg.style.objectFit = '';
    processedImg.style.margin = '';
    processedImg.style.display = 'none';
    processedImg.style.opacity = '';
    processedImg.style.transition = '';
    
    // 重置进度条
    progressBar.style.width = '0%';
}

// 下载处理后的图片
function downloadProcessedImage() {
    if (!window.currentProcessedBlob) {
        alert('没有可下载的处理后图片！');
        return;
    }
    
    const url = window.URL.createObjectURL(window.currentProcessedBlob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `processed_${currentImageFile.name}`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
    
    showNotification('图片下载成功！', 'success');
}

// 兼容旧的processImage函数调用
function processImage() {
    processImageWithPreview();
}

// 复制Token到剪贴板
function copyToken() {
    const tokenElement = document.getElementById('userToken');
    const token = tokenElement.textContent.trim();
    
    if (navigator.clipboard && window.isSecureContext) {
        navigator.clipboard.writeText(token).then(() => {
            showCopySuccess('Token 已复制到剪贴板！');
        }).catch(err => {
            fallbackCopyTextToClipboard(token);
        });
    } else {
        fallbackCopyTextToClipboard(token);
    }
}

// 复制cURL命令到剪贴板
function copyCurlCommand() {
    const tokenElement = document.getElementById('userToken');
    const token = tokenElement.textContent.trim();
    const baseURL = window.location.origin;
    
    const curlCommand = `curl -X POST "${baseURL}/api/attack" \\
  -H "Authorization: Bearer ${token}" \\
  -F "image=@your_image.jpg" \\
  -F "attackLevel=0.65" \\
  --output processed_image.jpg`;
    
    if (navigator.clipboard && window.isSecureContext) {
        navigator.clipboard.writeText(curlCommand).then(() => {
            showCopySuccess('cURL 命令已复制到剪贴板！');
        }).catch(err => {
            fallbackCopyTextToClipboard(curlCommand);
        });
    } else {
        fallbackCopyTextToClipboard(curlCommand);
    }
}

// 备用复制方法（兼容旧浏览器）
function fallbackCopyTextToClipboard(text) {
    const textArea = document.createElement("textarea");
    textArea.value = text;
    textArea.style.top = "0";
    textArea.style.left = "0";
    textArea.style.position = "fixed";
    
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    
    try {
        const successful = document.execCommand('copy');
        if (successful) {
            showCopySuccess('内容已复制到剪贴板！');
        } else {
            showCopyError('复制失败，请手动复制');
        }
    } catch (err) {
        showCopyError('复制失败，请手动复制');
    }
    
    document.body.removeChild(textArea);
}

// 显示复制成功提示
function showCopySuccess(message) {
    showNotification(message, 'success');
}

// 显示复制错误提示
function showCopyError(message) {
    showNotification(message, 'error');
}

// 显示通知
function showNotification(message, type = 'success') {
    // 移除已存在的通知
    const existingNotification = document.querySelector('.copy-notification');
    if (existingNotification) {
        existingNotification.remove();
    }
    
    const notification = document.createElement('div');
    notification.className = 'copy-notification';
    notification.textContent = message;
    
    // 设置样式
    const bgColor = type === 'success' ? 'var(--success)' : 'var(--error)';
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        background: ${bgColor};
        color: white;
        padding: 12px 20px;
        border-radius: 8px;
        box-shadow: var(--shadow);
        z-index: 10000;
        font-size: 0.9rem;
        font-weight: 500;
        opacity: 0;
        transform: translateX(100%);
        transition: all 0.3s ease;
    `;
    
    document.body.appendChild(notification);
    
    // 显示动画
    setTimeout(() => {
        notification.style.opacity = '1';
        notification.style.transform = 'translateX(0)';
    }, 10);
    
    // 自动隐藏
    setTimeout(() => {
        notification.style.opacity = '0';
        notification.style.transform = 'translateX(100%)';
        setTimeout(() => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 300);
    }, 3000);
}

// 重置API Token
function resetAPIToken() {
    if (!confirm('确定要重置 API Token 吗？重置后旧的 Token 将失效。')) {
        return;
    }
    
    fetch('/admin/reset-api-token', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.code === 200) {
            // 更新页面上的Token显示
            const tokenElement = document.getElementById('userToken');
            tokenElement.textContent = data.data.api_token;
            
            // 更新示例代码中的Token
            updateExampleCode(data.data.api_token);
            
            showCopySuccess('API Token 已重置！');
        } else {
            showCopyError('重置失败：' + (data.message || '未知错误'));
        }
    })
    .catch(error => {
        console.error('Error:', error);
        showCopyError('重置失败，请重试');
    });
}

// 更新示例代码中的Token
function updateExampleCode(newToken) {
    // 获取当前的baseURL
    const baseURL = window.location.origin;
    
    // 更新cURL示例代码
    const codeElement = document.querySelector('pre code');
    if (codeElement) {
        codeElement.textContent = `curl -X POST "${baseURL}/api/attack" \\
  -H "Authorization: Bearer ${newToken}" \\
  -F "image=@your_image.jpg" \\
  -F "attackLevel=0.65" \\
  --output processed_image.jpg`;
    }
}
