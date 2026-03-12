#!/bin/bash
# ============================================
# 咨询公司业务管理系统 - 健康检查脚本
# ============================================
# 用法: ./scripts/healthcheck.sh
# ============================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# 检查 HTTP 服务
check_http() {
    local url=$1
    local name=$2
    
    if curl -sf "$url" &> /dev/null; then
        print_success "$name 运行正常"
        return 0
    else
        print_error "$name 无法访问"
        return 1
    fi
}

# 检查 Docker 容器
check_container() {
    local container=$1
    
    if docker ps --format "{{.Names}}" | grep -q "^${container}$"; then
        local status=$(docker inspect --format='{{.State.Status}}' "$container" 2>/dev/null)
        local health=$(docker inspect --format='{{.State.Health.Status}}' "$container" 2>/dev/null || echo "none")
        
        if [ "$status" = "running" ]; then
            if [ "$health" = "none" ] || [ "$health" = "healthy" ]; then
                print_success "$container 运行正常 (状态: $status, 健康: $health)"
                return 0
            else
                print_warning "$container 健康检查异常 (健康: $health)"
                return 1
            fi
        else
            print_error "$container 未运行 (状态: $status)"
            return 1
        fi
    else
        print_error "$container 容器不存在"
        return 1
    fi
}

# 检查数据库连接
check_database() {
    print_info "检查数据库连接..."
    
    if docker ps --format "{{.Names}}" | grep -q "^consulting-postgres$"; then
        if docker-compose exec -T postgres pg_isready -U consulting_user -d consulting_system &> /dev/null; then
            print_success "PostgreSQL 连接正常"
            return 0
        else
            print_error "PostgreSQL 连接失败"
            return 1
        fi
    else
        print_error "PostgreSQL 容器未运行"
        return 1
    fi
}

# 检查 Redis 连接
check_redis() {
    print_info "检查 Redis 连接..."
    
    if docker ps --format "{{.Names}}" | grep -q "^consulting-redis$"; then
        if docker-compose exec -T redis redis-cli ping &> /dev/null | grep -q "PONG"; then
            print_success "Redis 连接正常"
            return 0
        else
            print_error "Redis 连接失败"
            return 1
        fi
    else
        print_error "Redis 容器未运行"
        return 1
    fi
}

# 主检查函数
main() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  系统健康检查${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
    
    local failed=0
    
    # 检查容器状态
    print_info "检查容器状态..."
    check_container "consulting-postgres" || ((failed++))
    check_container "consulting-redis" || ((failed++))
    check_container "consulting-backend" || ((failed++))
    check_container "consulting-frontend" || ((failed++))
    
    echo ""
    
    # 检查服务连接
    check_database || ((failed++))
    check_redis || ((failed++))
    
    echo ""
    
    # 检查 HTTP 服务
    print_info "检查 HTTP 服务..."
    check_http "http://localhost:8080/health" "后端 API" || ((failed++))
    check_http "http://localhost/health" "前端应用" || ((failed++))
    
    echo ""
    echo -e "${BLUE}========================================${NC}"
    
    if [ $failed -eq 0 ]; then
        print_success "所有检查通过! 系统运行正常"
        exit 0
    else
        print_error "发现 $failed 个问题，请检查服务状态"
        echo ""
        echo "查看日志:"
        echo "  make logs"
        echo "  make logs-backend"
        echo "  make logs-frontend"
        exit 1
    fi
}

# 执行主函数
main "$@"
