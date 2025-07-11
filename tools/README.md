# DataX Admin 数据重置工具

这是一个独立的数据重置工具，用于清空并重新初始化 DataX Admin 系统的用户、角色和权限数据。

## 功能特性

- 🗄️ **多数据库支持**: 支持 MySQL 和 PostgreSQL
- 📋 **配置文件选择**: 可以指定不同的配置文件
- 🔄 **安全重置**: 按照外键依赖顺序安全清理数据
- 🔢 **序列重置**: 自动重置自增ID和序列
- ⚠️ **交互确认**: 防止误操作的确认机制
- 📖 **详细帮助**: 完整的使用说明

## 使用方法

### 基本用法

```bash
# 使用默认MySQL配置
go run tools/reset_data.go

# 指定数据库类型
go run tools/reset_data.go -db mysql
go run tools/reset_data.go -db postgres

# 指定配置文件
go run tools/reset_data.go -config config-mysql.yaml
go run tools/reset_data.go -config config-postgres.yaml

# 显示帮助信息
go run tools/reset_data.go -help
```

### 编译后使用

```bash
# 编译工具
go build -o reset_tool tools/reset_data.go

# 运行工具
./reset_tool -db mysql
./reset_tool -config config-postgres.yaml
```

## 命令行参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-db` | string | mysql | 数据库类型 (mysql\|postgres) |
| `-config` | string | 自动选择 | 配置文件路径 |
| `-help` | bool | false | 显示帮助信息 |

## 配置文件映射

当使用 `-db` 参数时，工具会自动选择对应的配置文件：

- `mysql` → `config-mysql.yaml`
- `postgres` → `config-postgres.yaml`
- 默认 → `config.yaml`

## 执行流程

1. **参数解析**: 解析命令行参数
2. **配置文件检查**: 验证配置文件是否存在
3. **配置初始化**: 加载指定的配置文件
4. **数据库连接**: 连接到目标数据库
5. **信息显示**: 显示数据库连接信息
6. **用户确认**: 要求用户确认操作
7. **数据清理**: 按顺序删除数据
8. **序列重置**: 重置自增ID/序列
9. **数据初始化**: 重新创建默认数据

## 数据清理顺序

工具按照以下顺序清理数据，确保外键约束不会出错：

1. `user_roles` - 用户角色关联表
2. `role_permissions` - 角色权限关联表  
3. `users` - 用户表
4. `roles` - 角色表
5. `permissions` - 权限表

## 序列重置

### MySQL
```sql
ALTER TABLE users AUTO_INCREMENT = 1;
ALTER TABLE roles AUTO_INCREMENT = 1;
ALTER TABLE permissions AUTO_INCREMENT = 1;
```

### PostgreSQL
```sql
SELECT setval('users_id_seq', 1, false);
SELECT setval('roles_id_seq', 1, false);
SELECT setval('permissions_id_seq', 1, false);
```

## 默认数据

重置完成后，系统会创建以下默认数据：

### 默认角色
- **超级管理员** (admin): 拥有所有权限
- **普通用户** (user): 基础权限

### 默认用户
- **用户名**: admin
- **密码**: admin123
- **角色**: 超级管理员

### 权限菜单
- 完整的菜单权限体系 (57个权限项)
- 包括所有CRUD操作权限

## 使用示例

### 示例1: 重置MySQL数据
```bash
$ go run tools/reset_data.go -db mysql
📋 使用配置文件: config-mysql.yaml
🗄️  数据库类型: mysql
🔗 数据库连接: root@localhost:3306/datax_admin

⚠️  警告：此操作将删除所有用户、角色和权限数据！
确认要继续吗？(y/N): y
🔄 开始重置数据...
数据清空完成
✅ 数据重置完成！
👤 默认管理员账号：
   用户名: admin
   密码: admin123
```

### 示例2: 重置PostgreSQL数据
```bash
$ go run tools/reset_data.go -db postgres
📋 使用配置文件: config-postgres.yaml
🗄️  数据库类型: postgres
🔗 数据库连接: postgres@localhost:5432/datax_admin

⚠️  警告：此操作将删除所有用户、角色和权限数据！
确认要继续吗？(y/N): y
🔄 开始重置数据...
数据清空完成
✅ 数据重置完成！
👤 默认管理员账号：
   用户名: admin
   密码: admin123
```

## 注意事项

1. **数据备份**: 执行重置前请务必备份重要数据
2. **权限确认**: 确保数据库用户有删除表数据的权限
3. **服务停止**: 建议在重置前停止应用服务
4. **配置检查**: 确认配置文件中的数据库连接信息正确

## 故障排除

### 常见错误

**配置文件不存在**
```
❌ 配置文件不存在: config-postgres.yaml
```
解决方案: 检查配置文件是否存在，或使用 `-config` 参数指定正确的文件路径

**数据库连接失败**
```
failed to connect to database
```
解决方案: 检查数据库服务是否运行，配置文件中的连接信息是否正确

**权限不足**
```
删除用户失败: permission denied
```
解决方案: 确保数据库用户有删除表数据的权限

## 相关文件

- `tools/reset_data.go` - 重置工具源码
- `config.yaml` - 默认配置文件
- `config-mysql.yaml` - MySQL配置文件
- `config-postgres.yaml` - PostgreSQL配置文件
