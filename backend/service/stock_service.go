package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"stock-app/database"
	"stock-app/model"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GetAllStocks 获取所有股票
func GetAllStocks() ([]model.Stock, error) {
	var stocks []model.Stock
	result := database.DB.Order("order_index ASC").Find(&stocks)
	return stocks, result.Error
}

// AddStock 添加股票
func AddStock(code string) error {
	// 从API获取股票信息
	stockInfo, err := fetchStockInfo(code)
	if err != nil {
		return err
	}

	// 获取当前最大的order_index
	var maxOrder int
	database.DB.Model(&model.Stock{}).Select("COALESCE(MAX(order_index), 0)").Scan(&maxOrder)

	// 设置添加日期和初始价格
	stockInfo.AddDate = time.Now()
	stockInfo.AddPrice = stockInfo.Price
	stockInfo.TotalChange = 0
	stockInfo.OrderIndex = maxOrder + 1

	// 保存到数据库
	result := database.DB.Create(stockInfo)
	return result.Error
}

// DeleteStock 删除股票
func DeleteStock(id uint) error {
	result := database.DB.Delete(&model.Stock{}, id)
	return result.Error
}

// MoveStockUp 上移股票
func MoveStockUp(id uint) error {
	var currentStock model.Stock
	if err := database.DB.First(&currentStock, id).Error; err != nil {
		return err
	}

	// 找到当前股票前面的股票
	var prevStock model.Stock
	if err := database.DB.Where("order_index < ?", currentStock.OrderIndex).Order("order_index DESC").First(&prevStock).Error; err != nil {
		return err // 已经是第一个了
	}

	// 交换order_index
	currentStock.OrderIndex, prevStock.OrderIndex = prevStock.OrderIndex, currentStock.OrderIndex

	// 更新数据库
	database.DB.Save(&currentStock)
	database.DB.Save(&prevStock)

	return nil
}

// MoveStockDown 下移股票
func MoveStockDown(id uint) error {
	var currentStock model.Stock
	if err := database.DB.First(&currentStock, id).Error; err != nil {
		return err
	}

	// 找到当前股票后面的股票
	var nextStock model.Stock
	if err := database.DB.Where("order_index > ?", currentStock.OrderIndex).Order("order_index ASC").First(&nextStock).Error; err != nil {
		return err // 已经是最后一个了
	}

	// 交换order_index
	currentStock.OrderIndex, nextStock.OrderIndex = nextStock.OrderIndex, currentStock.OrderIndex

	// 更新数据库
	database.DB.Save(&currentStock)
	database.DB.Save(&nextStock)

	return nil
}

// MoveStockToTop 置顶股票
func MoveStockToTop(id uint) error {
	var currentStock model.Stock
	if err := database.DB.First(&currentStock, id).Error; err != nil {
		return err
	}

	// 找到最小的order_index
	var minOrder int
	database.DB.Model(&model.Stock{}).Select("MIN(order_index)").Scan(&minOrder)

	// 如果已经是第一个了，不需要操作
	if currentStock.OrderIndex == minOrder {
		return nil
	}

	// 将所有order_index小于当前股票的记录都+1
	database.DB.Model(&model.Stock{}).Where("order_index < ?", currentStock.OrderIndex).Update("order_index", database.DB.Raw("order_index + 1"))

	// 将当前股票设置为最小order_index
	currentStock.OrderIndex = minOrder
	database.DB.Save(&currentStock)

	return nil
}

// InitializeOrderIndices 初始化已有股票的order_index
func InitializeOrderIndices() error {
	var stocks []model.Stock
	// 获取所有order_index为0的股票，按ID排序
	result := database.DB.Where("order_index = 0").Order("id ASC").Find(&stocks)
	if result.Error != nil {
		return result.Error
	}

	// 为这些股票设置递增的order_index
	for i, stock := range stocks {
		stock.OrderIndex = i + 1
		database.DB.Save(&stock)
	}

	log.Printf("已初始化%d个股票的排序索引", len(stocks))
	return nil
}

