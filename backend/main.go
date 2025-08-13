package main

import (
	"log"
	"stock-app/database"
	"stock-app/router"
	"stock-app/service"
)

func main() {
	// 初始化数据库
	database.InitDB()

	// 初始化已有股票的排序索引
	if err := service.InitializeOrderIndices(); err != nil {
		log.Printf("初始化排序索引失败: %v", err)
	}

	// 启动股票价格定时更新
	service.StartPriceUpdater()

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	log.Println("服务器启动在 http://localhost:8080")
	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}