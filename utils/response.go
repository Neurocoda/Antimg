package utils

import (
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/bmp"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

func SendImageResponse(c *gin.Context, format string, img image.Image) {
	// 设置下载头部，强制下载而不是在浏览器中显示
	filename := "processed_image." + format
	if format == "jpg" {
		filename = "processed_image.jpeg"
	}

	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Writer.Header().Set("Cache-Control", "no-cache")

	switch format {
	case "jpeg", "jpg":
		jpeg.Encode(c.Writer, img, &jpeg.Options{Quality: 90})
	case "png":
		png.Encode(c.Writer, img)
	case "bmp":
		// BMP输出
		bmp.Encode(c.Writer, img)
	case "webp":
		// WebP格式转换为JPEG输出（因为Go标准库不支持WebP编码）
		jpeg.Encode(c.Writer, img, &jpeg.Options{Quality: 90})
	default:
		// 默认使用JPEG格式输出
		jpeg.Encode(c.Writer, img, &jpeg.Options{Quality: 90})
	}
}
