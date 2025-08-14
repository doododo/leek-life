<template>
  <div class="container mx-auto px-4 py-8">
    <!-- 顶部导航栏 -->
    <div class="flex items-center justify-between mb-8">
      <!-- Logo - 左上角 -->
      <div class="text-2xl font-bold text-gray-800">
        Leek Life
      </div>

      <!-- 搜索区域 - 中间靠右 -->
      <div class="flex items-center space-x-4">
        <div class="relative w-96">
          <input
            v-model="searchText"
            type="text"
            placeholder="输入股票代码或名称搜索 (如 00700 或 腾讯)"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            @input="handleSearch"
            @focus="showDropdown = true"
            @blur="hideDropdown"
          >
          
          <!-- 搜索下拉框 -->
          <div 
            v-if="showDropdown && searchResults.length > 0" 
            class="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-60 overflow-y-auto"
          >
            <div
              v-for="(result, index) in searchResults"
              :key="index"
              class="px-4 py-3 hover:bg-gray-50 cursor-pointer border-b border-gray-100 last:border-b-0"
              @click="selectStock(result)"
            >
              <div class="font-medium text-gray-900">{{ result.label }}</div>
              <div v-if="result.description" class="text-sm text-gray-500">{{ result.description }}</div>
            </div>
          </div>
          
          <!-- 无搜索结果 -->
          <div 
            v-if="showDropdown && searchText && searchResults.length === 0 && !isSearching" 
            class="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg p-4 text-center text-gray-500"
          >
            未找到相关股票，请尝试其他关键词
          </div>
          
          <!-- 搜索中状态 -->
          <div 
            v-if="isSearching" 
            class="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg p-4 text-center text-gray-500"
          >
            搜索中...
          </div>
          
          <!-- 提示信息 -->
          <div 
            v-if="showDropdown && searchResults.length === 1 && searchResults[0].code === null" 
            class="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg p-4 text-center text-gray-500"
          >
            {{ searchResults[0].label }}
          </div>
        </div>

        <!-- GitHub图标 - 最右侧 -->
        <a 
          href="https://github.com/your-username/stock-tracker" 
          target="_blank" 
          rel="noopener noreferrer"
          class="text-gray-600 hover:text-gray-800 transition-colors"
          title="View on GitHub"
        >
          <svg 
            class="w-8 h-8" 
            fill="currentColor" 
            viewBox="0 0 24 24" 
            xmlns="http://www.w3.org/2000/svg"
          >
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
          </svg>
        </a>
      </div>
    </div>

    <!-- 股票列表 -->
    <div class="overflow-x-auto">
      <table class="min-w-full bg-white border border-gray-200 rounded-lg">
        <thead>
          <tr class="bg-gray-100 text-gray-700">
            <th class="py-3 px-4 text-left border-b">股票名称</th>
            <th class="py-3 px-4 text-left border-b">股票代码</th>
            <th class="py-3 px-4 text-left border-b">当前价格</th>
            <th class="py-3 px-4 text-left border-b">今日涨幅</th>
            <th class="py-3 px-4 text-left border-b">关注价格</th>
            <th class="py-3 px-4 text-left border-b">关注涨幅</th>
            <th class="py-3 px-4 text-left border-b">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="stock in stocks" :key="stock.id" class="hover:bg-gray-50 transition-colors">
            <td class="py-3 px-4 border-b">{{ stock.name }}</td>
            <td class="py-3 px-4 border-b">{{ stock.code }}</td>
            <td class="py-3 px-4 border-b">{{ formatPriceSmart(stock.price, stock) }}</td>
            <td
              class="py-3 px-4 border-b"
              :class="stock.today_change > 0 ? 'text-red-600' : stock.today_change < 0 ? 'text-green-600' : 'text-gray-600'"
            >
              {{ stock.today_change >= 0 ? '+' : '' }}{{ formatPercentage(stock.today_change) }}
            </td>
            <td class="py-3 px-4 border-b">{{ formatPriceSmart(stock.add_price, stock) }}</td>
            <td
              class="py-3 px-4 border-b"
              :class="stock.total_change > 0 ? 'text-red-600' : stock.total_change < 0 ? 'text-green-600' : 'text-gray-600'"
            >
              <div>
                {{ stock.total_change >= 0 ? '+' : '' }}{{ formatPercentage(stock.total_change) }}
              </div>
              <div class="text-xs text-gray-500 mt-1">
                ({{ new Date(stock.add_date).toLocaleDateString('zh-CN') }})
              </div>
            </td>
            <td class="py-3 px-4 border-b">
              <div class="flex gap-1">
                <button
                  @click="handleMoveToTop(stock.id)"
                  class="px-2 py-1 bg-blue-600 text-white text-xs rounded hover:bg-blue-700 transition-colors"
                  title="置顶"
                >
                  ⇈
                </button>
                <button
                  @click="handleMoveUp(stock.id)"
                  class="px-2 py-1 bg-green-600 text-white text-xs rounded hover:bg-green-700 transition-colors"
                  title="上移"
                >
                  ↑
                </button>
                <button
                  @click="handleMoveDown(stock.id)"
                  class="px-2 py-1 bg-yellow-600 text-white text-xs rounded hover:bg-yellow-700 transition-colors"
                  title="下移"
                >
                  ↓
                </button>
                <button
                  @click="handleDeleteStock(stock.id)"
                  class="px-2 py-1 bg-red-600 text-white text-xs rounded hover:bg-red-700 transition-colors"
                  title="取消关注"
                >
                  ✕
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 空状态 -->
    <div v-if="stocks.length === 0" class="text-center py-12 text-gray-500">
      <p>暂无关注的股票，请添加股票代码</p>
    </div>
  </div>
