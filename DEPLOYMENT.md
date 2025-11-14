# Last Admin 部署指南

本文档详细说明了 Last Admin 的部署方式和流程。

## 部署方式概览

Last Admin 提供了多种部署方式，适应不同的使用场景：

| 方式 | 适用场景 | 优点 | 缺点 |
|------|--------|------|------|
| 快速安装 | 快速部署，无需克隆仓库 | 简单快速，自动下载最新脚本 | 需要网络连接 |
| 本地部署 | 已克隆仓库，本地开发 | 灵活，可修改脚本 | 需要克隆整个仓库 |
| Docker Compose | 容器化部署 | 环境隔离，易于扩展 | 需要 Docker 环境 |

## 快速安装（推荐）

### 前置要求

- Linux 或 macOS 系统
- curl 或 wget 命令行工具
- Python 3.8+
- Docker 和 Docker Compose（用于容器化部署）

### 安装步骤

#### 使用 curl（推荐）

```bash
curl -fsSL https://raw.githubusercontent.com/Wenpiner/last-admin-core/main/install-bootstrap.sh | bash
```

#### 使用 wget

```bash
wget -qO- https://raw.githubusercontent.com/Wenpiner/last-admin-core/main/install-bootstrap.sh | bash
```

#### 手动下载并运行

```bash
# 下载脚本
curl -fsSL -o install-bootstrap.sh https://raw.githubusercontent.com/Wenpiner/last-admin-core/main/install-bootstrap.sh

# 赋予执行权限
chmod +x install-bootstrap.sh

# 运行脚本
./install-bootstrap.sh
```

### 工作流程

1. **获取最新版本**：从 GitHub API 获取最新的 Release 版本
2. **下载部署包**：从 GitHub Release 下载 `deploy-scripts-*.tar.gz`
3. **解压部署包**：将脚本包解压到临时目录
4. **运行安装脚本**：执行解压后的 `install.sh`
5. **检测平台**：自动检测 Linux 或 macOS
6. **检查 Python**：验证 Python 3 是否已安装，如未安装则自动安装
7. **安装依赖**：创建虚拟环境并安装 Python 依赖
8. **启动向导**：运行交互式安装向导

## 本地部署

### 前置要求

- 已克隆 Last Admin 仓库
- Python 3.8+
- Docker 和 Docker Compose

### 安装步骤

```bash
# 进入 deploy 目录
cd last-admin-core/deploy

# 运行安装脚本
bash install.sh
```

### 使用 Makefile

```bash
# 进入 deploy 目录
cd last-admin-core/deploy

# 查看可用命令
make help

# 运行安装
make install

# 运行测试
make test

# 清理临时文件
make clean
```

## Docker Compose 部署

### 前置要求

- Docker 和 Docker Compose 已安装
- 足够的磁盘空间

### 快速启动

```bash
cd last-admin-core/deploy

# 启动所有服务
docker compose up -d

# 查看日志
docker compose logs -f

# 停止服务
docker compose down
```

### 服务说明

Docker Compose 配置包含以下服务：

- **PostgreSQL**：数据库服务
- **Redis**：缓存服务
- **API 服务**：REST API 服务（端口 8889）
- **RPC 服务**：gRPC 服务（端口 8080）

## 安装向导配置

运行安装脚本后，会进入交互式安装向导，需要配置以下内容：

### 第 1 阶段：项目基本信息

- **项目名称**：用于标识部署实例
- **部署环境**：开发、测试或生产

### 第 2 阶段：Docker 网络配置

- 选择使用现有 Docker 网络或创建新网络

### 第 3 阶段：Docker 镜像配置

- **API 镜像**：指定 API 服务镜像和标签
- **RPC 镜像**：指定 RPC 服务镜像和标签

### 第 4 阶段：组件部署方案

- **数据库部署模式**：Docker 或外部
- **Redis 部署模式**：Docker 或外部

### 第 5 阶段：端口配置

- **API 端口**：默认 8889
- **RPC 端口**：默认 8080
- **数据库端口**：默认 5432
- **Redis 端口**：默认 6379

### 第 6 阶段：数据库配置

- **数据库类型**：PostgreSQL、MySQL 或 SQLite3
- **数据库用户**：默认 postgres
- **数据库密码**：默认 postgres123
- **数据库名称**：自动生成或自定义
- **SSL 模式**：disable、require 或 prefer

### 第 7 阶段：Redis 配置

- **Redis 密码**：默认 redis123
- **Redis 数据库编号**：默认 0
- **连接池大小**：默认 10

### 第 8 阶段：认证配置

- **JWT 密钥**：自动生成
- **Token 过期时间**：默认 360000 秒
- **OAuth 密钥**：自动生成

### 第 9 阶段：验证码配置

- **验证码类型**：选择验证码实现方式
- **存储类型**：内存或 Redis

### 第 10 阶段：部署服务

- 自动生成 docker-compose.yml
- 拉取 Docker 镜像
- 启动服务

## 配置文件

安装完成后，会在 deploy 目录生成以下文件：

- `.env`：环境变量配置文件
- `docker-compose.yml`：Docker Compose 配置文件

## 故障排除

### 问题：无法下载脚本

**解决方案**：
- 检查网络连接
- 尝试使用代理或 VPN
- 手动下载脚本后运行

### 问题：Python 版本不兼容

**解决方案**：
- 确保 Python 版本 >= 3.8
- 使用 `python3 --version` 检查版本

### 问题：Docker 镜像拉取失败

**解决方案**：
- 检查 Docker 是否正常运行
- 检查网络连接和 Docker 镜像仓库访问权限
- 尝试手动拉取镜像：`docker pull <image-name>`

## 更新部署脚本

当 deploy 目录有更新时，GitHub Actions 会自动打包新的脚本包并发布到 Release。

### 手动触发打包

访问 GitHub Actions 工作流页面，手动触发 "Package Deploy Scripts" 工作流。

### 获取最新脚本

快速安装方式会自动下载最新的脚本包，无需手动更新。

## 安全建议

1. **修改默认密码**：安装后立即修改数据库和 Redis 密码
2. **使用 SSL/TLS**：在生产环境中启用 SSL 连接
3. **网络隔离**：将数据库和 Redis 放在私有网络中
4. **定期备份**：定期备份数据库和配置文件
5. **监控日志**：定期检查应用和系统日志

## 相关文档

- [README.md](./README.md) - 项目概览
- [deploy/](./deploy/) - 部署脚本目录

