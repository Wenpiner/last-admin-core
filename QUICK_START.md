# Last Admin 快速开始指南

## 🚀 最快的方式（推荐）

无需克隆仓库，一行命令快速部署：

### Linux/macOS

```bash
curl -fsSL https://raw.githubusercontent.com/Wenpiner/last-admin/main/last-admin-core/install-bootstrap.sh | bash
```

或使用 wget：

```bash
wget -qO- https://raw.githubusercontent.com/Wenpiner/last-admin/main/last-admin-core/install-bootstrap.sh | bash
```

## 📋 前置要求

- **网络**：能访问 GitHub
- **工具**：curl 或 wget
- **Python**：3.8+ （脚本会自动安装）
- **Docker**：用于容器化部署（可选）

## 🔧 安装步骤

### 步骤 1：运行引导脚本

```bash
curl -fsSL https://raw.githubusercontent.com/Wenpiner/last-admin-core/main/install-bootstrap.sh | bash
```

引导脚本会自动：
- ✅ 获取最新的 Release 版本
- ✅ 从 GitHub Release 下载部署包
- ✅ 解压部署包
- ✅ 运行安装脚本

### 步骤 2：按照向导配置

安装向导会引导你配置：

1. **项目信息**
   - 项目名称
   - 部署环境（开发/测试/生产）

2. **Docker 网络**
   - 选择现有网络或创建新网络

3. **Docker 镜像**
   - API 服务镜像
   - RPC 服务镜像

4. **部署方案**
   - 数据库部署模式（Docker/外部）
   - Redis 部署模式（Docker/外部）

5. **端口配置**
   - API 端口（默认 8889）
   - RPC 端口（默认 8080）
   - 数据库端口（默认 5432）
   - Redis 端口（默认 6379）

6. **数据库配置**
   - 数据库类型（PostgreSQL/MySQL/SQLite3）
   - 用户名、密码、数据库名

7. **Redis 配置**
   - Redis 密码
   - 数据库编号
   - 连接池大小

8. **认证配置**
   - JWT 密钥（自动生成）
   - Token 过期时间
   - OAuth 密钥（自动生成）

9. **验证码配置**
   - 验证码类型
   - 存储方式（内存/Redis）

10. **部署**
    - 生成 docker-compose.yml
    - 拉取 Docker 镜像
    - 启动服务

### 步骤 3：验证部署

```bash
# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 访问 API
curl http://localhost:8889/api/v1/health
```

## 📁 本地部署（已克隆仓库）

```bash
cd last-admin/last-admin-core/deploy
bash install.sh
```

## 🐳 Docker Compose 快速启动

```bash
cd last-admin/last-admin-core/deploy
docker-compose up -d
```

## 🔍 常见问题

### Q: 脚本下载失败怎么办？

A: 检查网络连接，或手动下载脚本：

```bash
curl -fsSL -o install-bootstrap.sh https://raw.githubusercontent.com/Wenpiner/last-admin/main/last-admin-core/install-bootstrap.sh
chmod +x install-bootstrap.sh
./install-bootstrap.sh
```

### Q: Python 版本不兼容？

A: 确保 Python 版本 >= 3.8：

```bash
python3 --version
```

### Q: Docker 镜像拉取失败？

A: 检查 Docker 是否运行，或尝试手动拉取：

```bash
docker pull <image-name>
```

### Q: 如何修改配置？

A: 编辑 `deploy/.env` 文件，然后重启服务：

```bash
docker-compose restart
```

## 📚 更多信息

- 详细部署指南：[DEPLOYMENT.md](./DEPLOYMENT.md)
- 项目文档：[README.md](./README.md)

## 🆘 获取帮助

如遇到问题，请：

1. 查看日志：`docker-compose logs`
2. 检查配置：`cat deploy/.env`
3. 提交 Issue：https://github.com/Wenpiner/last-admin/issues

## ✨ 下一步

部署完成后，你可以：

1. 访问 API：http://localhost:8889
2. 查看 Swagger 文档：http://localhost:8889/swagger
3. 配置管理员账户
4. 开始使用 Last Admin

祝你使用愉快！🎉

