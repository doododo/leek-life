package model

import (
	"time"
)

// Stock 股票模型
type Stock struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Price        float64   `json:"price"`
	TodayChange  float64   `json:"today_change"`
	AddDate      time.Time `json:"add_date"`
	AddPrice     float64   `json:"add_price"`
	TotalChange  float64   `json:"total_change"`
	OrderIndex   int       `json:"order_index"`
}