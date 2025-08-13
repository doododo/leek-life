package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"stock-app/service"
)

// GetStocks 获取所有股票
func GetStocks(c *gin.Context) {
	stocks, err := service.GetAllStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stocks)
}

// AddStock 添加股票
func AddStock(c *gin.Context) {
	var request struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}

	if request.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "股票代码不能为空"})
		return
	}

	err := service.AddStock(request.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "添加成功"})
}

// DeleteStock 删除股票
func DeleteStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的股票ID"})
		return
	}

	err = service.DeleteStock(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消关注成功"})
}

// MoveStockUp 上移股票
func MoveStockUp(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的股票ID"})
		return
	}

	err = service.MoveStockUp(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "上移成功"})
}

// MoveStockDown 下移股票
func MoveStockDown(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的股票ID"})
		return
	}

	err = service.MoveStockDown(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "下移成功"})
}

// MoveStockToTop 置顶股票
func MoveStockToTop(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的股票ID"})
		return
	}

	err = service.MoveStockToTop(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "置顶成功"})
}