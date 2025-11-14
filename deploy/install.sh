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
DEPLOY_SCRIPTS_DIR="."
INSTALL_DIR=""

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

# 安装 python3-venv
install_python_venv() {
    echo -e "${YELLOW}正在安装 python3-venv...${NC}"

    if command -v apt-get &> /dev/null; then
        # Debian/Ubuntu
        if sudo apt-get update && sudo apt-get install -y python3-venv; then
            echo -e "${GREEN}✓ python3-venv 安装成功${NC}"
            return 0
        fi
    elif command -v yum &> /dev/null; then
        # CentOS/RHEL
        if sudo yum install -y python3-venv; then
            echo -e "${GREEN}✓ python3-venv 安装成功${NC}"
            return 0
        fi
    elif command -v pacman &> /dev/null; then
        # Arch Linux
        if sudo pacman -S python-venv; then
            echo -e "${GREEN}✓ python3-venv 安装成功${NC}"
            return 0
        fi
    elif command -v brew &> /dev/null; then
        # macOS
        if brew install python3; then
            echo -e "${GREEN}✓ python3-venv 安装成功${NC}"
            return 0
        fi
    fi

    # 如果自动安装失败，提示用户手动安装
    echo -e "${RED}✗ 无法自动安装 python3-venv（可能没有 sudo 权限）${NC}"
    echo -e "${YELLOW}请手动安装 python3-venv：${NC}"
    echo -e "${YELLOW}  Ubuntu/Debian: sudo apt install python3-venv${NC}"
    echo -e "${YELLOW}  CentOS/RHEL: sudo yum install python3-venv${NC}"
    echo -e "${YELLOW}  Arch Linux: sudo pacman -S python-venv${NC}"
    echo -e "${YELLOW}  macOS: brew install python3${NC}"
    return 1
}

# 安装 Python 依赖
install_dependencies() {
    echo -e "${YELLOW}正在安装 Python 依赖...${NC}"

    # 进入 deploy 目录
    cd "$DEPLOY_SCRIPTS_DIR" || return 1

    if [ ! -d "venv" ]; then
        echo -e "${YELLOW}正在创建虚拟环境...${NC}"
        if ! python3 -m venv venv; then
            echo -e "${RED}✗ 虚拟环境创建失败${NC}"
            install_python_venv || return 1

            # 重新尝试创建虚拟环境
            echo -e "${YELLOW}正在重新创建虚拟环境...${NC}"
            if ! python3 -m venv venv; then
                echo -e "${RED}✗ 虚拟环境创建仍然失败${NC}"
                return 1
            fi
        fi
        echo -e "${GREEN}✓ 虚拟环境已创建${NC}"
    fi

    if [ ! -f "venv/bin/activate" ]; then
        echo -e "${RED}✗ 虚拟环境激活脚本不存在${NC}"
        return 1
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
    python3 install.py "$INSTALL_DIR"
}

# 选择安装目录
prompt_install_dir() {
    echo -e "${YELLOW}请选择安装目录${NC}"
    echo -e "  1. /opt/last-admin (推荐)"
    echo -e "  2. 自定义目录"

    read -p "请选择 [1-2] (默认: 1): " choice
    choice=${choice:-1}

    case $choice in
        1)
            INSTALL_DIR="/opt/last-admin"
            ;;
        2)
            read -p "请输入安装目录 (默认: /opt/last-admin): " custom_dir
            INSTALL_DIR=${custom_dir:-/opt/last-admin}
            ;;
        *)
            echo -e "${RED}✗ 选择无效，使用默认目录${NC}"
            INSTALL_DIR="/opt/last-admin"
            ;;
    esac

    echo -e "${GREEN}✓ 安装目录已设置为: $INSTALL_DIR${NC}"

    # 创建安装目录
    if ! mkdir -p "$INSTALL_DIR"; then
        echo -e "${RED}✗ 无法创建安装目录: $INSTALL_DIR${NC}"
        echo -e "${YELLOW}请确保有足够的权限，或尝试使用 sudo${NC}"
        return 1
    fi
}

# 主函数
main() {
    echo -e "${GREEN}================================${NC}"
    echo -e "${GREEN}Last Admin 安装向导${NC}"
    echo -e "${GREEN}================================${NC}\n"

    detect_platform

    # 提示用户选择安装目录
    if ! prompt_install_dir; then
        exit 1
    fi

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