</template>

<!-- Tailwind CSS 通过CDN加载 -->

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { getStocks, addStock, deleteStock, moveStockUp, moveStockDown, moveStockToTop } from '../api/stockApi'
import { StockSearchService } from '../api/stockSearchService'
import { formatPriceSmart, formatPercentage } from '../utils/priceFormatter'

const stocks = ref([])
const searchText = ref('')
const searchResults = ref([])
const showDropdown = ref(false)
const isSearching = ref(false)
const selectedStock = ref(null)
let refreshInterval = null
let searchTimeout = null

// 检查是否在交易时间内（工作日9:00-16:30）
const isTradingHours = () => {
  const now = new Date()
  
  // 检查是否为工作日（周一到周五）
  const weekday = now.getDay()
  if (weekday === 0 || weekday === 6) { // 0=Sunday, 6=Saturday
    return false
  }
  
  // 检查时间是否在9:00-16:30之间
  const hour = now.getHours()
  const minute = now.getMinutes()
  
  // 开始时间: 9:00
  if (hour < 9) {
    return false
  }
  
  // 结束时间: 16:30
  if (hour > 16 || (hour === 16 && minute > 30)) {
    return false
  }
  
  return true
}

// 获取股票列表
const fetchStocks = async () => {
  try {
    const response = await getStocks()
    stocks.value = response.data
  } catch (error) {
    console.error('获取股票列表失败:', error)
    alert('获取股票列表失败，请重试')
  }
}

// 搜索股票
const handleSearch = () => {
  if (!searchText.value.trim()) {
    searchResults.value = []
    return
  }

  // 清除之前的搜索定时器
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }

  // 设置防抖搜索
  searchTimeout = setTimeout(async () => {
    await searchStocks(searchText.value)
  }, 300)
}

// 搜索股票方法
const searchStocks = async (searchText) => {
  if (!searchText.trim()) {
    searchResults.value = []
    return
  }

  isSearching.value = true
  
  try {
    const results = await StockSearchService.getStockSuggestList(searchText)
    
    // 过滤掉提示信息
    searchResults.value = results.filter(item => item.code)
    
    // 如果没有搜索结果
    if (searchResults.value.length === 0 && searchText.trim()) {
      searchResults.value = [{ 
        label: '未找到相关股票', 
        description: '请尝试其他关键词',
        code: null 
      }]
    }
  } catch (error) {
    console.error('搜索股票失败:', error)
    searchResults.value = [{ 
      label: '搜索失败，请重试', 
      description: '网络错误',
      code: null 
    }]
  } finally {
    isSearching.value = false
  }
}

// 选择股票
const selectStock = async (stock) => {
  // 如果没有股票代码，不执行添加
  if (!stock || !stock.code) {
    return
  }
  
  selectedStock.value = stock
  searchText.value = stock.label.split(' | ')[0]
  showDropdown.value = false
  
  // 直接添加选中的股票
  try {
    await addStock(stock.code)
    searchText.value = ''
    searchResults.value = []
    selectedStock.value = null
    await fetchStocks()
    alert('添加成功')
  } catch (error) {
    console.error('添加股票失败:', error)
    alert('添加股票失败，请重试')
  }
}

// 隐藏下拉框（延迟关闭，以便点击事件生效）
const hideDropdown = () => {
  setTimeout(() => {
    showDropdown.value = false
  }, 200)
}

// 添加股票
const handleAddStock = async () => {
  if (!searchText.value.trim()) {
    alert('请输入股票代码')
    return
  }

  // 如果输入的是完整代码，直接添加
  try {
    await addStock(searchText.value)
    searchText.value = ''
    searchResults.value = []
    await fetchStocks()
    alert('添加成功')
  } catch (error) {
    console.error('添加股票失败:', error)
    alert('添加股票失败，请检查代码是否正确')
  }
}

// 删除股票
const handleDeleteStock = async (id) => {
  if (confirm('确定要取消关注此股票吗？')) {
    try {
      await deleteStock(id)
      await fetchStocks()
      alert('取消关注成功')
    } catch (error) {
      console.error('删除股票失败:', error)
      alert('删除股票失败，请重试')
    }
  }
}

// 上移股票
const handleMoveUp = async (id) => {
  try {
    await moveStockUp(id)
    await fetchStocks()
  } catch (error) {
    console.error('上移失败:', error)
    alert('上移失败，可能已经是第一个了')
  }
}

// 下移股票
const handleMoveDown = async (id) => {
  try {
    await moveStockDown(id)
    await fetchStocks()
  } catch (error) {
    console.error('下移失败:', error)
    alert('下移失败，可能已经是最后一个了')
  }
}

// 置顶股票
const handleMoveToTop = async (id) => {
  try {
    await moveStockToTop(id)
    await fetchStocks()
  } catch (error) {
    console.error('置顶失败:', error)
    alert('置顶失败，请重试')
  }
}

// 初始化
onMounted(() => {
  fetchStocks()
  // 每5秒检查一次，只在交易时间内轮询
  refreshInterval = setInterval(() => {
    if (isTradingHours()) {
      fetchStocks()
    }
  }, 5000)
})

// 清理
onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>