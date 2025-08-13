package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"stock-app/controller"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8099", "http://192.168.65.1:8099"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// API路由组
	api := r.Group("/api")
	{
		// 股票相关路由
		api.GET("/stocks", controller.GetStocks)
		api.POST("/stocks", controller.AddStock)
		api.DELETE("/stocks/:id", controller.DeleteStock)
		
		// 股票排序路由
		api.PUT("/stocks/:id/move-up", controller.MoveStockUp)
		api.PUT("/stocks/:id/move-down", controller.MoveStockDown)
		api.PUT("/stocks/:id/move-top", controller.MoveStockToTop)
	}

	return r
}