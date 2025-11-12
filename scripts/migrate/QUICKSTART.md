# 快速开始指南 (Quick Start)

## 最快的方式（5 秒钟）

```bash
cd last-admin-core
make ent-migrate
# 按 Enter 多次使用所有默认值
# 最后输入 y 确认
```

完成！数据库已初始化。

## 三种常见场景

### 场景 1：本地开发（使用默认 PostgreSQL）

```bash
make ent-migrate
```

然后：
1. 按 Enter 选择 PostgreSQL（默认）
2. 按 Enter 使用所有默认值
3. 输入 `y` 确认

✅ 完成！

### 场景 2：连接到远程数据库

```bash
make ent-migrate
```

然后：
1. 按 Enter 选择 PostgreSQL
2. 输入远程主机地址（例如：`db.example.com`）
3. 输入端口（例如：`5432`）
4. 输入用户名和密码
5. 输入数据库名称
6. 输入 `y` 确认

✅ 完成！

### 场景 3：使用 SQLite3（快速测试）

```bash
make ent-migrate
```

然后：
1. 输入 `3` 选择 SQLite3
2. 按 Enter 使用默认文件路径或输入自定义路径
3. 输入 `y` 确认

✅ 完成！

## 支持的数据库

| 数据库 | 推荐用途 | 默认端口 |
|--------|---------|---------|
| PostgreSQL | 生产环境、开发环境 | 5432 |
| MySQL | 生产环境、开发环境 | 3306 |
| SQLite3 | 测试、演示 | N/A |

## 默认配置速查表

### PostgreSQL
- 主机: `127.0.0.1`
- 端口: `5432`
- 用户: `postgres`
- 密码: `postgres`
- 数据库: `last_admin`

### MySQL
- 主机: `127.0.0.1`
- 端口: `3306`
- 用户: `root`
- 密码: `root`
- 数据库: `last_admin`

### SQLite3
- 文件: `last_admin.db`

## 常见问题

**Q: 如何取消迁移？**
A: 在最后的确认提示输入 `n` 或按 Enter。

**Q: 如何修改配置？**
A: 在提示时输入新值，或按 Enter 使用默认值。

**Q: 密码会被保存吗？**
A: 不会。密码只在当前迁移过程中使用。

**Q: 迁移会删除数据吗？**
A: 不会。迁移只创建或更新表结构。

**Q: 如何查看详细文档？**
A: 查看 [README.md](./README.md)

## 故障排除

### 连接失败
- 检查数据库服务是否运行
- 验证主机地址和端口
- 检查用户名和密码

### 权限错误
- 确保数据库用户有创建表的权限
- 检查数据库是否存在

### 其他问题
- 查看 [README.md](./README.md) 的故障排除部分
- 查看 [EXAMPLES.md](./EXAMPLES.md) 的示例

## 下一步

迁移完成后，你可以：

1. 启动应用服务
2. 访问 API 文档
3. 开始开发

## 更多帮助

- 📖 [详细文档](./README.md)
- 📝 [使用示例](./EXAMPLES.md)
- 🔧 [实现说明](./IMPLEMENTATION.md)

