package main

import (
	"log"

	"Antimg/config"
	"Antimg/routes"
)

func main() {
	// 初始化配置
	config.Init()

	// 设置路由
	r := routes.SetupRoutes()

	// 启动服务器
	port := ":" + config.AppConfig.Port
	log.Printf("🚀 Antimg 服务器启动成功")
	log.Printf("📡 端口: %s", port)
	log.Printf("👤 管理员账号: %s / %s", config.AppConfig.AdminUsername, config.AppConfig.AdminPassword)
	log.Printf("🌐 Web界面: http://localhost%s", port)
	log.Printf("📚 API接口: http://localhost%s/api", port)
	log.Printf("💡 使用 Ctrl+C 停止服务器")

	if err := r.Run(port); err != nil {
		log.Fatal("❌ 服务器启动失败:", err)
	}
}
