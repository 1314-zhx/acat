# 第一阶段：构建
FROM golang:1.23-alpine AS builder

WORKDIR /acat
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# 第二阶段：运行
FROM alpine:latest

# 安装必要依赖
RUN apk --no-cache add ca-certificates tzdata

# 创建非 root 用户
RUN adduser -D -s /bin/sh appuser

# 创建应用目录并授权
RUN mkdir -p /acat/logs /acat/uploads && \
    chown -R appuser:appuser /acat

# 切换到非 root 用户
USER appuser

# 设置工作目录
WORKDIR /acat

# 复制二进制文件
COPY --from=builder /acat/main .

# 声明端口
EXPOSE 9090

# 启动应用
CMD ["./main"]