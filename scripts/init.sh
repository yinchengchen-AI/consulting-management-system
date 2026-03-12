#!/bin/bash
# ============================================
# 咨询公司业务管理系统 - 初始化脚本
# ============================================
# 用法: ./scripts/init.sh
# 说明: 首次运行项目时执行此脚本
# ============================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# 切换到项目目录
cd "$PROJECT_DIR"

# 打印信息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

# 打印成功
print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# 打印警告
print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 打印错误
print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 打印分隔线
print_line() {
    echo -e "${BLUE}========================================${NC}"
}

# 检查命令是否存在
check_command() {
    if ! command -v "$1" &> /dev/null; then
        return 1
    fi
    return 0
}

# 检查 Docker
check_docker() {
    print_info "检查 Docker 环境..."
    
    if ! check_command docker; then
        print_error "Docker 未安装"
        echo "请访问 https://docs.docker.com/get-docker/ 安装 Docker"
        exit 1
    fi
    
    if ! check_command docker-compose; then
        print_error "Docker Compose 未安装"
        echo "请访问 https://docs.docker.com/compose/install/ 安装 Docker Compose"
        exit 1
    fi
    
    # 检查 Docker 服务
    if ! docker info &> /dev/null; then
        print_error "Docker 服务未运行"
        echo "请启动 Docker 服务:"
        echo "  - macOS/Windows: 启动 Docker Desktop"
        echo "  - Linux: sudo systemctl start docker"
        exit 1
    fi
    
    print_success "Docker 环境检查通过"
}

# 检查 Go
check_go() {
    print_info "检查 Go 环境..."
    
    if ! check_command go; then
        print_warning "Go 未安装（可选，仅本地开发需要）"
        return
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go 版本: $GO_VERSION"
}

# 检查 Node.js
check_node() {
    print_info "检查 Node.js 环境..."
    
    if ! check_command node; then
        print_warning "Node.js 未安装（可选，仅本地开发需要）"
        return
    fi
    
    NODE_VERSION=$(node --version)
    print_success "Node.js 版本: $NODE_VERSION"
}

# 创建目录结构
create_directories() {
    print_info "创建目录结构..."
    
    mkdir -p data/postgres
    mkdir -p data/redis
    mkdir -p backups
    mkdir -p backend/uploads
    mkdir -p backend/logs
    mkdir -p docker/nginx/ssl
    
    print_success "目录结构创建完成"
}

# 创建环境文件
create_env_file() {
    print_info "创建环境配置文件..."
    
    if [ -f .env ]; then
        print_warning ".env 文件已存在"
        read -p "是否覆盖? [y/N] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "保留现有 .env 文件"
            return
        fi
    fi
    
    cp .env.example .env
    print_success ".env 文件已创建"
    print_warning "请编辑 .env 文件，修改以下配置:"
    echo "  - DB_PASSWORD: 数据库密码"
    echo "  - JWT_SECRET: JWT 密钥"
    echo "  - REDIS_PASSWORD: Redis 密码（可选）"
}

# 生成随机密钥
generate_secrets() {
    print_info "生成安全密钥..."
    
    JWT_SECRET=$(openssl rand -base64 32 2>/dev/null || head -c 32 /dev/urandom | base64)
    DB_PASSWORD=$(openssl rand -base64 16 2>/dev/null || head -c 16 /dev/urandom | base64)
    
    print_success "随机密钥已生成"
    print_info "JWT_SECRET: $JWT_SECRET"
    print_info "DB_PASSWORD: $DB_PASSWORD"
    print_warning "请将这些密钥复制到 .env 文件中"
}

# 启动基础服务
start_base_services() {
    print_info "启动基础服务（PostgreSQL + Redis）..."
    
    docker-compose up -d postgres redis
    
    print_info "等待数据库就绪..."
    sleep 5
    
    # 检查数据库是否就绪
    for i in {1..30}; do
        if docker-compose exec -T postgres pg_isready -U consulting_user -d consulting_system &> /dev/null; then
            print_success "数据库已就绪"
            return 0
        fi
        echo -n "."
        sleep 1
    done
    
    print_error "数据库启动超时"
    return 1
}

# 显示完成信息
show_completion() {
    print_line
    print_success "初始化完成!"
    print_line
    echo ""
    echo "下一步操作:"
    echo ""
    echo "1. 编辑 .env 文件，配置环境变量"
    echo "   vim .env"
    echo ""
    echo "2. 启动所有服务"
    echo "   make up"
    echo "   或"
    echo "   ./scripts/deploy.sh prod"
    echo ""
    echo "3. 访问应用"
    echo "   - 前端: http://localhost"
    echo "   - 后端: http://localhost:8080"
    echo "   - API 文档: http://localhost:8080/swagger/index.html"
    echo ""
    echo "常用命令:"
    echo "   make help          # 查看所有命令"
    echo "   make logs          # 查看日志"
    echo "   make status        # 查看服务状态"
    echo ""
}

# 主函数
main() {
    print_line
    echo -e "${GREEN}  咨询公司业务管理系统 - 初始化脚本${NC}"
    print_line
    echo ""
    
    # 检查环境
    check_docker
    check_go
    check_node
    
    # 创建目录
    create_directories
    
    # 创建环境文件
    create_env_file
    
    # 生成密钥
    generate_secrets
    
    # 询问是否启动基础服务
    echo ""
    read -p "是否立即启动基础服务（PostgreSQL + Redis）? [Y/n] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Nn]$ ]]; then
        start_base_services
    fi
    
    # 显示完成信息
    show_completion
}

# 执行主函数
main "$@"
