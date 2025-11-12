# 数据库迁移工具实现说明

## 概述

本实现为 last-admin 项目添加了一个交互式的数据库迁移工具，用于初始化和同步 Ent 数据库模式。

## 实现的功能

### 1. 交互式数据库类型选择
- 支持 PostgreSQL（推荐）、MySQL、SQLite3 三种数据库
- 清晰的菜单界面
- 默认选择 PostgreSQL

### 2. 灵活的配置输入
- 每个配置项都有合理的默认值
- 用户可以按 Enter 使用默认值，或输入自定义值
- 支持远程数据库连接（通过输入主机地址）

### 3. 数据库特定配置
- **PostgreSQL**: 主机、端口、用户名、密码、数据库名、SSL 模式
- **MySQL**: 主机、端口、用户名、密码、数据库名
- **SQLite3**: 数据库文件路径

### 4. 配置确认机制
- 迁移前显示所有配置信息
- 密码以星号显示，保护隐私
- 用户必须明确确认才能执行迁移

### 5. 数据库连接验证
- 自动测试数据库连接
- 提供清晰的错误信息
- 连接成功后才执行迁移

### 6. 自动化迁移执行
- 使用 Ent 的 Schema.Create 方法
- 支持外键、列删除、索引删除选项
- 详细的执行日志

### 7. 多数据库驱动支持
- PostgreSQL: github.com/jackc/pgx/v5/stdlib
- MySQL: github.com/go-sql-driver/mysql
- SQLite3: github.com/mattn/go-sqlite3

## 文件结构

```
last-admin-core/
├── Makefile                          # 更新：添加 ent-migrate 目标
├── scripts/
│   └── migrate/
│       ├── main.go                   # 迁移工具主程序
│       ├── README.md                 # 详细使用文档
│       ├── EXAMPLES.md               # 使用示例
│       └── IMPLEMENTATION.md         # 本文件
```

## 核心代码结构

### DatabaseConfig 结构体
```go
type DatabaseConfig struct {
    DBType       string  // 数据库类型
    Host         string  // 主机地址
    Port         int     // 端口号
    Username     string  // 用户名
    Password     string  // 密码
    DatabaseName string  // 数据库名
    SSLMode      string  // SSL 模式（仅 PostgreSQL）
}
```

### 主要函数

1. **selectDatabaseType()** - 数据库类型选择
2. **promptInput()** - 字符串输入提示
3. **promptIntInput()** - 整数输入提示
4. **promptYesNo()** - 是/否确认
5. **printConfig()** - 显示配置信息
6. **performMigration()** - 执行数据库迁移

## 使用流程

```
启动工具
    ↓
选择数据库类型
    ↓
输入数据库配置（每项可选）
    ↓
显示配置确认
    ↓
用户确认
    ↓
连接数据库
    ↓
执行迁移
    ↓
显示结果
```

## 默认配置值

### PostgreSQL
```yaml
Host: 127.0.0.1
Port: 5432
Username: postgres
Password: postgres
DatabaseName: last_admin
SSLMode: disable
```

### MySQL
```yaml
Host: 127.0.0.1
Port: 3306
Username: root
Password: root
DatabaseName: last_admin
```

### SQLite3
```yaml
DatabaseName: last_admin.db
```

## 集成方式

### 1. Makefile 集成
```makefile
.PHONY: ent-migrate
ent-migrate: # 同步 Ent 到数据库
	@go run ./scripts/migrate/main.go
```

### 2. 使用方式
```bash
cd last-admin-core
make ent-migrate
```

## 技术细节

### 数据库连接字符串生成

#### PostgreSQL DSN
```
postgresql://username:password@host:port/database?sslmode=disable
```

#### MySQL DSN
```
username:password@tcp(host:port)/database?parseTime=True
```

#### SQLite3 DSN
```
file:path/to/database.db?_busy_timeout=100000&_fk=1
```

### 迁移选项

```go
schema.WithForeignKeys(false)   // 不创建外键约束
schema.WithDropColumn(true)     // 删除未使用的列
schema.WithDropIndex(true)      // 删除未使用的索引
```

## 错误处理

工具提供以下错误处理：

1. **连接错误** - 数据库连接失败
2. **认证错误** - 用户名/密码错误
3. **迁移错误** - 表创建失败
4. **输入错误** - 无效的数字输入

所有错误都会显示详细的错误信息。

## 安全考虑

1. **密码保护** - 密码在显示时以星号替代
2. **无持久化** - 密码不会被保存到任何文件
3. **确认机制** - 迁移前需要用户确认
4. **连接验证** - 执行迁移前验证数据库连接

## 扩展性

工具设计易于扩展：

1. 可以添加更多数据库类型
2. 可以添加更多配置选项
3. 可以集成到 CI/CD 流程
4. 可以添加日志记录功能

## 测试

工具已通过以下测试：

- ✅ PostgreSQL 连接和迁移
- ✅ MySQL 连接和迁移
- ✅ SQLite3 连接和迁移
- ✅ 默认值使用
- ✅ 自定义值输入
- ✅ 取消操作
- ✅ 错误处理

## 依赖

- Go 1.18+
- entgo.io/ent
- github.com/go-sql-driver/mysql
- github.com/jackc/pgx/v5
- github.com/mattn/go-sqlite3

## 相关文档

- [README.md](./README.md) - 详细使用文档
- [EXAMPLES.md](./EXAMPLES.md) - 使用示例
- [../../rpc/ent/migrate/migrate.go](../../rpc/ent/migrate/migrate.go) - Ent 迁移模块

## 后续改进建议

1. 添加迁移历史记录
2. 支持迁移回滚
3. 添加数据库备份功能
4. 支持环境变量配置
5. 添加配置文件支持
6. 集成更多数据库类型
7. 添加性能优化选项

