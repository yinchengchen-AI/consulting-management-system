#!/bin/bash
# ============================================
# 咨询公司业务管理系统 - 部署脚本
# ============================================
# 用法: ./scripts/deploy.sh [dev|prod|stop|restart|logs]
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

# 检查环境文件
check_env() {
    if [ ! -f .env ]; then
        print_warning ".env 文件不存在，正在从 .env.example 创建..."
        cp .env.example .env
        print_success ".env 文件已创建，请根据需要修改配置"
    fi
}

# 检查 Docker
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi

    # 检查 Docker 服务是否运行
    if ! docker info &> /dev/null; then
        print_error "Docker 服务未运行，请启动 Docker 服务"
        exit 1
    fi
}

# 开发环境部署
deploy_dev() {
    print_info "启动开发环境..."
    check_env
    check_docker
    
    print_info "启动数据库和缓存服务..."
    docker-compose up -d postgres redis
    
    print_info "等待数据库就绪..."
    sleep 5
    
    print_success "开发环境已启动!"
    echo ""
    echo "服务地址:"
    echo "  - PostgreSQL: localhost:5432"
    echo "  - Redis: localhost:6379"
    echo ""
    echo "请在单独的终端中运行:"
    echo "  make dev-backend   # 启动后端开发服务器"
    echo "  make dev-frontend  # 启动前端开发服务器"
}

# 生产环境部署
deploy_prod() {
    print_info "启动生产环境..."
    check_env
    check_docker
    
    print_info "构建并启动服务..."
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
    
    print_info "等待服务就绪..."
    sleep 10
    
    print_success "生产环境已启动!"
    echo ""
    echo "服务地址:"
    echo "  - 前端应用: http://localhost"
    echo "  - 后端 API: http://localhost:8080"
    echo "  - API 文档: http://localhost:8080/swagger/index.html"
}

# 停止服务
stop_services() {
    print_info "停止所有服务..."
    docker-compose down
    print_success "所有服务已停止"
}

# 重启服务
restart_services() {
    print_info "重启所有服务..."
    docker-compose restart
    print_success "所有服务已重启"
}

# 查看日志
show_logs() {
    if [ -z "$2" ]; then
        docker-compose logs -f --tail=100
    else
        docker-compose logs -f --tail=100 "$2"
    fi
}

# 显示状态
show_status() {
    print_info "服务状态:"
    docker-compose ps
}

# 备份数据
backup_data() {
    print_info "备份数据库..."
    mkdir -p backups
    
    BACKUP_FILE="backups/consulting_system_$(date +%Y%m%d_%H%M%S).sql"
    docker-compose exec -T postgres pg_dump -U consulting_user consulting_system > "$BACKUP_FILE"
    
    print_success "数据库备份完成: $BACKUP_FILE"
}

# 显示帮助
show_help() {
    echo "咨询公司业务管理系统 - 部署脚本"
    echo ""
    echo "用法: $0 [命令] [选项]"
    echo ""
    echo "命令:"
    echo "  dev       启动开发环境"
    echo "  prod      启动生产环境"
    echo "  stop      停止所有服务"
    echo "  restart   重启所有服务"
    echo "  logs      查看日志 [service_name]"
    echo "  status    查看服务状态"
    echo "  backup    备份数据库"
    echo "  help      显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 dev              # 启动开发环境"
    echo "  $0 prod             # 启动生产环境"
    echo "  $0 logs             # 查看所有日志"
    echo "  $0 logs backend     # 查看后端日志"
    echo "  $0 backup           # 备份数据库"
}

# 主函数
main() {
    case "${1:-help}" in
        dev)
            deploy_dev
            ;;
        prod)
            deploy_prod
            ;;
        stop)
            stop_services
            ;;
        restart)
            restart_services
            ;;
        logs)
            show_logs "$@"
            ;;
        status)
            show_status
            ;;
        backup)
            backup_data
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
