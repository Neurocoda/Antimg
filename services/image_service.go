package services

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"math/rand"
	"time"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
	_ "image/png"
)

type ImageService struct {
	rng *rand.Rand
}

func NewImageService() *ImageService {
	// Use new random number generator API
	return &ImageService{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ProcessImage 处理上传的图片，带超时控制
func (s *ImageService) ProcessImage(src io.Reader, attackLevel float64) (image.Image, string, error) {
	// 创建带超时的上下文 (30秒超时)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 使用通道来处理超时
	type result struct {
		img    image.Image
		format string
		err    error
	}

	resultChan := make(chan result, 1)

	go func() {
		// 解码图片，同时获取格式信息
		img, format, err := image.Decode(src)
		if err != nil {
			resultChan <- result{nil, "", err}
			return
		}

		// 执行水印攻击
		processedImg := s.attackWatermark(img, attackLevel)
		resultChan <- result{processedImg, format, nil}
	}()

	// 等待结果或超时
	select {
	case res := <-resultChan:
		return res.img, res.format, res.err
	case <-ctx.Done():
		return nil, "", errors.New("图片处理超时，请尝试较小的图片或降低攻击强度")
	}
}

// attackWatermark 执行水印攻击算法
func (s *ImageService) attackWatermark(img image.Image, attackLevel float64) image.Image {
	// 强化攻击算法 - 多轮攻击
	result := img

	// 第一轮：强力几何攻击
	result = s.applyAggressiveGeometricAttack(result, attackLevel)

	// 第二轮：强力噪声攻击
	result = s.applyAggressiveNoiseAttack(result, attackLevel)

	// 第三轮：强力频域攻击
	result = s.applyAggressiveFrequencyAttack(result, attackLevel)

	// 第四轮：强力压缩攻击
	result = s.applyAggressiveCompressionAttack(result, attackLevel)

	// 第五轮：强力颜色攻击
	result = s.applyAggressiveColorAttack(result, attackLevel)

	// 最终轮：混合攻击
	if attackLevel > 0.7 {
		result = s.applyFinalMixedAttack(result, attackLevel)
	}

	return result
}

// applyAggressiveNoiseAttack 强力噪声攻击
func (s *ImageService) applyAggressiveNoiseAttack(img image.Image, level float64) image.Image {
	result := imaging.Clone(img)

	// 强力亮度攻击
	brightnessChange := (s.rng.Float64() - 0.5) * level * 60 // 最大±30亮度变化
	result = imaging.AdjustBrightness(result, brightnessChange)

	// 强力对比度攻击
	contrastChange := (s.rng.Float64() - 0.5) * level * 80 // 最大±40对比度变化
	result = imaging.AdjustContrast(result, contrastChange)

	// 多次随机调整
	rounds := int(level*3) + 1
	for i := 0; i < rounds; i++ {
		brightness := (s.rng.Float64() - 0.5) * level * 20
		contrast := (s.rng.Float64() - 0.5) * level * 30
		result = imaging.AdjustBrightness(result, brightness)
		result = imaging.AdjustContrast(result, contrast)
	}

	return result
}

// applyAggressiveFrequencyAttack 强力频域攻击
func (s *ImageService) applyAggressiveFrequencyAttack(img image.Image, level float64) image.Image {
	result := img

	// 强力模糊攻击
	blurRadius := level * 3.0 // 大幅增加模糊半径
	if blurRadius > 0.5 {
		result = imaging.Blur(result, blurRadius)
	}

	// 强力锐化攻击
	if level > 0.3 {
		sharpenAmount := level * 5.0 // 大幅增加锐化强度
		result = imaging.Sharpen(result, sharpenAmount)
	}

	// 交替模糊和锐化
	rounds := int(level*2) + 1
	for i := 0; i < rounds; i++ {
		if i%2 == 0 {
			result = imaging.Blur(result, level*2.0)
		} else {
			result = imaging.Sharpen(result, level*3.0)
		}
	}

	// 最终强力模糊
	if level > 0.7 {
		result = imaging.Blur(result, level*4.0)
	}

	return result
}

// applyAggressiveGeometricAttack 强力几何攻击
func (s *ImageService) applyAggressiveGeometricAttack(img image.Image, level float64) image.Image {
	bounds := img.Bounds()
	result := img

	// 强力旋转攻击
	if level > 0.2 {
		angle := (s.rng.Float64() - 0.5) * level * 15 // 大幅增加旋转角度
		result = imaging.Rotate(result, angle, color.Transparent)
	}

	// 强力缩放攻击
	if level > 0.3 {
		scaleFactor := 1.0 + (s.rng.Float64()-0.5)*level*0.2 // 大幅增加缩放范围
		newWidth := int(float64(bounds.Dx()) * scaleFactor)
		newHeight := int(float64(bounds.Dy()) * scaleFactor)
		result = imaging.Resize(result, newWidth, newHeight, imaging.Lanczos)
		// 裁剪回原始大小
		result = imaging.CropCenter(result, bounds.Dx(), bounds.Dy())
	}

	// 多轮几何变换
	rounds := int(level*3) + 1
	for i := 0; i < rounds; i++ {
		// 随机旋转
		angle := (s.rng.Float64() - 0.5) * level * 8
		result = imaging.Rotate(result, angle, color.Transparent)

		// 随机缩放
		scale := 1.0 + (s.rng.Float64()-0.5)*level*0.1
		newW := int(float64(bounds.Dx()) * scale)
		newH := int(float64(bounds.Dy()) * scale)
		result = imaging.Resize(result, newW, newH, imaging.Lanczos)
		result = imaging.CropCenter(result, bounds.Dx(), bounds.Dy())
	}

	// 最终强力变换
	if level > 0.8 {
		finalAngle := (s.rng.Float64() - 0.5) * level * 20
		result = imaging.Rotate(result, finalAngle, color.Transparent)
	}

	return result
}

// applyAggressiveCompressionAttack 强力压缩攻击
func (s *ImageService) applyAggressiveCompressionAttack(img image.Image, level float64) image.Image {
	// 根据攻击等级调整JPEG质量 - 更激进
	quality := 100 - int(level*60) // 质量从100降到40
	if quality < 30 {
		quality = 30
	}

	// 多轮压缩攻击
	result := img
	compressionRounds := int(level*5) + 1 // 最多6轮压缩

	for i := 0; i < compressionRounds; i++ {
		// 每轮都降低质量
		currentQuality := quality - i*5
		if currentQuality < 20 {
			currentQuality = 20
		}

		var buf bytes.Buffer
		jpeg.Encode(&buf, result, &jpeg.Options{Quality: currentQuality})

		// 解码回图片
		decodedImg, err := jpeg.Decode(&buf)
		if err != nil {
			return img // 如果失败，返回原图
		}
		result = decodedImg
	}

	return result
}

// applyAggressiveColorAttack 强力颜色攻击
func (s *ImageService) applyAggressiveColorAttack(img image.Image, level float64) image.Image {
	result := img

	// 强力亮度和对比度攻击
	rounds := int(level*4) + 1
	for i := 0; i < rounds; i++ {
		brightnessChange := (s.rng.Float64() - 0.5) * level * 50 // 大幅增加亮度变化
		contrastChange := (s.rng.Float64() - 0.5) * level * 60   // 大幅增加对比度变化
		result = imaging.AdjustBrightness(result, brightnessChange)
		result = imaging.AdjustContrast(result, contrastChange)
	}

	return result
}

// applyFinalMixedAttack 最终混合攻击
func (s *ImageService) applyFinalMixedAttack(img image.Image, level float64) image.Image {
	result := img

	// 最终破坏性攻击组合
	for i := 0; i < 3; i++ {
		// 强力模糊
		result = imaging.Blur(result, level*5.0)

		// 强力锐化
		result = imaging.Sharpen(result, level*6.0)

		// 强力亮度对比度调整
		brightness := (s.rng.Float64() - 0.5) * level * 40
		contrast := (s.rng.Float64() - 0.5) * level * 50
		result = imaging.AdjustBrightness(result, brightness)
		result = imaging.AdjustContrast(result, contrast)

		// 旋转攻击
		angle := (s.rng.Float64() - 0.5) * level * 10
		result = imaging.Rotate(result, angle, color.Transparent)

		// 压缩攻击
		quality := 50 - int(level*30)
		if quality < 15 {
			quality = 15
		}
		var buf bytes.Buffer
		jpeg.Encode(&buf, result, &jpeg.Options{Quality: quality})
		decodedImg, err := jpeg.Decode(&buf)
		if err == nil {
			result = decodedImg
		}
	}

	return result
}
