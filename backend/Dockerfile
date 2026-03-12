# ============================================
# 后端 Dockerfile
# 咨询公司业务管理系统
# ============================================
# 多阶段构建，优化镜像大小
# ============================================

# ============================================
# 阶段1：构建阶段
# ============================================
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖（利用 Docker 缓存层）
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 构建应用
# - CGO_ENABLED=0: 禁用 CGO，生成静态二进制文件
# - GOOS=linux: 目标操作系统为 Linux
# - -a: 强制重新构建所有包
# - -installsuffix cgo: 使用 cgo 后缀
# - -ldflags: 链接器标志，减小二进制文件大小
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-w -s -X main.Version=1.0.0 -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o main cmd/main.go

# ============================================
# 阶段2：运行阶段
# ============================================
FROM alpine:3.19

# 安装运行时依赖
RUN apk --no-cache add ca-certificates curl tzdata

# 设置时区
ENV TZ=Asia/Shanghai
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 创建非 root 用户运行应用（安全最佳实践）
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 创建必要的目录
RUN mkdir -p /app/uploads /app/logs && \
    chown -R appuser:appgroup /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件（如果有）
COPY --from=builder /app/config ./config

# 设置文件权限
RUN chmod +x /app/main && \
    chown -R appuser:appgroup /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 启动应用
CMD ["./main"]
