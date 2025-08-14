// 价格格式化工具

/**
 * 计算价格应该保留的小数位
 * @param {string} open - 开盘价
 * @param {string} yestclose - 昨收价
 * @param {string} price - 当前价格
 * @param {string} high - 最高价
 * @param {string} low - 最低价
 * @returns {number} 应该保留的小数位
 */
export const calcFixedPriceNumber = (
  open = '0',
  yestclose = '0',
  price = '0',
  high = '0',
  low = '0'
) => {
  let reg = /0+$/g;
  open = open.toString().replace(reg, '');
  yestclose = yestclose.toString().replace(reg, '');
  price = price.toString().replace(reg, '');
  high = high.toString().replace(reg, '');
  low = low.toString().replace(reg, '');
  
  let o = open.indexOf('.') === -1 ? 0 : open.length - open.indexOf('.') - 1;
  let yc = yestclose.indexOf('.') === -1 ? 0 : yestclose.length - yestclose.indexOf('.') - 1;
  let p = price.indexOf('.') === -1 ? 0 : price.length - price.indexOf('.') - 1;
  let h = high.indexOf('.') === -1 ? 0 : high.length - high.indexOf('.') - 1;
  let l = low.indexOf('.') === -1 ? 0 : low.length - low.indexOf('.') - 1;
  let max = Math.max(o, yc, p, h, l);
  
  // 对于ETF类股票，强制使用3位小数
  const priceNum = parseFloat(price);
  if (priceNum > 0 && priceNum < 10) {
    return 3; // ETF类通常价格较低，使用3位小数
  }
  
  if (max > 3) {
    max = 2; // 接口返回的指数数值的小数位为4，但习惯两位小数
  }
  
  return max;
};

/**
 * 格式化价格显示
 * @param {number|string} price - 价格
 * @param {number} decimals - 小数位，可选
 * @returns {string} 格式化后的价格
 */
export const formatPrice = (price, decimals = null) => {
  if (price === null || price === undefined || price === '') {
    return '0.00';
  }
  
  const numPrice = parseFloat(price);
  if (isNaN(numPrice)) {
    return '0.00';
  }
  
  if (decimals !== null && decimals >= 0) {
    return numPrice.toFixed(decimals);
  }
  
  // 默认保留2位小数
  return numPrice.toFixed(2);
};

/**
 * 根据股票数据智能确定小数位并格式化价格
 * @param {number|string} price - 价格
 * @param {Object} stockData - 股票完整数据
 * @returns {string} 智能格式化后的价格
 */
export const formatPriceSmart = (price, stockData = {}) => {
  if (price === null || price === undefined || price === '') {
    return '0.00';
  }
  
  const numPrice = parseFloat(price);
  if (isNaN(numPrice)) {
    return '0.00';
  }
  
  // 如果提供了完整的股票数据，使用智能计算
  if (stockData && stockData.open && stockData.yestclose) {
    const decimals = calcFixedPriceNumber(
      stockData.open.toString(),
      stockData.yestclose.toString(),
      price.toString(),
      stockData.high?.toString() || price.toString(),
      stockData.low?.toString() || price.toString()
    );
    return numPrice.toFixed(decimals);
  }
  
  // 默认保留2位小数
  return numPrice.toFixed(2);
};

/**
 * 格式化百分比
 * @param {number|string} value - 百分比值
 * @param {number} decimals - 小数位
 * @returns {string} 格式化后的百分比
 */
export const formatPercentage = (value, decimals = 2) => {
  if (value === null || value === undefined || value === '') {
    return '0.00%';
  }
  
  const numValue = parseFloat(value);
  if (isNaN(numValue)) {
    return '0.00%';
  }
  
  return numValue.toFixed(decimals) + '%';
};