// isTradingHours 检查是否在交易时间内（工作日9:00-16:30）
func isTradingHours() bool {
	now := time.Now()

	// 检查是否为工作日（周一到周五）
	weekday := now.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	// 检查时间是否在9:00-16:30之间
	hour := now.Hour()
	minute := now.Minute()

	// 开始时间: 9:00
	if hour < 9 {
		return false
	}

	// 结束时间: 16:10
	if hour > 16 || (hour == 16 && minute > 30) {
		return false
	}

	return true
}

// StockData 股票数据结构
type StockData struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	YestClose float64 `json:"yestclose"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Volume    float64 `json:"volume"`
	Amount    float64 `json:"amount"`
	Buy1      float64 `json:"buy1"`
	Sell1     float64 `json:"sell1"`
	Time      string  `json:"time"`
}

// fetchBatchStockData 批量获取股票数据
func fetchBatchStockData(codes []string) (map[string]*StockData, error) {
	if len(codes) == 0 {
		return make(map[string]*StockData), nil
	}

	// 构建请求参数，使用 r_ 前缀
	rCodes := make([]string, len(codes))
	for i, code := range codes {
		rCodes[i] = "r_" + code
	}

	url := fmt.Sprintf("https://qt.gtimg.cn/q=%s&fmt=json", strings.Join(rCodes, ","))

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,en;q=0.6")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 处理编码
	bodyStr := string(body)
	if !utf8.ValidString(bodyStr) {
		bodyStr = convertToUTF8(body)
	}

	// 解析JSON响应
	var responseData map[string]interface{}
	if err := json.Unmarshal([]byte(bodyStr), &responseData); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	result := make(map[string]*StockData)

	// 处理每个股票的数据
	for _, code := range codes {
		rCode := "r_" + code
		if stockArr, exists := responseData[rCode]; exists {
			if arr, ok := stockArr.([]interface{}); ok && len(arr) > 37 {
				stockData := &StockData{
					Code: code,
				}

				// 解析各个字段
				if arr[1] != nil {
					stockData.Name = fmt.Sprintf("%v", arr[1])
				}
				if arr[3] != nil {
					if price, err := strconv.ParseFloat(fmt.Sprintf("%v", arr[3]), 64); err == nil {
						stockData.Price = price
					}
				}
				if arr[4] != nil {
					if yestClose, err := strconv.ParseFloat(fmt.Sprintf("%v", arr[4]), 64); err == nil {
						stockData.YestClose = yestClose
					}
				}
				if arr[5] != nil {
					if open, err := strconv.ParseFloat(fmt.Sprintf("%v", arr[5]), 64); err == nil {
						stockData.Open = open
					}
				}
				if len(arr) > 33 && arr[33] != nil {
					if high, err := strconv.ParseFloat(fmt.Sprintf("%v", arr[33]), 64); err == nil {
						stockData.High = high
					}
				}
				if len(arr) > 34 && arr[34] != nil {
					if low, err := strconv.ParseFloat(fmt.Sprintf("%v", arr[34]), 64); err == nil {
						stockData.Low = low
					}
				}
				if len(arr) > 36 && arr[36] != nil {
					if volume, err := strconv.ParseFloat(fmt.Sprintf("%v", arr[36]), 64); err == nil {
						stockData.Volume = volume
					}
				}
				if len(arr) > 37 && arr[37] != nil {
					if amount, err := strconv.ParseFloat(fmt.Sprintf("%v", arr[37]), 64); err == nil {
						stockData.Amount = amount
					}
				}
				if len(arr) > 30 && arr[30] != nil {
					stockData.Time = fmt.Sprintf("%v", arr[30])
				}

				result[code] = stockData
			}
		}
	}

	return result, nil
}

// fetchStockInfo 从腾讯财经API获取股票信息
func fetchStockInfo(code string) (*model.Stock, error) {
	url := "https://qt.gtimg.cn/q=" + code

	// 创建HTTP客户端请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头，确保正确处理编码
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,en;q=0.6")
	req.Header.Set("Accept-Charset", "utf-8,gb2312,gbk;q=0.7,*;q=0.3")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 确保字符串是有效的UTF-8
	bodyStr := string(body)
	if !utf8.ValidString(bodyStr) {
		// 尝试从GBK转换到UTF-8
		bodyStr = convertToUTF8(body)
	}

	parts := strings.Split(bodyStr, "~")
	if len(parts) < 40 {
		return nil, fmt.Errorf("解析数据失败: 数据格式不正确，期望至少40个字段，得到%d个", len(parts))
	}

	// 根据腾讯API的实际数据结构解析
	// 通过分析实际数据发现正确的字段位置：
	// parts[1] = 股票名称
	// parts[3] = 当前价格
	// parts[4] = 昨收价格
	// 通过计算得到正确的涨跌幅

	price, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return nil, fmt.Errorf("解析价格失败: %v", err)
	}

	yesterdayPrice, err := strconv.ParseFloat(parts[4], 64)
	if err != nil {
		return nil, fmt.Errorf("解析昨日价格失败: %v", err)
	}

	// 计算今日涨跌幅
	var todayChange float64 = 0
	if yesterdayPrice != 0 {
		todayChange = ((price - yesterdayPrice) / yesterdayPrice) * 100
	}

	// 添加调试日志
	log.Printf("股票 %s: 当前价格=%.3f, 昨收价格=%.3f, 计算涨跌幅=%.2f%%",
		code, price, yesterdayPrice, todayChange)

	return &model.Stock{
		Name:        strings.TrimSpace(parts[1]),
		Code:        code,
		Price:       price,
		TodayChange: todayChange,
	}, nil
}

// convertToUTF8 尝试将GBK编码的字节转换为UTF-8字符串
func convertToUTF8(data []byte) string {
	// 尝试从GBK转换
	reader := transform.NewReader(strings.NewReader(string(data)), simplifiedchinese.GBK.NewDecoder())
	result, err := io.ReadAll(reader)
	if err != nil {
		// 如果转换失败，返回原始字符串
		return string(data)
	}
	return string(result)
}

// StartPriceUpdater 启动定时更新股票价格的goroutine
func StartPriceUpdater() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			// 只在交易时间内更新
			if isTradingHours() {
				updateAllStockPrices()
			}
		}
	}()
}

// updateAllStockPrices 更新所有股票的价格
func updateAllStockPrices() {
	var stocks []model.Stock
	result := database.DB.Find(&stocks)
	if result.Error != nil {
		log.Printf("获取股票列表失败: %v", result.Error)
		return
	}

	if len(stocks) == 0 {
		return
	}

	// 收集所有股票代码
	codes := make([]string, len(stocks))
	stockMap := make(map[string]*model.Stock)
	for i := range stocks {
		codes[i] = stocks[i].Code
		stockMap[stocks[i].Code] = &stocks[i]
	}

	// 批量获取股票数据
	stockDataMap, err := fetchBatchStockData(codes)
	if err != nil {
		log.Printf("批量获取股票数据失败: %v", err)
		return
	}

	// 更新每个股票的数据
	for code, stockData := range stockDataMap {
		if stock, exists := stockMap[code]; exists {
			// 更新价格和今日涨幅
			if stockData.Price > 0 {
				stock.Price = stockData.Price
			}

			// 计算今日涨幅
			if stockData.YestClose > 0 {
				stock.TodayChange = ((stockData.Price - stockData.YestClose) / stockData.YestClose) * 100
			}

			// 计算关注至今涨幅
			if stock.AddPrice > 0 {
				stock.TotalChange = ((stockData.Price - stock.AddPrice) / stock.AddPrice) * 100
			}

			// 保存更新
			database.DB.Save(stock)
		}
	}

	log.Printf("批量更新了%d个股票的价格", len(stockDataMap))
}
