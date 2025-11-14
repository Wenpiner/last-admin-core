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

# GitHub 镜像地址列表（前缀）
# 空字符串表示直接使用 github.com
GITHUB_MIRRORS=(
    ""
    "https://fastgit.cc/"
    "https://gh.ddlc.top/"
    "https://github.boki.moe/"
    "https://ghp.keleyaa.com/"
    "https://ghproxy.cfd/"
    "https://xget.xi-xu.me/gh/"
)

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}Last Admin 快速安装${NC}"
echo -e "${GREEN}================================${NC}\n"

# 检查必要的命令
if ! command -v curl &> /dev/null; then
    echo -e "${RED}✗ 未检测到 curl，请先安装 curl${NC}"
    exit 1
fi

if ! command -v sha256sum &> /dev/null && ! command -v shasum &> /dev/null; then
    echo -e "${RED}✗ 未检测到 sha256sum 或 shasum，请先安装${NC}"
    exit 1
fi

# 创建部署目录
echo -e "${YELLOW}正在初始化部署目录: $DEPLOY_HOME${NC}"
mkdir -p "$DEPLOY_HOME"
echo -e "${GREEN}✓ 部署目录已创建${NC}"

# 获取最新的 Release 版本和 SHA256
echo -e "${YELLOW}正在获取最新的 Release 版本...${NC}"

RELEASE_INFO=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest")
LATEST_RELEASE=$(echo "$RELEASE_INFO" | grep '"tag_name"' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo -e "${RED}✗ 无法获取最新的 Release 版本${NC}"
    exit 1
fi

VERSION_NUM=${LATEST_RELEASE#v}
echo -e "${GREEN}✓ 最新版本: $LATEST_RELEASE${NC}"

# 从 Release assets 中提取 SHA256 文件 URL
# 查找 deploy.tar.gz.sha256 的下载 URL
SHA256_FILE_URL=$(echo "$RELEASE_INFO" | grep -o '"browser_download_url": "[^"]*deploy\.tar\.gz\.sha256"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$SHA256_FILE_URL" ]; then
    echo -e "${YELLOW}⚠ 未找到 SHA256 文件，将跳过校验${NC}"
    SKIP_SHA256_CHECK=true
else
    # 下载 SHA256 文件
    EXPECTED_SHA256=$(curl -s "$SHA256_FILE_URL")
    if [ -z "$EXPECTED_SHA256" ]; then
        echo -e "${YELLOW}⚠ 无法读取 SHA256 值，将跳过校验${NC}"
        SKIP_SHA256_CHECK=true
    else
        echo -e "${GREEN}✓ SHA256: ${EXPECTED_SHA256:0:16}...${NC}"
        SKIP_SHA256_CHECK=false
    fi
fi

# 下载部署包（带重试机制）
echo -e "${YELLOW}正在下载部署包...${NC}"

# 构建基础 URL
BASE_URL="https://github.com/${GITHUB_REPO}/releases/download/${LATEST_RELEASE}/deploy.tar.gz"

DOWNLOAD_SUCCESS=false
for mirror in "${GITHUB_MIRRORS[@]}"; do
    # 构建完整的下载 URL
    DOWNLOAD_URL="${mirror}${BASE_URL}"
    MIRROR_NAME="$mirror"

    echo -e "${YELLOW}尝试从 $MIRROR_NAME 下载...${NC}"

    if curl -fsSL -o "${DEPLOY_HOME}/deploy.tar.gz" "$DOWNLOAD_URL" 2>/dev/null; then
        # 验证 SHA256
        if [ "$SKIP_SHA256_CHECK" = false ]; then
            if command -v sha256sum &> /dev/null; then
                ACTUAL_SHA256=$(sha256sum "${DEPLOY_HOME}/deploy.tar.gz" | awk '{print $1}')
            else
                ACTUAL_SHA256=$(shasum -a 256 "${DEPLOY_HOME}/deploy.tar.gz" | awk '{print $1}')
            fi

            if [ "$ACTUAL_SHA256" = "$EXPECTED_SHA256" ]; then
                echo -e "${GREEN}✓ SHA256 校验通过${NC}"
                DOWNLOAD_SUCCESS=true
                break
            else
                echo -e "${YELLOW}⚠ SHA256 校验失败，尝试下一个镜像${NC}"
                rm -f "${DEPLOY_HOME}/deploy.tar.gz"
                continue
            fi
        else
            echo -e "${GREEN}✓ 部署包下载完成${NC}"
            DOWNLOAD_SUCCESS=true
            break
        fi
    fi
done

if [ "$DOWNLOAD_SUCCESS" = false ]; then
    echo -e "${RED}✗ 从所有镜像下载部署包都失败了${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 部署包下载完成${NC}"

# 解压部署包
echo -e "${YELLOW}正在解压部署包...${NC}"

if ! tar -xzf "${DEPLOY_HOME}/deploy.tar.gz" -C "$DEPLOY_HOME"; then
    echo -e "${RED}✗ 解压失败${NC}"
    exit 1
fi

if [ ! -d "${DEPLOY_HOME}/deploy" ]; then
    echo -e "${RED}✗ 解压后未找到 deploy 目录${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 解压完成${NC}"

# 清理旧的 tar.gz 文件
rm -f "${DEPLOY_HOME}/deploy.tar.gz"

# 进入 deploy 目录并运行安装脚本
cd "${DEPLOY_HOME}/deploy"
chmod +x install.sh
exec ./install.sh

