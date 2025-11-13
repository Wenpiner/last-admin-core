#!/bin/bash
# install.sh - macOS/Linux 通用安装脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 平台检测
detect_platform() {
    case "$(uname -s)" in
        Linux*)     PLATFORM="Linux";;
        Darwin*)    PLATFORM="macOS";;
        *)          PLATFORM="UNKNOWN";;
    esac
    echo -e "${GREEN}检测到平台: $PLATFORM${NC}"
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
    
    source venv/bin/activate
    python3 install.py
}

# 主函数
main() {
    echo -e "${GREEN}================================${NC}"
    echo -e "${GREEN}Last Admin 安装向导${NC}"
    echo -e "${GREEN}================================${NC}\n"
    
    detect_platform
    
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

