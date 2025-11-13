#!/bin/bash
# install-bootstrap.sh - Last Admin 快速安装脚本
# 用法: curl -fsSL https://raw.githubusercontent.com/Wenpiner/last-admin/main/last-admin-core/install-bootstrap.sh | bash
# 或: wget -qO- https://raw.githubusercontent.com/Wenpiner/last-admin/main/last-admin-core/install-bootstrap.sh | bash

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 配置
GITHUB_REPO="Wenpiner/last-admin"
GITHUB_RAW_URL="https://raw.githubusercontent.com/${GITHUB_REPO}/main/last-admin-core"
TEMP_DIR=$(mktemp -d)

# 清理临时文件
cleanup() {
    if [ -d "$TEMP_DIR" ]; then
        rm -rf "$TEMP_DIR"
    fi
}

trap cleanup EXIT

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}Last Admin 快速安装${NC}"
echo -e "${GREEN}================================${NC}\n"

# 检查 curl 是否已安装
if ! command -v curl &> /dev/null; then
    echo -e "${RED}✗ 未检测到 curl，请先安装 curl${NC}"
    exit 1
fi

echo -e "${YELLOW}正在下载安装脚本...${NC}"

# 下载 install.sh 到临时目录
if ! curl -fsSL -o "${TEMP_DIR}/install.sh" "${GITHUB_RAW_URL}/deploy/install.sh"; then
    echo -e "${RED}✗ 下载安装脚本失败${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 安装脚本下载完成${NC}"

# 赋予执行权限并运行
chmod +x "${TEMP_DIR}/install.sh"
exec "${TEMP_DIR}/install.sh"

