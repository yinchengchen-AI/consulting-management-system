# ============================================
# 咨询公司业务管理系统 - Makefile
# ============================================
# 版本: 1.0.0
# 描述: 项目自动化构建和部署脚本
# ============================================

.PHONY: help dev build up down logs ps clean install test migrate seed lint format

# 默认目标
.DEFAULT_GOAL := help

# ============================================
# 颜色定义
# ============================================
BLUE := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
NC := \033[0m

# ============================================
# 帮助信息
# ============================================
help: ## 显示帮助信息
	@echo "$(BLUE)========================================$(NC)"
	@echo "$(BLUE)  咨询公司业务管理系统 - 可用命令$(NC)"
	@echo "$(BLUE)========================================$(NC)"
	@echo ""
	@echo "$(GREEN)Docker 命令:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E 'docker|up|down|build|logs' | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)开发命令:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E 'dev|install' | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)测试命令:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E 'test|coverage' | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)数据库命令:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E 'migrate|seed|db' | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)代码质量命令:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E 'lint|format' | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# ============================================
# Docker 命令
# ============================================
dev: ## 启动开发环境（包含热重载）
	@echo "$(GREEN)正在启动开发环境...$(NC)"
	@echo "$(YELLOW)启动数据库和缓存服务...$(NC)"
	docker-compose up -d postgres redis
	@echo "$(GREEN)等待数据库就绪...$(NC)"
	@sleep 5
	@echo "$(GREEN)开发环境已启动!$(NC)"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379"
	@echo ""
	@echo "$(YELLOW)请在单独的终端中运行:$(NC)"
	@echo "  make dev-backend   # 启动后端开发服务器"
	@echo "  make dev-frontend  # 启动前端开发服务器"

up: ## 启动所有 Docker 服务（生产模式）
	@echo "$(GREEN)正在启动所有服务...$(NC)"
	@echo "$(YELLOW)检查环境文件...$(NC)"
	@if [ ! -f .env ]; then \
		echo "$(YELLOW)警告: .env 文件不存在，使用 .env.example 创建$(NC)"; \
		cp .env.example .env; \
	fi
	@echo "$(YELLOW)构建并启动服务...$(NC)"
	docker-compose up -d --build
	@echo "$(GREEN)========================================$(NC)"
	@echo "$(GREEN)  所有服务已启动!$(NC)"
	@echo "$(GREEN)========================================$(NC)"
	@echo "  - 前端应用: http://localhost"
	@echo "  - 后端 API: http://localhost:8080"
	@echo "  - API 文档: http://localhost:8080/swagger/index.html"
	@echo ""
	@echo "$(YELLOW)查看日志: make logs$(NC)"

down: ## 停止所有 Docker 服务
	@echo "$(YELLOW)正在停止所有服务...$(NC)"
	docker-compose down
	@echo "$(GREEN)所有服务已停止$(NC)"

stop: ## 停止所有 Docker 服务（不删除容器）
	@echo "$(YELLOW)正在停止所有服务...$(NC)"
	docker-compose stop

restart: ## 重启所有 Docker 服务
	@echo "$(YELLOW)正在重启所有服务...$(NC)"
	docker-compose restart

build: ## 重新构建所有 Docker 镜像
	@echo "$(GREEN)正在构建 Docker 镜像...$(NC)"
	docker-compose build --no-cache
	@echo "$(GREEN)镜像构建完成$(NC)"

rebuild: down build up ## 完全重建并启动所有服务

logs: ## 查看所有服务日志
	docker-compose logs -f --tail=100

logs-backend: ## 查看后端服务日志
	docker-compose logs -f --tail=100 backend

logs-frontend: ## 查看前端服务日志
	docker-compose logs -f --tail=100 frontend

logs-db: ## 查看数据库日志
	docker-compose logs -f --tail=100 postgres

logs-redis: ## 查看 Redis 日志
	docker-compose logs -f --tail=100 redis

ps: ## 查看运行中的容器
	@echo "$(BLUE)运行中的容器:$(NC)"
	docker-compose ps

status: ps ## 查看服务状态（别名）

# ============================================
# 开发命令
# ============================================
dev-backend: ## 本地启动后端开发服务器
	@echo "$(GREEN)启动后端开发服务器...$(NC)"
	cd backend && go run cmd/main.go

dev-frontend: ## 本地启动前端开发服务器
	@echo "$(GREEN)启动前端开发服务器...$(NC)"
	cd frontend && npm run dev

install-backend: ## 安装后端依赖
	@echo "$(GREEN)安装后端依赖...$(NC)"
	cd backend && go mod download

install-frontend: ## 安装前端依赖
	@echo "$(GREEN)安装前端依赖...$(NC)"
	cd frontend && npm install

install: install-backend install-frontend ## 安装所有依赖

# ============================================
# 测试命令
# ============================================
test: ## 运行所有测试
	@echo "$(GREEN)运行后端测试...$(NC)"
	cd backend && go test -v ./...
	@echo "$(GREEN)运行前端测试...$(NC)"
	cd frontend && npm test

test-backend: ## 运行后端测试
	@echo "$(GREEN)运行后端测试...$(NC)"
	cd backend && go test -v ./...

test-frontend: ## 运行前端测试
	@echo "$(GREEN)运行前端测试...$(NC)"
	cd frontend && npm test

coverage: ## 生成测试覆盖率报告
	@echo "$(GREEN)生成后端测试覆盖率报告...$(NC)"
	cd backend && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)覆盖率报告已生成: backend/coverage.html$(NC)"

# ============================================
# 数据库命令
# ============================================
migrate: migrate-up ## 执行数据库迁移（默认）

