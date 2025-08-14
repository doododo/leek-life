# 韭菜生活 - 股票关注看板

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)

一个现代化的股票关注看板系统，支持港股实时数据监控、智能轮询和股票关注列表管理。专为中文用户设计的简洁界面，帮助您更好地跟踪关注股票的表现。

![项目截图](https://bu.dusays.com/2025/08/14/689dabd07d5ad.png)

## 🎯 项目特色

- 📊 **实时港股数据** - 支持腾讯控股(hk00700)等港股实时价格监控
- ⏰ **智能交易时间轮询** - 仅在港股交易时间(9:30-16:00)内自动刷新
- 📝 **便捷股票管理** - 一键添加、删除、置顶股票关注列表
- 📈 **清晰涨跌展示** - 今日涨跌和关注以来总涨跌一目了然
- 🎨 **中文友好界面** - 基于Tailwind CSS的现代化响应式设计
- 🐳 **容器化部署** - Docker Compose一键部署，开箱即用

## 🚀 快速部署指南

### 方法1：Docker Compose部署（推荐）

1. **克隆项目**
```bash
git clone git@github.com:doododo/leek-life.git
cd leek-life
```

2. **一键启动**
```bash
docker-compose up -d
```

3. **访问应用**
- 系统界面：http://localhost:8099

### 方法2：本地开发部署

#### 后端服务 (Go)

1. **环境准备**
```bash
cd backend
```

2. **安装依赖**
```bash
go mod download
```

3. **启动服务**
```bash
go run main.go
```

服务将在 http://localhost:8080 启动

#### 前端服务 (Vue 3)

1. **环境准备**
```bash
cd frontend
```

2. **安装依赖**
```bash
npm install
```

3. **开发模式** (带代理配置)
```bash
npm run dev
```

服务将在 http://localhost:5173 启动，已配置代理到后端 http://localhost:8080

4. **构建生产版本**
```bash
npm run build
```

5. **本地预览生产版本**
```bash
npm run preview
```

**注意**：本地开发时，前端开发服务器已配置代理，无需额外的Nginx配置。API请求会自动转发到后端服务。

## 📱 使用教程

### 添加股票关注

1. 在页面顶部输入框输入股票代码
2. 支持的格式：
   - 港股：`hk00700`（腾讯控股示例）
3. 点击"添加"按钮，系统自动获取股票信息

### 管理关注列表

| 操作 | 图标 | 说明 |
|------|------|------|
| **置顶** | ⇈ | 将股票移至列表顶部 |
| **上移** | ↑ | 向上移动一位 |
| **下移** | ↓ | 向下移动一位 |
| **删除** | ✕ | 取消关注 |

### 数据指标说明

- **当前价格**：实时股票价格
- **关注价格**：您添加关注时的价格
- **今日涨跌**：当日价格变化百分比
- **关注涨跌**：从关注至今的总变化
- **关注日期**：添加关注的具体日期

### 智能刷新机制

- **交易时间**：工作日9:30-16:00，每5秒自动刷新
- **非交易时间**：停止自动刷新，节省资源

## 🔧 技术架构

### 后端技术栈
- **Go 1.23** - 高性能后端语言
- **Gin Web框架** - 轻量级HTTP框架
- **GORM** - 优雅的ORM库
- **SQLite** - 轻量级数据库，开箱即用
- **CGO支持** - 完整支持SQLite功能

### 前端技术栈
- **Vue 3** - 渐进式JavaScript框架
- **Vite** - 现代化构建工具
- **Tailwind CSS** - 实用优先的CSS框架
- **Axios** - HTTP客户端
- **Nginx** - 生产环境Web服务器

### 部署架构
- **Docker容器化** - 环境一致性保证
- **Docker Compose** - 服务编排管理
- **Debian基础镜像** - 解决Alpine CGO兼容性问题
- **多阶段构建** - 优化镜像大小

## 📊 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   用户浏览器     │────│   Nginx前端     │────│   Go后端服务     │
│   localhost:8099 │    │   localhost:8099│    │   localhost:8098│
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                        │
                                │                        │
                       ┌─────────────────┐    ┌─────────────────┐
                       │   静态资源服务   │    │   SQLite数据库  │
                       │   (Vue.js构建)  │    │   (数据持久化)  │
                       └─────────────────┘    └─────────────────┘
```

## 🛠️ 开发指南

### 项目结构
```
Docker/stock/
├── backend/           # Go后端服务
│   ├── controller/    # 控制器层
│   ├── database/      # 数据库配置
│   ├── model/         # 数据模型
│   ├── router/        # 路由定义
│   ├── service/       # 业务逻辑层
│   └── main.go        # 程序入口
├── frontend/          # Vue前端应用
│   ├── src/
│   │   ├── components/ # Vue组件
│   │   ├── api/       # API接口封装
│   │   └── ...
│   ├── vite.config.js # Vite配置（含开发代理）
│   ├── nginx.conf     # 生产环境Nginx配置
│   └── Dockerfile     # 前端构建配置
└── docker-compose.yml # 服务编排
```

### API接口文档

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/stocks` | 获取股票列表 |
| POST | `/api/stocks` | 添加股票关注 |
| DELETE | `/api/stocks/:id` | 删除股票关注 |
| PUT | `/api/stocks/:id/move-up` | 上移股票位置 |
| PUT | `/api/stocks/:id/move-down` | 下移股票位置 |
| PUT | `/api/stocks/:id/move-top` | 置顶股票 |

### 环境变量

#### 后端配置
- `PORT`: 服务端口（默认8080）
- `GIN_MODE`: 运行模式（debug/release）

#### Docker配置
- 前端端口：8099
- 后端端口：8098
- 网络：stock-stock-network

## 🐛 常见问题

### Q1: 容器启动失败怎么办？
```bash
# 检查服务状态
docker-compose ps

# 查看详细日志
docker-compose logs backend
docker-compose logs frontend

# 重新构建
docker-compose build --no-cache
docker-compose up -d
```

### Q2: 数据库连接问题？
- 确保后端容器有/data目录写入权限
- 检查SQLite文件权限

### Q3: 前端API调用失败？
- 确认nginx代理配置正确
- 检查网络连接和端口映射

## 🤝 贡献指南

欢迎提交Issue和Pull Request！

### 开发流程
1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 代码规范
- 后端：遵循Go代码规范，使用gofmt格式化
- 前端：遵循Vue风格指南，使用ESLint检查

## 📄 开源协议

本项目采用 [MIT License](LICENSE) 开源协议，可自由使用、修改和分发。

## 🙏 致谢

感谢以下开源项目：
- [Gin](https://github.com/gin-gonic/gin) - Go Web框架
- [Vue.js](https://vuejs.org/) - 前端框架
- [Tailwind CSS](https://tailwindcss.com/) - CSS框架
- [SQLite](https://www.sqlite.org/) - 数据库

## 📞 支持

如有问题，请通过以下方式联系：
- 提交Issue：[GitHub Issues](https://github.com/doododo/leek-life/issues)
- 邮件反馈：saybye720@gmail.com

---

<div align="center">

**🎉 祝投资顺利，远离韭菜命运！** 

</div>