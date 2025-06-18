package main

import (
	"log"

	"github.com/Neurocoda/Antimg/config"
	"github.com/Neurocoda/Antimg/routes"
)

// 构建时注入的版本信息
var (
	Version   = "dev"
	BuildTime = "unknown"
	Revision  = "unknown"
)

func main() {
	// 打印版本信息
	log.Printf("🚀 Antimg v%s (built at %s, revision %s)", Version, BuildTime, Revision)

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
