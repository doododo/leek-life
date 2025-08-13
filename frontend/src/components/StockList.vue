<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold mb-6 text-center">股票看板</h1>

    <!-- 添加股票表单 -->
    <div class="flex justify-center mb-8">
      <div class="flex w-full max-w-md">
        <input
          v-model="stockCode"
          type="text"
          placeholder="输入股票代码 (如 hk00700)"
          class="flex-1 px-4 py-2 border border-gray-300 rounded-l-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
        <button
          @click="handleAddStock"
          class="px-4 py-2 bg-blue-600 text-white rounded-r-lg hover:bg-blue-700 transition-colors"
        >
          添加
        </button>
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
            <th class="py-3 px-4 text-left border-b">关注价格</th>
            <th class="py-3 px-4 text-left border-b">今日涨幅</th>
            <th class="py-3 px-4 text-left border-b">关注涨幅</th>
            <th class="py-3 px-4 text-left border-b">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="stock in stocks" :key="stock.id" class="hover:bg-gray-50 transition-colors">
            <td class="py-3 px-4 border-b">{{ stock.name }}</td>
            <td class="py-3 px-4 border-b">{{ stock.code }}</td>
            <td class="py-3 px-4 border-b">{{ stock.price.toFixed(2) }}</td>
            <td class="py-3 px-4 border-b">{{ stock.add_price.toFixed(2) }}</td>
            <td
              class="py-3 px-4 border-b"
              :class="stock.today_change > 0 ? 'text-red-600' : stock.today_change < 0 ? 'text-green-600' : 'text-gray-600'"
            >
              {{ stock.today_change >= 0 ? '+' : '' }}{{ stock.today_change.toFixed(2) }}%
            </td>
            <td
              class="py-3 px-4 border-b"
              :class="stock.total_change > 0 ? 'text-red-600' : stock.total_change < 0 ? 'text-green-600' : 'text-gray-600'"
            >
              <div>
                {{ stock.total_change >= 0 ? '+' : '' }}{{ stock.total_change.toFixed(2) }}%
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

const stocks = ref([])
const stockCode = ref('')
let refreshInterval = null

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

// 添加股票
const handleAddStock = async () => {
  if (!stockCode.value.trim()) {
    alert('请输入股票代码')
    return
  }

  try {
    await addStock(stockCode.value)
    stockCode.value = ''
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