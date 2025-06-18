/**
 * 客户端图像处理器
 * 在浏览器中本地执行图像处理算法，无需服务器端计算
 */
class ClientImageProcessor {
    constructor() {
        this.canvas = null;
        this.ctx = null;
        this.originalImageData = null;
        this.processedImageData = null;
    }

    /**
     * 初始化处理器
     */
    init() {
        // 创建离屏canvas用于图像处理
        this.canvas = document.createElement('canvas');
        this.ctx = this.canvas.getContext('2d');
    }

    /**
     * 加载图像文件
     * @param {File} file - 图像文件
     * @returns {Promise<ImageData>} 图像数据
     */
    async loadImage(file) {
        return new Promise((resolve, reject) => {
            const img = new Image();
            img.onload = () => {
                // 设置canvas尺寸
                this.canvas.width = img.width;
                this.canvas.height = img.height;
                
                // 绘制图像到canvas
                this.ctx.drawImage(img, 0, 0);
                
                // 获取图像数据
                this.originalImageData = this.ctx.getImageData(0, 0, img.width, img.height);
                resolve(this.originalImageData);
            };
            img.onerror = reject;
            img.src = URL.createObjectURL(file);
        });
    }

    /**
     * 处理图像 - 客户端水印攻击算法
     * @param {number} attackLevel - 攻击强度 (0-1)
     * @returns {Promise<Blob>} 处理后的图像Blob
     */
    async processImage(attackLevel) {
        if (!this.originalImageData) {
            throw new Error('请先加载图像');
        }

        // 复制原始图像数据
        const imageData = new ImageData(
            new Uint8ClampedArray(this.originalImageData.data),
            this.originalImageData.width,
            this.originalImageData.height
        );

        // 执行多轮攻击算法
        await this.applyGeometricAttack(imageData, attackLevel);
        await this.applyNoiseAttack(imageData, attackLevel);
        await this.applyFrequencyAttack(imageData, attackLevel);
        await this.applyColorAttack(imageData, attackLevel);
        
        if (attackLevel > 0.7) {
            await this.applyFinalMixedAttack(imageData, attackLevel);
        }

        this.processedImageData = imageData;

        // 将处理后的图像数据转换为Blob
        return this.imageDataToBlob(imageData);
    }

    /**
     * 几何攻击 - 旋转和缩放
     */
    async applyGeometricAttack(imageData, level) {
        const { width, height } = imageData;
        
        // 创建临时canvas进行几何变换
        const tempCanvas = document.createElement('canvas');
        const tempCtx = tempCanvas.getContext('2d');
        tempCanvas.width = width;
        tempCanvas.height = height;
        
        // 将图像数据绘制到临时canvas
        tempCtx.putImageData(imageData, 0, 0);
        
        // 多轮几何变换
        const rounds = Math.floor(level * 3) + 1;
        for (let i = 0; i < rounds; i++) {
            // 随机旋转
            const angle = (Math.random() - 0.5) * level * 8 * Math.PI / 180;
            
            // 随机缩放
            const scale = 1.0 + (Math.random() - 0.5) * level * 0.1;
            
            // 应用变换
            tempCtx.save();
            tempCtx.translate(width / 2, height / 2);
            tempCtx.rotate(angle);
            tempCtx.scale(scale, scale);
            tempCtx.translate(-width / 2, -height / 2);
            
            // 重新绘制
            const currentImageData = tempCtx.getImageData(0, 0, width, height);
            tempCtx.clearRect(0, 0, width, height);
            tempCtx.putImageData(currentImageData, 0, 0);
            
            tempCtx.restore();
        }
        
        // 获取变换后的图像数据
        const transformedData = tempCtx.getImageData(0, 0, width, height);
        imageData.data.set(transformedData.data);
    }

