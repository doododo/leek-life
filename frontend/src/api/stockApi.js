import axios from 'axios'

// 创建axios实例
const api = axios.create({
  baseURL: '/api',
  timeout: 5000
})

// 获取所有股票
export const getStocks = () => api.get('/stocks')

// 添加股票
export const addStock = (code) => api.post('/stocks', { code })

// 删除股票
export const deleteStock = (id) => api.delete(`/stocks/${id}`)

// 上移股票
export const moveStockUp = (id) => api.put(`/stocks/${id}/move-up`)

// 下移股票
export const moveStockDown = (id) => api.put(`/stocks/${id}/move-down`)

// 置顶股票
export const moveStockToTop = (id) => api.put(`/stocks/${id}/move-top`)