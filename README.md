# Last Admin Core

一个为中后台管理系统设计的现代化后端解决方案，基于 Go Zero 框架构建，提供完整的 RBAC（基于角色的访问控制）能力。

## 项目简介

Last Admin Core 是 Last Admin 开源方案的后端核心组件。它提供了一套可扩展的 API 和 RPC 服务基础设施，用于构建企业级管理系统，具备权限控制和多角色支持能力。

### 核心特性

- **RBAC 实现**：基于 [Casbin](https://casbin.org/) 的灵活角色权限控制
- **多角色支持**：
  - 多角色菜单管理
  - 多角色配置项管理
  - 细粒度权限控制
- **双服务架构**：
  - REST API 服务（面向客户端应用）
  - gRPC RPC 服务（内部服务通信）
- **企业级可靠性**：基于成熟技术栈，生产级别的稳定性

## 技术栈

### 核心框架
- **Go Zero**：高性能微服务框架
- **Go 1.24.2**

### 主要依赖
- **gRPC & Protocol Buffers**：高效的 RPC 通信
- **Ent**：数据库实体框架
- **Casbin**：权限管理库（RBAC 实现）
- **Redis**：缓存和分布式锁
- **PostgreSQL**：主数据库（支持 MySQL 和 SQLite3）
- **JWT**：令牌认证
- **OAuth2**：第三方认证支持

### 其他库
- **OTP**：一次性密码支持（双因素认证）
- **UUID**：唯一标识符生成
- **Validator**：输入验证和本地化支持

## 项目结构

```
last-admin-core/
├── api/                    # REST API 服务
│   ├── core.go            # API 入口
│   ├── desc/              # API 定义
│   ├── internal/          # 内部实现
│   ├── etc/               # 配置文件
│   └── swagger.json       # API 文档
├── rpc/                    # gRPC RPC 服务
│   ├── core.go            # RPC 入口
│   ├── core.proto         # Protocol Buffer 定义
│   ├── internal/          # 内部实现
│   ├── types/             # 生成的类型
│   ├── ent/               # 实体定义
│   └── etc/               # 配置文件
├── deploy/                # 部署脚本
│   ├── install.sh         # Linux/macOS 部署
│   ├── install.ps1        # Windows 部署
│   ├── docker-compose.yml # Docker Compose 配置
│   └── templates/         # 配置模板
├── scripts/               # 工具脚本
│   └── migrate/           # 数据库迁移
├── Dockerfile.api         # API 服务 Docker 镜像
├── Dockerfile.rpc         # RPC 服务 Docker 镜像
├── Makefile               # 构建和代码生成命令
└── go.mod                 # Go 模块定义
```

## 服务说明

### REST API 服务
为客户端应用提供 HTTP REST 接口，处理：
- 用户认证和授权
- 资源管理
- 配置管理
- 菜单和角色管理

**入口文件**：`api/core.go`

### gRPC RPC 服务
为内部服务通信提供高性能 gRPC 接口，包含：
- API 服务
- 字典服务
- 角色服务
- 菜单服务
- 部门服务
- 职位服务
- 用户服务
- OAuth 提供商服务
- 初始化服务
- 令牌服务
- 配置服务

**入口文件**：`rpc/core.go`

## 快速开始

### 前置要求
- Go 1.24.2 或更高版本
- PostgreSQL 数据库
- Redis 

### 部署

项目提供了自动化部署脚本，支持多种部署方式：

#### 方式 1：快速安装（推荐）

无需克隆整个仓库，直接下载并运行安装脚本：

**Linux/macOS：**
```bash
curl -fsSL https://raw.githubusercontent.com/Wenpiner/last-admin/main/last-admin-core/install-bootstrap.sh | bash
```

或使用 wget：
```bash
wget -qO- https://raw.githubusercontent.com/Wenpiner/last-admin/main/last-admin-core/install-bootstrap.sh | bash
```

#### 方式 2：本地部署

如果已克隆仓库，可以使用本地脚本：

**Linux/macOS：**
```bash
cd deploy
bash install.sh
```

**Windows：**
```powershell
cd deploy
.\install.ps1
```

#### 方式 3：Docker Compose

使用 Docker Compose 快速启动所有服务：
```bash
cd deploy
docker-compose up -d
```

### 部署文档

- **[快速开始](./QUICK_START.md)** - 一行命令快速部署指南
- **[详细部署指南](./DEPLOYMENT.md)** - 完整的部署配置和故障排除
- **[部署脚本](./deploy/)** - 部署脚本和配置文件

## 开发指南

### 代码生成

生成 RPC 服务代码：
```bash
make gen-rpc
```

生成 API 服务代码：
```bash
make gen-api
```

生成 Swagger 文档：
```bash
make swagger
```

生成 Ent 实体代码：
```bash
make ent-generate
```

执行数据库迁移：
```bash
make ent-migrate
```

### 配置文件

API 和 RPC 服务均使用 YAML 配置文件，位置如下：
- API：`api/etc/core.yaml`
- RPC：`rpc/etc/core.yaml`

## API 文档

API 文档通过 Swagger UI 提供。部署完成后，可通过配置的 API 端点访问 Swagger 界面。

本地生成 Swagger 文档：
```bash
make swagger
make swagger-serve
```

Swagger UI 将在 `http://localhost:36666` 可访问

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

### MIT 许可证说明

MIT 许可证是一个宽松的开源许可证，允许你：
- ✅ 将软件用于任何目的
- ✅ 复制、修改和分发软件
- ✅ 在私有或商业项目中使用软件

唯一的要求是：
- 📋 包含许可证副本和版权声明

更多信息，请访问 [opensource.org/licenses/MIT](https://opensource.org/licenses/MIT)

## 支持

如有问题或建议，请参考项目仓库。

---

**Last Admin Core** - 用 Go 构建企业级管理系统

