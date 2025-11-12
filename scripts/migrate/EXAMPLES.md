# 数据库迁移工具使用示例

## 示例 1：使用默认 PostgreSQL 配置

```bash
$ cd last-admin-core
$ make ent-migrate
```

然后按照提示操作：

```
========================================
   数据库迁移工具 - Database Migration
========================================

请选择数据库类型:
1. PostgreSQL (推荐)
2. MySQL
3. SQLite3

请输入选项 (1-3) [默认: 1]: 
# 按 Enter 使用默认值 PostgreSQL

请输入数据库配置信息 (按 Enter 使用默认值):

数据库主机地址 [127.0.0.1]: 
# 按 Enter 使用默认值

数据库端口 [5432]: 
# 按 Enter 使用默认值

数据库用户名 [postgres]: 
# 按 Enter 使用默认值

数据库密码 [postgres]: 
# 按 Enter 使用默认值

数据库名称 [last_admin]: 
# 按 Enter 使用默认值

SSL 模式 (disable/require/prefer) [disable]: 
# 按 Enter 使用默认值

========================================
   确认数据库配置
========================================
数据库类型: postgres
主机地址: 127.0.0.1
端口: 5432
用户名: postgres
密码: ********
数据库名: last_admin
SSL 模式: disable
========================================
是否继续执行迁移? (y/n) [默认: n]: y
# 输入 y 确认

正在执行数据库迁移...
✓ 数据库连接成功
✓ 正在创建/更新数据库表...
✓ 数据库表创建/更新完成
✅ 数据库迁移成功!
```

## 示例 2：使用自定义 PostgreSQL 配置

```bash
$ make ent-migrate
```

选择 PostgreSQL 后，输入自定义配置：

```
请输入数据库配置信息 (按 Enter 使用默认值):

数据库主机地址 [127.0.0.1]: db.example.com
# 输入远程数据库地址

数据库端口 [5432]: 5433
# 输入自定义端口

数据库用户名 [postgres]: admin
# 输入自定义用户名

数据库密码 [postgres]: MySecurePassword123
# 输入自定义密码

数据库名称 [last_admin]: my_app_db
# 输入自定义数据库名

SSL 模式 (disable/require/prefer) [disable]: require
# 选择 SSL 模式
```

## 示例 3：使用 MySQL

```bash
$ make ent-migrate
```

选择 MySQL：

```
请选择数据库类型:
1. PostgreSQL (推荐)
2. MySQL
3. SQLite3

请输入选项 (1-3) [默认: 1]: 2
# 选择 MySQL

请输入数据库配置信息 (按 Enter 使用默认值):

数据库主机地址 [127.0.0.1]: 
# 按 Enter 使用默认值

数据库端口 [3306]: 
# 按 Enter 使用默认值

数据库用户名 [root]: 
# 按 Enter 使用默认值

数据库密码 [root]: 
# 按 Enter 使用默认值

数据库名称 [last_admin]: 
# 按 Enter 使用默认值
```

## 示例 4：使用 SQLite3

```bash
$ make ent-migrate
```

选择 SQLite3：

```
请选择数据库类型:
1. PostgreSQL (推荐)
2. MySQL
3. SQLite3

请输入选项 (1-3) [默认: 1]: 3
# 选择 SQLite3

请输入数据库配置信息 (按 Enter 使用默认值):

数据库文件路径 [last_admin.db]: ./data/app.db
# 输入自定义文件路径

========================================
   确认数据库配置
========================================
数据库类型: sqlite3
数据库文件: ./data/app.db
========================================
是否继续执行迁移? (y/n) [默认: n]: y
# 输入 y 确认

正在执行数据库迁移...
✓ 数据库连接成功
✓ 正在创建/更新数据库表...
✓ 数据库表创建/更新完成
✅ 数据库迁移成功!
```

## 示例 5：自动化脚本中使用（非交互模式）

虽然工具主要设计为交互式，但可以通过管道输入来自动化：

```bash
# 使用默认 PostgreSQL 配置并自动确认
echo -e "1\n\n\n\n\n\n\ny" | make ent-migrate

# 使用 MySQL 并自动确认
echo -e "2\n\n\n\n\n\ny" | make ent-migrate

# 使用 SQLite3 并自动确认
echo -e "3\n\n\ny" | make ent-migrate

# 使用自定义配置
echo -e "1\ndb.example.com\n5433\nadmin\npassword\nmy_db\nrequire\ny" | make ent-migrate
```

## 示例 6：Docker 环境中使用

如果在 Docker 容器中运行，确保数据库服务可访问：

```bash
# 在 Docker 容器中
docker exec last-admin-core make ent-migrate

# 或者
docker exec last-admin-core sh -c 'echo -e "1\n\n\n\n\n\n\ny" | make ent-migrate'
```

## 常见场景

### 场景 1：开发环境初始化

使用所有默认值快速初始化本地开发数据库：

```bash
make ent-migrate
# 按 Enter 多次使用所有默认值
# 最后输入 y 确认
```

### 场景 2：生产环境部署

使用自定义配置连接到生产数据库：

```bash
make ent-migrate
# 输入生产环境的数据库配置
# 仔细检查配置后输入 y 确认
```

### 场景 3：测试环境

使用 SQLite3 进行快速测试：

```bash
make ent-migrate
# 选择 3 (SQLite3)
# 输入测试数据库文件路径
# 输入 y 确认
```

### 场景 4：CI/CD 流程

在自动化脚本中使用：

```bash
#!/bin/bash
set -e

cd last-admin-core

# 使用默认配置自动迁移
echo -e "1\n\n\n\n\n\n\ny" | make ent-migrate

echo "Database migration completed successfully"
```

## 故障排除示例

### 问题：连接被拒绝

```
❌ 迁移失败: failed to connect to database: connection refused
```

**解决方案：**
- 检查数据库服务是否运行
- 验证主机地址和端口是否正确
- 检查防火墙设置

### 问题：认证失败

```
❌ 迁移失败: failed to connect to database: authentication failed
```

**解决方案：**
- 验证用户名和密码
- 检查数据库用户权限
- 确保用户有创建表的权限

### 问题：数据库不存在

```
❌ 迁移失败: database "last_admin" does not exist
```

**解决方案：**
- 手动创建数据库
- 或使用具有创建数据库权限的用户

## 更多帮助

查看 [README.md](./README.md) 了解更多详细信息。

