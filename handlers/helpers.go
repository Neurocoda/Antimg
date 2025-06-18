package handlers

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// parseAttackLevel 解析攻击强度参数
func parseAttackLevel(c *gin.Context) (float64, error) {
	attackLevelStr := c.PostForm("attackLevel")
	if attackLevelStr == "" {
		return 0.5, nil // 默认值
	}

	level, err := strconv.ParseFloat(attackLevelStr, 64)
	if err != nil {
		return 0, errors.New("攻击强度必须是数字")
	}

	if level < 0 || level > 1 {
		return 0, errors.New("攻击强度必须在0.0-1.0之间")
	}

	return level, nil
}

// validateImageFile 验证上传的图片文件
func validateImageFile(header *multipart.FileHeader) error {
	// 检查文件大小 (最大100MB)
	const maxFileSize = 100 << 20 // 100MB
	if header.Size > maxFileSize {
		return errors.New("文件大小超过限制，最大支持100MB")
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".bmp":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return errors.New("不支持的图片格式，仅支持: jpg, jpeg, png, bmp, webp")
	}

	// 检查MIME类型
	contentType := header.Header.Get("Content-Type")
	allowedMimes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/bmp":  true,
		"image/webp": true,
	}

	if contentType != "" && !allowedMimes[contentType] {
		return errors.New("无效的图片MIME类型")
	}

	return nil
}
