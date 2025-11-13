#!/bin/bash
# install.sh - macOS/Linux 通用安装脚本
# 支持从 GitHub Release 下载 deploy 脚本包或使用本地脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 配置
GITHUB_REPO="Wenpiner/last-admin"
DEPLOY_SCRIPTS_DIR="deploy"
TEMP_DIR=$(mktemp -d)

# 清理临时文件
cleanup() {
    if [ -d "$TEMP_DIR" ]; then
        rm -rf "$TEMP_DIR"
    fi
}

trap cleanup EXIT

# 平台检测
detect_platform() {
    case "$(uname -s)" in
        Linux*)     PLATFORM="Linux";;
        Darwin*)    PLATFORM="macOS";;
        *)          PLATFORM="UNKNOWN";;
    esac
    echo -e "${GREEN}检测到平台: $PLATFORM${NC}"
}

# 获取最新的 Release 版本
get_latest_release() {
    echo -e "${YELLOW}正在获取最新的 Release 版本...${NC}"

    LATEST_RELEASE=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name"' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$LATEST_RELEASE" ]; then
        echo -e "${RED}✗ 无法获取最新的 Release 版本${NC}"
        return 1
    fi

    echo -e "${GREEN}✓ 最新版本: $LATEST_RELEASE${NC}"
    echo "$LATEST_RELEASE"
}

# 下载 deploy 脚本包
download_deploy_scripts() {
    local version=$1
    local version_num=${version#v}
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/deploy-scripts-${version_num}.tar.gz"

    echo -e "${YELLOW}正在从 GitHub 下载 deploy 脚本包...${NC}"
    echo -e "${YELLOW}下载地址: $download_url${NC}"

    if ! curl -fsSL -o "${TEMP_DIR}/deploy-scripts.tar.gz" "$download_url"; then
        echo -e "${RED}✗ 下载失败${NC}"
        return 1
    fi

    echo -e "${GREEN}✓ 下载完成${NC}"
}

# 解压 deploy 脚本包
extract_deploy_scripts() {
    echo -e "${YELLOW}正在解压 deploy 脚本包...${NC}"

    if ! tar -xzf "${TEMP_DIR}/deploy-scripts.tar.gz" -C "$TEMP_DIR"; then
        echo -e "${RED}✗ 解压失败${NC}"
        return 1
    fi

    # 将解压的 deploy 目录复制到当前目录
    if [ -d "${TEMP_DIR}/deploy" ]; then
        cp -r "${TEMP_DIR}/deploy" .
        echo -e "${GREEN}✓ 解压完成${NC}"
        return 0
    else
        echo -e "${RED}✗ 解压后未找到 deploy 目录${NC}"
        return 1
    fi
}

# 检查是否需要下载脚本
check_and_download_scripts() {
    if [ ! -d "$DEPLOY_SCRIPTS_DIR" ]; then
        echo -e "${YELLOW}未检测到本地 deploy 目录，将从 GitHub 下载...${NC}"

        local latest_version
        latest_version=$(get_latest_release) || return 1

        download_deploy_scripts "$latest_version" || return 1
        extract_deploy_scripts || return 1
    else
        echo -e "${GREEN}✓ 检测到本地 deploy 目录${NC}"
    fi
}

# 检查 Python 是否已安装
check_python() {
    if command -v python3 &> /dev/null; then
        PYTHON_VERSION=$(python3 --version 2>&1 | awk '{print $2}')
        echo -e "${GREEN}✓ Python 3 已安装: $PYTHON_VERSION${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠ Python 3 未安装${NC}"
        return 1
    fi
}

# macOS 安装 Python
install_python_macos() {
    echo -e "${YELLOW}正在为 macOS 安装 Python 3...${NC}"
    
    if ! command -v brew &> /dev/null; then
        echo -e "${YELLOW}正在安装 Homebrew...${NC}"
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    fi
    
    brew install python3
    echo -e "${GREEN}✓ Python 3 安装完成${NC}"
}

# Linux 安装 Python
install_python_linux() {
    echo -e "${YELLOW}正在为 Linux 安装 Python 3...${NC}"
    
    if command -v apt-get &> /dev/null; then
        sudo apt-get update
        sudo apt-get install -y python3 python3-pip python3-venv
    elif command -v yum &> /dev/null; then
        sudo yum install -y python3 python3-pip
    elif command -v pacman &> /dev/null; then
        sudo pacman -S python python-pip
    else
        echo -e "${RED}✗ 无法识别的包管理器${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✓ Python 3 安装完成${NC}"
}

# 安装 Python 依赖
install_dependencies() {
    echo -e "${YELLOW}正在安装 Python 依赖...${NC}"

    # 进入 deploy 目录
    cd "$DEPLOY_SCRIPTS_DIR" || return 1

    if [ ! -d "venv" ]; then
        python3 -m venv venv
        echo -e "${GREEN}✓ 虚拟环境已创建${NC}"
    fi

    source venv/bin/activate

    pip install --upgrade pip
    pip install -r requirements.txt

    echo -e "${GREEN}✓ Python 依赖安装完成${NC}"
}

# 运行主安装脚本
run_installer() {
    echo -e "${YELLOW}正在启动安装向导...${NC}"

    # 确保在 deploy 目录中
    if [ ! -f "install.py" ]; then
        cd "$DEPLOY_SCRIPTS_DIR" || return 1
    fi

    source venv/bin/activate
    python3 install.py
}

# 主函数
main() {
    echo -e "${GREEN}================================${NC}"
    echo -e "${GREEN}Last Admin 安装向导${NC}"
    echo -e "${GREEN}================================${NC}\n"

    detect_platform

    # 检查并下载 deploy 脚本（如果需要）
    check_and_download_scripts || exit 1

    if ! check_python; then
        case "$PLATFORM" in
            macOS)
                install_python_macos
                ;;
            Linux)
                install_python_linux
                ;;
            *)
                echo -e "${RED}✗ 不支持的平台${NC}"
                exit 1
                ;;
        esac
    fi

    install_dependencies
    run_installer
}

trap 'echo -e "${RED}✗ 安装失败${NC}"; exit 1' ERR

main