migrate-up: ## 执行数据库迁移（升级）
	@echo "$(GREEN)执行数据库迁移...$(NC)"
	@if [ -f backend/cmd/migrate/main.go ]; then \
		cd backend && go run cmd/migrate/main.go up; \
	else \
		echo "$(YELLOW)迁移文件不存在，跳过...$(NC)"; \
	fi

migrate-down: ## 回滚数据库迁移
	@echo "$(YELLOW)回滚数据库迁移...$(NC)"
	@if [ -f backend/cmd/migrate/main.go ]; then \
		cd backend && go run cmd/migrate/main.go down; \
	else \
		echo "$(YELLOW)迁移文件不存在，跳过...$(NC)"; \
	fi

migrate-create: ## 创建新的迁移文件 (用法: make migrate-create name=create_users_table)
	@if [ -z "$(name)" ]; then \
		echo "$(RED)请提供迁移名称: make migrate-create name=your_migration_name$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)创建迁移文件: $(name)$(NC)"
	@if [ -f backend/cmd/migrate/main.go ]; then \
		cd backend && go run cmd/migrate/main.go create $(name); \
	else \
		echo "$(YELLOW)迁移工具不存在，跳过...$(NC)"; \
	fi

seed: ## 填充数据库种子数据
	@echo "$(GREEN)填充种子数据...$(NC)"
	@if [ -f backend/cmd/seed/main.go ]; then \
		cd backend && go run cmd/seed/main.go; \
	else \
		echo "$(YELLOW)种子文件不存在，跳过...$(NC)"; \
	fi

db-reset: migrate-down migrate-up seed ## 重置数据库并填充种子数据

db-shell: ## 进入数据库命令行
	@echo "$(GREEN)进入 PostgreSQL 命令行...$(NC)"
	docker-compose exec postgres psql -U consulting_user -d consulting_system

db-backup: ## 备份数据库
	@echo "$(GREEN)备份数据库...$(NC)"
	@mkdir -p backups
	docker-compose exec postgres pg_dump -U consulting_user consulting_system > backups/consulting_system_$$(date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)数据库备份完成$(NC)"

# ============================================
# 代码质量命令
# ============================================
lint: lint-backend lint-frontend ## 运行所有代码检查

lint-backend: ## 运行后端代码检查
	@echo "$(GREEN)运行后端代码检查...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd backend && golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint 未安装，使用 go vet...$(NC)"; \
		cd backend && go vet ./...; \
	fi

lint-frontend: ## 运行前端代码检查
	@echo "$(GREEN)运行前端代码检查...$(NC)"
	cd frontend && npm run lint

format: ## 格式化代码
	@echo "$(GREEN)格式化后端代码...$(NC)"
	cd backend && gofmt -w .
	@echo "$(GREEN)格式化前端代码...$(NC)"
	cd frontend && npm run format 2>/dev/null || echo "$(YELLOW)前端格式化脚本不存在$(NC)"

# ============================================
# 构建命令
# ============================================
build-backend: ## 构建后端二进制文件
	@echo "$(GREEN)构建后端二进制文件...$(NC)"
	cd backend && go build -o bin/server cmd/main.go
	@echo "$(GREEN)后端构建完成: backend/bin/server$(NC)"

build-frontend: ## 构建前端生产版本
	@echo "$(GREEN)构建前端生产版本...$(NC)"
	cd frontend && npm run build
	@echo "$(GREEN)前端构建完成: frontend/dist/$(NC)"

# ============================================
# 清理命令
# ============================================
clean: ## 清理 Docker 容器、镜像和数据卷
	@echo "$(RED)正在清理 Docker 资源...$(NC)"
	docker-compose down -v --rmi all 2>/dev/null || true
	docker system prune -f
	@echo "$(GREEN)清理完成$(NC)"

clean-data: ## 清理所有数据（包括数据库）
	@echo "$(RED)警告: 这将删除所有数据！$(NC)"
	@read -p "确定要继续吗? [y/N] " confirm && [ $$confirm = y ] || exit 1
	docker-compose down -v
	rm -rf data/
	@echo "$(GREEN)数据已清理$(NC)"

# ============================================
# 部署命令
# ============================================
docker-push: ## 推送 Docker 镜像到仓库
	@echo "$(GREEN)推送 Docker 镜像...$(NC)"
	docker-compose push

deploy: build up ## 构建并部署

# ============================================
# 工具命令
# ============================================
swagger: ## 生成 Swagger 文档
	@echo "$(GREEN)生成 Swagger 文档...$(NC)"
	@if command -v swag >/dev/null 2>&1; then \
		cd backend && swag init -g cmd/main.go -o docs; \
	else \
		echo "$(YELLOW)swag 未安装，正在安装...$(NC)"; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		cd backend && swag init -g cmd/main.go -o docs; \
	fi

generate: ## 运行代码生成工具
	@echo "$(GREEN)运行代码生成...$(NC)"
	cd backend && go generate ./...

# ============================================
# 实用工具
# ============================================
env: ## 创建环境文件
	@if [ ! -f .env ]; then \
		echo "$(GREEN)创建 .env 文件...$(NC)"; \
		cp .env.example .env; \
		echo "$(GREEN).env 文件已创建，请编辑配置$(NC)"; \
	else \
		echo "$(YELLOW).env 文件已存在$(NC)"; \
	fi

shell-backend: ## 进入后端容器
	docker-compose exec backend sh

shell-frontend: ## 进入前端容器
	docker-compose exec frontend sh

shell-db: ## 进入数据库容器
	docker-compose exec postgres sh

# ============================================
# 快捷命令
# ============================================
start: up ## 启动服务（up 的别名）

stop-all: down ## 停止服务（down 的别名）
