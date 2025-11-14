#!/bin/bash
# install-bootstrap.sh - Last Admin 快速安装脚本
# 用法: curl -fsSL https://raw.githubusercontent.com/Wenpiner/last-admin-core/main/install-bootstrap.sh | bash
# 或: wget -qO- https://raw.githubusercontent.com/Wenpiner/last-admin-core/main/install-bootstrap.sh | bash

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 配置
GITHUB_REPO="Wenpiner/last-admin-core"
DEPLOY_HOME="${HOME}/.last-admin/deploy"

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}Last Admin 快速安装${NC}"
echo -e "${GREEN}================================${NC}\n"

# 检查 curl 是否已安装
if ! command -v curl &> /dev/null; then
    echo -e "${RED}✗ 未检测到 curl，请先安装 curl${NC}"
    exit 1
fi

# 创建部署目录
echo -e "${YELLOW}正在初始化部署目录: $DEPLOY_HOME${NC}"
mkdir -p "$DEPLOY_HOME"
echo -e "${GREEN}✓ 部署目录已创建${NC}"

# 获取最新的 Release 版本
echo -e "${YELLOW}正在获取最新的 Release 版本...${NC}"

LATEST_RELEASE=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name"' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo -e "${RED}✗ 无法获取最新的 Release 版本${NC}"
    exit 1
fi

VERSION_NUM=${LATEST_RELEASE#v}
echo -e "${GREEN}✓ 最新版本: $LATEST_RELEASE${NC}"

# 下载部署包
echo -e "${YELLOW}正在下载部署包...${NC}"

DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${LATEST_RELEASE}/deploy-scripts-${VERSION_NUM}.tar.gz"

if ! curl -fsSL -o "${DEPLOY_HOME}/deploy-scripts.tar.gz" "$DOWNLOAD_URL"; then
    echo -e "${RED}✗ 下载部署包失败${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 部署包下载完成${NC}"

# 解压部署包
echo -e "${YELLOW}正在解压部署包...${NC}"

if ! tar -xzf "${DEPLOY_HOME}/deploy-scripts.tar.gz" -C "$DEPLOY_HOME"; then
    echo -e "${RED}✗ 解压失败${NC}"
    exit 1
fi

if [ ! -d "${DEPLOY_HOME}/deploy" ]; then
    echo -e "${RED}✗ 解压后未找到 deploy 目录${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 解压完成${NC}"

# 清理旧的 tar.gz 文件
rm -f "${DEPLOY_HOME}/deploy-scripts.tar.gz"

# 进入 deploy 目录并运行安装脚本
cd "${DEPLOY_HOME}/deploy"
chmod +x install.sh
exec ./install.sh

