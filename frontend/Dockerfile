# ============================================
# 前端 Dockerfile
# 咨询公司业务管理系统
# ============================================
# 多阶段构建：构建阶段 + 运行阶段
# ============================================

# ============================================
# 阶段1：构建阶段
# ============================================
FROM node:20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装构建工具和依赖
RUN apk add --no-cache git

# 安装 pnpm（更快的包管理器）
RUN npm install -g pnpm

# 复制 package.json 和 lock 文件
# 利用 Docker 缓存层，只有这些文件变化时才重新安装依赖
COPY package.json pnpm-lock.yaml* ./

# 安装依赖
# --frozen-lockfile: 确保使用 lock 文件中的精确版本
RUN pnpm install --frozen-lockfile

# 复制源代码
COPY . .

# 构建生产环境代码
# 设置环境变量
ENV NODE_ENV=production
RUN pnpm build

# ============================================
# 阶段2：运行阶段 (Nginx)
# ============================================
FROM nginx:1.25-alpine

# 安装基础工具
RUN apk add --no-cache curl ca-certificates

# 创建 nginx 用户和组（如果不存在）
RUN addgroup -g 101 -S nginx 2>/dev/null || true && \
    adduser -S -D -H -u 101 -h /var/cache/nginx -s /sbin/nologin -G nginx nginx 2>/dev/null || true

# 删除默认的 Nginx 配置
RUN rm -rf /etc/nginx/conf.d/default.conf

# 复制自定义 Nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 从构建阶段复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 设置文件权限
RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chmod -R 755 /usr/share/nginx/html && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /var/log/nginx && \
    chown -R nginx:nginx /etc/nginx/conf.d

# 创建必要的目录并设置权限
RUN mkdir -p /var/run/nginx && \
    chown -R nginx:nginx /var/run/nginx

# 切换到非 root 用户
USER nginx

# 暴露端口
EXPOSE 80

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:80/health || exit 1

# 启动 Nginx
CMD ["nginx", "-g", "daemon off;"]
