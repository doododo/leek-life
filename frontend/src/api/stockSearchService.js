import axios from 'axios'

// 股票搜索服务
export class StockSearchService {
  // 搜索股票建议列表
  static async getStockSuggestList(searchText = '') {
    if (!searchText) {
      return [{ label: '请输入关键词查询，如：0000001 或 上证指数' }]
    }

    const result = []
    try {
      // 这里应该调用真实的腾讯数据源API
      const mockStocks = await this.searchStockList(searchText)
      
      mockStocks.forEach((item) => {
        const { code, name, market } = item
        const _code = `${market}${code}`
        
        if (['sz', 'sh', 'bj'].includes(market)) {
          result.push({
            label: `${_code} | ${name}`,
            description: 'A股',
            code: _code
          })
        } else if (['hk'].includes(market)) {
          // 港股个股 || 港股指数
          result.push({
            label: `${_code} | ${name}`,
            description: '港股',
            code: _code
          })
        } else if (['us'].includes(market)) {
          const usCode = _code.split('.')[0] // 去除美股指数.后的内容
          result.push({
            label: `${usCode} | ${name}`,
            description: '美股',
            code: usCode
          })
        }
      })
      
      return result
    } catch (err) {
      console.error('searchStockList error: ', searchText)
      console.error(err)
      return [{ label: '股票查询失败，请重试' }]
    }
  }

  // 腾讯数据源搜索
  static async searchStockList(searchText) {
    if (!searchText.trim()) {
      return []
    }

    try {
      const searchUrl = 'https://proxy.finance.qq.com/ifzqgtimg/appstock/smartbox/search/get'
      
      const stockResponse = await axios.get(searchUrl, {
        params: {
          q: searchText,
        }
      })
      
      const stockListArray = stockResponse?.data?.data?.stock || []
      
      const stockList = stockListArray.map((stockItemArr) => {
        return {
          code: stockItemArr[1].toLowerCase(),
          name: stockItemArr[2],
          market: stockItemArr[0],
          abbreviation: stockItemArr[3],
        }
      })
      
      console.log('stockList: ', stockList, searchText)
      return stockList.slice(0, 10) // 限制返回结果数量
    } catch (error) {
      console.error('腾讯数据源搜索失败:', error)
      
      // 如果真实API失败，回退到模拟数据
      const fallbackDatabase = [
        { code: '00700', name: '腾讯控股', market: 'hk' },
        { code: '03690', name: '美团-W', market: 'hk' },
        { code: '09988', name: '阿里巴巴-SW', market: 'hk' },
        { code: '000001', name: '平安银行', market: 'sz' },
        { code: '000002', name: '万科A', market: 'sz' },
        { code: '600000', name: '浦发银行', market: 'sh' },
        { code: '600519', name: '贵州茅台', market: 'sh' },
        { code: 'AAPL', name: '苹果公司', market: 'us' },
        { code: 'MSFT', name: '微软公司', market: 'us' },
        { code: 'BABA', name: '阿里巴巴', market: 'us' }
      ]

      const filtered = fallbackDatabase.filter(stock => {
        const searchLower = searchText.toLowerCase()
        return stock.code.toLowerCase().includes(searchLower) || 
               stock.name.toLowerCase().includes(searchLower)
      })

      return filtered.slice(0, 10)
    }
  }


}