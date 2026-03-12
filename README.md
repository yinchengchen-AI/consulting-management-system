# 咨询公司业务管理系统

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/React-18-61DAFB?style=flat-square&logo=react" alt="React">
  <img src="https://img.shields.io/badge/TypeScript-5.9+-3178C6?style=flat-square&logo=typescript" alt="TypeScript">
  <img src="https://img.shields.io/badge/PostgreSQL-15-336791?style=flat-square&logo=postgresql" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Redis-7-DC382D?style=flat-square&logo=redis" alt="Redis">
  <img src="https://img.shields.io/badge/Docker-20.10+-2496ED?style=flat-square&logo=docker" alt="Docker">
</p>

## 项目简介

咨询公司业务管理系统是一个基于 **Go + React + TypeScript** 的全栈企业级管理系统，专为咨询公司设计，提供项目管理、客户管理、合同管理、财务管理等核心功能。

### 核心特性

- **现代化技术栈**: 采用 Go 1.22 + React 18 + TypeScript 5.9 构建
- **微服务架构**: 支持水平扩展，易于维护
- **响应式设计**: 支持桌面端和移动端访问
- **权限管理**: 基于 RBAC 的细粒度权限控制
- **数据可视化**: 集成 ECharts 提供丰富的图表展示
- **容器化部署**: 一键 Docker 部署，开箱即用

## 技术栈

### 后端技术栈

| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.22+ | 编程语言 |
| Gin | v1.9+ | Web 框架 |
| GORM | v1.25+ | ORM 框架 |
| PostgreSQL | 15 | 关系型数据库 |
| Redis | 7.x | 缓存数据库 |
| JWT | v5 | 身份认证 |
| Swagger | - | API 文档 |

### 前端技术栈

| 技术 | 版本 | 说明 |
|------|------|------|
| React | 18 | UI 框架 |
| TypeScript | 5.9+ | 类型系统 |
| Vite | 7.x | 构建工具 |
| Ant Design | 5.29+ | UI 组件库 |
| Ant Design ProComponents | 2.8+ | 高级组件 |
| React Query | 5.90+ | 数据获取 |
| Zustand | 5.x | 状态管理 |
| React Router | 7.x | 路由管理 |
| ECharts | 5.x | 图表库 |
| Axios | 1.13+ | HTTP 客户端 |

## 功能模块

### 1. 系统管理
- [x] 用户管理（增删改查、角色分配）
- [x] 角色管理（权限配置）
- [x] 部门管理
- [x] 操作日志

### 2. 客户管理
- [x] 客户信息维护
- [x] 客户分类标签
- [x] 联系人管理
- [x] 客户跟进记录

### 3. 项目管理
- [x] 项目立项
- [x] 项目进度跟踪
- [x] 项目成员管理
- [x] 项目文档管理
- [x] 项目里程碑

### 4. 合同管理
- [x] 合同起草
- [x] 合同审批流程
- [x] 合同归档
- [x] 合同提醒

### 5. 财务管理
- [x] 收入管理
- [x] 支出管理
- [x] 发票管理
- [x] 财务报表
- [x] 数据统计分析

### 6. 服务管理
- [x] 服务类型管理
- [x] 服务订单管理
- [x] 服务交付跟踪

## 快速开始

### 环境要求

- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **Go**: 1.22+ (本地开发)
- **Node.js**: 20+ (本地开发)

### 方式一：Docker 一键启动（推荐）

```bash
# 1. 克隆项目
git clone <repository-url>
cd consulting-system

# 2. 创建环境文件
cp .env.example .env

# 3. 启动所有服务
make up

# 4. 查看服务状态
make status
```

服务启动后访问：
- 前端应用: http://localhost
- 后端 API: http://localhost:8080
- API 文档: http://localhost:8080/swagger/index.html

### 方式二：本地开发环境

#### 后端开发

```bash
# 1. 进入后端目录
cd backend

# 2. 安装依赖
go mod download

# 3. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置数据库连接

# 4. 启动开发服务器
go run cmd/main.go

# 后端服务运行在 http://localhost:8080
```

#### 前端开发

```bash
# 1. 进入前端目录
cd frontend

# 2. 安装依赖
npm install

# 3. 启动开发服务器
npm run dev

# 前端服务运行在 http://localhost:5173
```

## 项目结构

