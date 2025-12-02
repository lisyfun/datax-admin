## 目标与原则
- 新增“Redis 管理”原生功能页面与后端接口，遵循现有分层与命名习惯。
- 不改变既有功能与权限设计；新增文件与路由、权限点保持体系一致。
- 前端使用 `Vue3 + TS + Pinia + Arco`，路由静态声明，菜单由后端权限树控制显示。

## 前端实现
1. 路由
- 在 `web/src/router/index.ts` 静态注册 `path: '/tools/redis'`，与现有 `tools` 组保持一致。
- 侧边栏点击与渲染沿用现有逻辑：`web/src/components/Sidebar/index.vue:200`。
2. 页面结构（新建目录 `web/src/views/tools/redis/`）
- `Index.vue`：整体布局与面包屑。
- `Browser.vue`：Key 列表与筛选（pattern、type）、游标分页（SCAN）。
- `KeyDetail.vue`：按类型查看/编辑（string/hash/list/set/zset），TTL 显示/修改。
- `ConnectionDrawer.vue`：连接新增/编辑/测试。
3. 状态与 API
- `web/src/stores/redis.ts`：当前连接、最近使用、加载状态。
- `web/src/api/redis.ts`：`listKeys/getValue/setValue/deleteKey/getTTL/expire` 与连接管理 `list/create/update/delete/test`。
- 类型 `web/src/types/redis.ts`：`RedisConnection/RedisKey/RedisValue`。

## 后端实现
1. 依赖
- 引入官方客户端 `github.com/redis/go-redis/v9`（最新版）。参考快速用法与基础操作已核对最新文档。
2. 模型（GORM）
- 新增 `models/redis_connection.go`：`id/name/host/port/username/password/db/use_tls/created_at/updated_at/deleted_at`。
- 密码字段加密存储或安全配置管理，严禁明文日志。
3. 服务层 `services/redis.go`
- 连接池封装：按连接ID缓存 `*redis.Client`，统一 `context.Context` 与超时。
- 键操作：`ListKeys(SCAN) / Type / Get / Set / Del / TTL / Expire`。
- 结构操作：`String/Hash/List/Set/ZSet` 读写；必要时细化方法。
4. 控制器 `controllers/redis.go`
- 连接管理：`GET/POST/PUT/DELETE /api/v1/redis/connections`，`POST /api/v1/redis/connections/test`。
- 键操作：`GET /api/v1/redis/keys`（pattern/cursor/count/type）、`GET /api/v1/redis/keys/:key`、`POST /api/v1/redis/keys`、`DELETE /api/v1/redis/keys/:key`、`GET /api/v1/redis/keys/:key/ttl`、`PUT /api/v1/redis/keys/:key/ttl`。
5. 路由注册
- 在 `routes/routes.go` 增加 `RegisterRedisRoutes`，与既有 `api/v1` 分组一致。

## 权限与菜单
- 在 `permissions` 新增菜单：`name='Redis 管理'`、`code='tools.redis'`、`type='menu'`、`path='tools/redis'`、`icon='icon-database'`、`hidden=0`、`cache=0`。
- 可选按钮权限：`tools.redis.query/set/del/ttl/conn.manage`；在 `role_permissions` 给指定角色授权。
- 菜单树由服务层构建并通过 `GET /api/v1/user/menus` 提供，前端侧边栏渲染。

## 安全与性能
- 敏感信息不入日志；对写操作实施权限点与后端鉴权。
- SCAN 游标分页与 `count` 限制，避免阻塞；统一超时与错误处理。
- 仅内网访问或网关鉴权；跨域严格限制。

## 测试与验收
- 单元测试：服务层键与结构操作、连接测试；模拟错误与超时。
- 接口联调：连接创建→列表→键浏览/编辑/TTL→删除流程。
- 前端自测：列表筛选/分页、详情编辑、连接抽屉交互；菜单显示与跳转。

## 任务拆解（任务化）
- 后端：新增 Redis 连接模型与表结构
- 后端：集成 go-redis 并封装连接池
- 后端：实现键操作接口（SCAN/GET/SET/DEL/TTL/EXPIRE）
- 后端：实现结构化数据接口（hash/list/set/zset）
- 后端：实现连接管理 CRUD 与测试接口
- 后端：注册路由并添加权限点
- 前端：新增路由与页面骨架（tools/redis）
- 前端：实现 Key 列表筛选与分页
- 前端：实现 Key 详情查看与编辑
- 前端：实现连接管理抽屉与测试
- 前端：实现 API 模块与 Pinia 状态
- 权限：新增菜单与（可选）按钮权限并授权
- 验证：联调、E2E 流程与单元测试

## 交付物
- 前端：`web/src/views/tools/redis/*`、`web/src/api/redis.ts`、`web/src/stores/redis.ts`、`web/src/types/redis.ts`、`web/src/router/index.ts` 路由项。
- 后端：`models/redis_connection.go`、`services/redis.go`、`controllers/redis.go`、`routes/routes.go` 注册项。
- 数据：`permissions` 与 `role_permissions` 插入脚本（或前端权限页配置）。

请确认上述方案与任务拆解；确认后我将按任务逐项落地并验证。