    /**
     * 噪声攻击 - 亮度和对比度调整
     */
    async applyNoiseAttack(imageData, level) {
        const data = imageData.data;
        const rounds = Math.floor(level * 3) + 1;
        
        for (let round = 0; round < rounds; round++) {
            // 随机亮度和对比度调整
            const brightness = (Math.random() - 0.5) * level * 60;
            const contrast = 1.0 + (Math.random() - 0.5) * level * 0.8;
            
            for (let i = 0; i < data.length; i += 4) {
                // 应用对比度
                data[i] = Math.max(0, Math.min(255, (data[i] - 128) * contrast + 128));     // R
                data[i + 1] = Math.max(0, Math.min(255, (data[i + 1] - 128) * contrast + 128)); // G
                data[i + 2] = Math.max(0, Math.min(255, (data[i + 2] - 128) * contrast + 128)); // B
                
                // 应用亮度
                data[i] = Math.max(0, Math.min(255, data[i] + brightness));
                data[i + 1] = Math.max(0, Math.min(255, data[i + 1] + brightness));
                data[i + 2] = Math.max(0, Math.min(255, data[i + 2] + brightness));
            }
        }
    }

    /**
     * 频域攻击 - 模糊和锐化
     */
    async applyFrequencyAttack(imageData, level) {
        // 应用高斯模糊
        if (level > 0.2) {
            const blurRadius = Math.floor(level * 3) + 1;
            this.applyGaussianBlur(imageData, blurRadius);
        }
        
        // 应用锐化
        if (level > 0.3) {
            const sharpenStrength = level * 2;
            this.applySharpen(imageData, sharpenStrength);
        }
    }

    /**
     * 高斯模糊实现
     */
    applyGaussianBlur(imageData, radius) {
        const { width, height, data } = imageData;
        const output = new Uint8ClampedArray(data);
        
        // 简化的高斯模糊 - 使用盒式滤波器近似
        const boxSize = radius * 2 + 1;
        const weight = 1 / (boxSize * boxSize);
        
        for (let y = 0; y < height; y++) {
            for (let x = 0; x < width; x++) {
                let r = 0, g = 0, b = 0;
                
                for (let dy = -radius; dy <= radius; dy++) {
                    for (let dx = -radius; dx <= radius; dx++) {
                        const ny = Math.max(0, Math.min(height - 1, y + dy));
                        const nx = Math.max(0, Math.min(width - 1, x + dx));
                        const idx = (ny * width + nx) * 4;
                        
                        r += data[idx] * weight;
                        g += data[idx + 1] * weight;
                        b += data[idx + 2] * weight;
                    }
                }
                
                const idx = (y * width + x) * 4;
                output[idx] = Math.max(0, Math.min(255, r));
                output[idx + 1] = Math.max(0, Math.min(255, g));
                output[idx + 2] = Math.max(0, Math.min(255, b));
            }
        }
        
        data.set(output);
    }

    /**
     * 锐化滤波器
     */
    applySharpen(imageData, strength) {
        const { width, height, data } = imageData;
        const output = new Uint8ClampedArray(data);
        
        // 锐化卷积核
        const kernel = [
            0, -strength, 0,
            -strength, 1 + 4 * strength, -strength,
            0, -strength, 0
        ];
        
        for (let y = 1; y < height - 1; y++) {
            for (let x = 1; x < width - 1; x++) {
                let r = 0, g = 0, b = 0;
                
                for (let ky = -1; ky <= 1; ky++) {
                    for (let kx = -1; kx <= 1; kx++) {
                        const idx = ((y + ky) * width + (x + kx)) * 4;
                        const kernelIdx = (ky + 1) * 3 + (kx + 1);
                        
                        r += data[idx] * kernel[kernelIdx];
                        g += data[idx + 1] * kernel[kernelIdx];
                        b += data[idx + 2] * kernel[kernelIdx];
                    }
                }
                
                const idx = (y * width + x) * 4;
                output[idx] = Math.max(0, Math.min(255, r));
                output[idx + 1] = Math.max(0, Math.min(255, g));
                output[idx + 2] = Math.max(0, Math.min(255, b));
            }
        }
        
        data.set(output);
    }

