#!/bin/bash

# 🧹 Podman/Docker 清理脚本
# 用于管理和清理悬挂镜像

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# 显示当前镜像状态
show_images() {
    print_message $BLUE "📊 当前镜像状态:"
    podman images
    
    local dangling_count=$(podman images -f "dangling=true" -q | wc -l)
    print_message $YELLOW "💡 悬挂镜像数量: $dangling_count"
}

# 清理悬挂镜像
clean_dangling() {
    print_message $BLUE "🧹 清理悬挂镜像..."
    
    local dangling_images=$(podman images -f "dangling=true" -q)
    
    if [ -z "$dangling_images" ]; then
        print_message $GREEN "✅ 没有发现悬挂镜像"
    else
        local count=$(echo "$dangling_images" | wc -l)
        print_message $YELLOW "🔍 发现 $count 个悬挂镜像，正在清理..."
        podman image prune -f
        print_message $GREEN "✅ 悬挂镜像清理完成"
    fi
}

# 清理未使用的镜像
clean_unused() {
    print_message $BLUE "🧹 清理所有未使用的镜像..."
    print_message $YELLOW "⚠️  警告: 这将删除所有未被容器使用的镜像！"
    
    read -p "是否继续? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        podman image prune -a -f
        print_message $GREEN "✅ 未使用镜像清理完成"
    else
        print_message $YELLOW "❌ 操作已取消"
    fi
}

# 系统级清理
system_clean() {
    print_message $BLUE "🧹 系统级清理..."
    print_message $RED "⚠️  危险: 这将清理所有未使用的容器、镜像、网络和卷！"
    
    read -p "确定要继续吗? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        podman system prune -a -f --volumes
        print_message $GREEN "✅ 系统清理完成"
    else
        print_message $YELLOW "❌ 操作已取消"
    fi
}

# 智能清理 - 保留最新的几个版本
smart_clean() {
    print_message $BLUE "🧠 智能清理 - 保留最新3个版本的镜像..."
    
    # 获取项目相关的镜像，按创建时间排序
    local project_images=$(podman images --format "table {{.Repository}}:{{.Tag}} {{.ID}} {{.Created}}" \
        | grep "test_podman_web-server" \
        | grep -v "REPOSITORY" \
        | sort -k3 -r)
    
    if [ -z "$project_images" ]; then
        print_message $YELLOW "没有找到项目相关的镜像"
        return
    fi
    
    local count=0
    local to_delete=()
    
    while IFS= read -r line; do
        count=$((count + 1))
        if [ $count -gt 3 ]; then
            local image_id=$(echo "$line" | awk '{print $2}')
            to_delete+=("$image_id")
        fi
    done <<< "$project_images"
    
    if [ ${#to_delete[@]} -eq 0 ]; then
        print_message $GREEN "✅ 只有3个或更少的镜像版本，无需清理"
    else
        print_message $YELLOW "🔍 将删除 ${#to_delete[@]} 个旧版本镜像"
        for image_id in "${to_delete[@]}"; do
            podman rmi "$image_id" 2>/dev/null || true
        done
        print_message $GREEN "✅ 智能清理完成"
    fi
    
    # 清理悬挂镜像
    clean_dangling
}

# 显示帮助信息
show_help() {
    echo "🧹 Podman/Docker 清理工具"
    echo
    echo "使用方法: $0 [选项]"
    echo
    echo "选项:"
    echo "  -s, --show      显示当前镜像状态"
    echo "  -d, --dangling  清理悬挂镜像 (<none> 镜像)"
    echo "  -u, --unused    清理所有未使用的镜像"
    echo "  -a, --all       系统级清理 (危险)"
    echo "  -i, --smart     智能清理 (推荐)"
    echo "  -h, --help      显示此帮助信息"
    echo
    echo "示例:"
    echo "  $0 -d           # 清理悬挂镜像"
    echo "  $0 --smart      # 智能清理"
    echo "  $0 --show       # 显示镜像状态"
}

# 主函数
main() {
    case "${1:-}" in
        -s|--show)
            show_images
            ;;
        -d|--dangling)
            clean_dangling
            ;;
        -u|--unused)
            clean_unused
            ;;
        -a|--all)
            system_clean
            ;;
        -i|--smart)
            smart_clean
            ;;
        -h|--help|*)
            show_help
            ;;
    esac
}

# 执行主函数
main "$@"
