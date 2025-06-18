package handlers

import (
	"net/http"
	"strconv"

	"github.com/Neurocoda/Antimg/models"
	"github.com/Neurocoda/Antimg/services"
	"github.com/Neurocoda/Antimg/utils"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	imageService *services.ImageService
}

func NewImageHandler() *ImageHandler {
	return &ImageHandler{
		imageService: services.NewImageService(),
	}
}

// API: 攻击水印
func (h *ImageHandler) AttackWatermark(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "文件上传失败")
		return
	}

	// 验证文件类型
	if err := validateImageFile(file); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	src, err := file.Open()
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "文件打开错误")
		return
	}
	defer src.Close()

	attackLevel, err := parseAttackLevel(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "攻击强度参数无效: "+err.Error())
		return
	}

	processedImg, format, err := h.imageService.ProcessImage(src, attackLevel)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "图片处理失败: "+err.Error())
		return
	}

	utils.SendImageResponse(c, format, processedImg)
}

// Web: 图像处理页面
func (h *ImageHandler) ProcessPage(c *gin.Context) {
	// 获取用户信息
	username := c.GetString("username")
	user, err := models.GetUserByUsername(username)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	
	// 构建基础URL
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := scheme + "://" + c.Request.Host
	
	c.HTML(http.StatusOK, "base.html", gin.H{
		"title":    "图像处理工作台 - Antimg",
		"username": username,
		"token":    user.APIToken, // 使用API Token而不是Web Token
		"baseURL":  baseURL,
		"page":     "process",
	})
}

// Web: 处理图片上传
func (h *ImageHandler) WebProcessImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.HTML(http.StatusBadRequest, "base.html", gin.H{
			"title":    "图像处理工作台 - Antimg",
			"username": c.GetString("username"),
			"error":    "文件上传失败",
			"page":     "process",
		})
		return
	}

	// 验证文件类型
	if err := validateImageFile(file); err != nil {
		c.HTML(http.StatusBadRequest, "base.html", gin.H{
			"title":    "图像处理工作台 - Antimg",
			"username": c.GetString("username"),
			"error":    err.Error(),
			"page":     "process",
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.HTML(http.StatusBadRequest, "base.html", gin.H{
			"title":    "图像处理工作台 - Antimg",
			"username": c.GetString("username"),
			"error":    "文件打开错误",
			"page":     "process",
		})
		return
	}
	defer src.Close()

	attackLevelStr := c.PostForm("attackLevel")
	attackLevel := 0.5 // 默认值
	if attackLevelStr != "" {
		if level, err := strconv.ParseFloat(attackLevelStr, 64); err == nil && level >= 0 && level <= 1 {
			attackLevel = level
		}
	}

	processedImg, format, err := h.imageService.ProcessImage(src, attackLevel)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "base.html", gin.H{
			"title":    "图像处理工作台 - Antimg",
			"username": c.GetString("username"),
			"error":    "图片处理失败: " + err.Error(),
			"page":     "process",
		})
		return
	}

	// 直接返回处理后的图片，保持原格式
	utils.SendImageResponse(c, format, processedImg)
}