    /**
     * 颜色攻击 - 色彩空间变换
     */
    async applyColorAttack(imageData, level) {
        const data = imageData.data;
        const rounds = Math.floor(level * 4) + 1;
        
        for (let round = 0; round < rounds; round++) {
            // 随机色彩调整
            const hueShift = (Math.random() - 0.5) * level * 60;
            const saturationScale = 1.0 + (Math.random() - 0.5) * level * 0.6;
            
            for (let i = 0; i < data.length; i += 4) {
                // RGB转HSV
                const r = data[i] / 255;
                const g = data[i + 1] / 255;
                const b = data[i + 2] / 255;
                
                const max = Math.max(r, g, b);
                const min = Math.min(r, g, b);
                const delta = max - min;
                
                let h = 0, s = 0, v = max;
                
                if (delta !== 0) {
                    s = delta / max;
                    
                    if (max === r) {
                        h = ((g - b) / delta) % 6;
                    } else if (max === g) {
                        h = (b - r) / delta + 2;
                    } else {
                        h = (r - g) / delta + 4;
                    }
                    h *= 60;
                }
                
                // 应用色彩变换
                h = (h + hueShift) % 360;
                if (h < 0) h += 360;
                s = Math.max(0, Math.min(1, s * saturationScale));
                
                // HSV转RGB
                const c = v * s;
                const x = c * (1 - Math.abs((h / 60) % 2 - 1));
                const m = v - c;
                
                let rNew, gNew, bNew;
                
                if (h < 60) {
                    [rNew, gNew, bNew] = [c, x, 0];
                } else if (h < 120) {
                    [rNew, gNew, bNew] = [x, c, 0];
                } else if (h < 180) {
                    [rNew, gNew, bNew] = [0, c, x];
                } else if (h < 240) {
                    [rNew, gNew, bNew] = [0, x, c];
                } else if (h < 300) {
                    [rNew, gNew, bNew] = [x, 0, c];
                } else {
                    [rNew, gNew, bNew] = [c, 0, x];
                }
                
                data[i] = Math.max(0, Math.min(255, (rNew + m) * 255));
                data[i + 1] = Math.max(0, Math.min(255, (gNew + m) * 255));
                data[i + 2] = Math.max(0, Math.min(255, (bNew + m) * 255));
            }
        }
    }

    /**
     * 最终混合攻击
     */
    async applyFinalMixedAttack(imageData, level) {
        // 组合多种攻击方式
        for (let i = 0; i < 2; i++) {
            await this.applyNoiseAttack(imageData, level * 0.3);
            this.applyGaussianBlur(imageData, Math.floor(level * 2));
            this.applySharpen(imageData, level * 1.5);
            await this.applyColorAttack(imageData, level * 0.4);
        }
    }

    /**
     * 将ImageData转换为Blob
     */
    async imageDataToBlob(imageData, format = 'image/jpeg', quality = 0.9) {
        return new Promise((resolve) => {
            // 创建临时canvas
            const tempCanvas = document.createElement('canvas');
            const tempCtx = tempCanvas.getContext('2d');
            tempCanvas.width = imageData.width;
            tempCanvas.height = imageData.height;
            
            // 绘制图像数据
            tempCtx.putImageData(imageData, 0, 0);
            
            // 转换为Blob
            tempCanvas.toBlob(resolve, format, quality);
        });
    }

    /**
     * 获取处理进度（模拟）
     */
    getProgress() {
        // 这里可以实现真实的进度跟踪
        return Math.random() * 100;
    }

    /**
     * 清理资源
     */
    cleanup() {
        this.canvas = null;
        this.ctx = null;
        this.originalImageData = null;
        this.processedImageData = null;
    }
}

// 导出处理器类
window.ClientImageProcessor = ClientImageProcessor;