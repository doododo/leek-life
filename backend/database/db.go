package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"stock-app/model"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() {
	var err error
	// 连接SQLite数据库
	DB, err = gorm.Open(sqlite.Open("stock.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	// 自动迁移创建表
	err = DB.AutoMigrate(&model.Stock{})
	if err != nil {
		log.Fatalf("无法创建表: %v", err)
	}
}