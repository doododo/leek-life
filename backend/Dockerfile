# 使用官方 Go 镜像作为构建阶段
FROM golang:1.23-bullseye AS builder

# 设置工作目录
WORKDIR /app

# 安装构建依赖
RUN apt-get update && apt-get install -y gcc libc6-dev

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用Debian作为最终镜像，支持CGO
FROM debian:bullseye-slim

# 设置时区为北京时间
ENV TZ=Asia/Shanghai
RUN apt-get update && apt-get install -y ca-certificates sqlite3 libsqlite3-dev tzdata && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone && \
    rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 创建数据目录
RUN mkdir -p /data

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./main"]