```
consulting-system/
├── backend/                    # Go 后端项目
│   ├── cmd/                    # 应用入口
│   │   ├── main.go            # 主程序
│   │   ├── migrate/           # 数据库迁移工具
│   │   └── seed/              # 数据种子工具
│   ├── config/                # 配置文件
│   ├── internal/              # 内部包
│   │   ├── handlers/          # HTTP 处理器
│   │   ├── services/          # 业务逻辑层
│   │   ├── models/            # 数据模型
│   │   ├── middleware/        # 中间件
│   │   └── database/          # 数据库连接
│   ├── pkg/                   # 公共包
│   │   ├── jwt/               # JWT 工具
│   │   ├── logger/            # 日志工具
│   │   └── response/          # 响应封装
│   ├── uploads/               # 上传文件目录
│   ├── Dockerfile             # 后端 Docker 配置
│   └── go.mod                 # Go 模块定义
│
├── frontend/                   # React 前端项目
│   ├── src/
│   │   ├── api/               # API 接口
│   │   ├── components/        # 公共组件
│   │   ├── pages/             # 页面组件
│   │   ├── router/            # 路由配置
│   │   ├── stores/            # 状态管理
│   │   ├── types/             # TypeScript 类型
│   │   └── utils/             # 工具函数
│   ├── public/                # 静态资源
│   ├── Dockerfile             # 前端 Docker 配置
│   └── package.json           # 依赖配置
│
├── docker/                     # Docker 配置
│   ├── nginx/                 # Nginx 配置
│   │   └── nginx.conf         # Nginx 主配置
│   └── postgres/              # PostgreSQL 配置
│       └── init/              # 初始化脚本
│
├── data/                       # 数据持久化目录
│   ├── postgres/              # PostgreSQL 数据
│   └── redis/                 # Redis 数据
│
├── docker-compose.yml          # Docker Compose 配置
├── Makefile                    # 自动化脚本
├── .env.example                # 环境变量模板
└── README.md                   # 项目说明
```

## 开发指南

### 常用命令

```bash
# 查看所有可用命令
make help

# 启动所有服务
make up

# 停止所有服务
make down

# 查看日志
make logs
make logs-backend    # 仅查看后端日志
make logs-frontend   # 仅查看前端日志

# 本地开发
make dev-backend     # 启动后端开发服务器
make dev-frontend    # 启动前端开发服务器

# 数据库操作
make migrate-up      # 执行数据库迁移
make migrate-down    # 回滚数据库迁移
make seed            # 填充种子数据
make db-reset        # 重置数据库

# 代码质量
make test            # 运行所有测试
make lint            # 运行代码检查
make format          # 格式化代码

# 构建
make build           # 构建 Docker 镜像
make build-backend   # 构建后端二进制文件
make build-frontend  # 构建前端生产版本

# 清理
make clean           # 清理 Docker 资源
make clean-data      # 清理所有数据（包括数据库）
```

### 数据库迁移

```bash
# 创建新迁移
make migrate-create name=create_users_table

# 执行迁移
make migrate-up

# 回滚迁移
make migrate-down

# 重置数据库
make db-reset
```

### API 文档

启动服务后，访问 Swagger UI 查看 API 文档：

```
http://localhost:8080/swagger/index.html
```

生成 Swagger 文档：

```bash
make swagger
```

## 生产环境部署

### 1. 环境准备

```bash
# 克隆项目
git clone <repository-url>
cd consulting-system

# 创建生产环境配置
cp .env.example .env
# 编辑 .env，修改以下配置：
# - APP_ENV=production
# - JWT_SECRET=your-strong-secret-key
# - DB_PASSWORD=your-strong-db-password
```

### 2. 启动服务

```bash
# 构建并启动
make up

# 或者使用 docker-compose 直接启动
docker-compose up -d --build
```

### 3. 配置 HTTPS（推荐）

将 SSL 证书放入 `docker/nginx/ssl/` 目录：

```bash
mkdir -p docker/nginx/ssl
cp your-cert.pem docker/nginx/ssl/cert.pem
cp your-key.pem docker/nginx/ssl/key.pem
```

取消 `docker/nginx/nginx.conf` 中 HTTPS 配置的注释。

### 4. 配置反向代理（可选）

如果使用外部 Nginx 作为反向代理：

```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:80;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 环境变量说明

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `APP_ENV` | 应用环境 | `development` |
| `APP_PORT` | 后端端口 | `8080` |
| `DB_HOST` | 数据库主机 | `postgres` |
| `DB_PORT` | 数据库端口 | `5432` |
| `DB_USER` | 数据库用户 | `consulting_user` |
| `DB_PASSWORD` | 数据库密码 | `consulting_pass` |
| `DB_NAME` | 数据库名称 | `consulting_system` |
| `REDIS_HOST` | Redis 主机 | `redis` |
| `REDIS_PORT` | Redis 端口 | `6379` |
| `JWT_SECRET` | JWT 密钥 | - |
| `JWT_EXPIRE_HOURS` | Token 过期时间 | `24` |

完整的环境变量配置请参考 `.env.example` 文件。

## 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

### 代码规范

- 后端代码遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 前端代码遵循 [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript)
- 提交信息遵循 [Conventional Commits](https://www.conventionalcommits.org/)

## 常见问题

### Q: 如何重置数据库？

```bash
make db-reset
```

### Q: 如何查看服务日志？

```bash
# 所有服务
make logs

# 特定服务
make logs-backend
make logs-frontend
make logs-db
```

### Q: 如何修改端口？

编辑 `.env` 文件，修改以下配置：

```env
APP_PORT=8080          # 后端端口
FRONTEND_PORT=80       # 前端端口
DB_PORT=5432           # 数据库端口
REDIS_PORT=6379        # Redis 端口
```

### Q: 如何备份数据库？

```bash
make db-backup
```

备份文件将保存在 `backups/` 目录。

## 更新日志

### v1.0.0 (2024-03-12)

- 初始版本发布
- 实现核心功能模块
- 支持 Docker 一键部署

## 许可证

[MIT](LICENSE)

## 联系方式

如有问题或建议，欢迎提交 Issue 或 Pull Request。

---

<p align="center">
   Made with ❤️ for Consulting Business
</p>
