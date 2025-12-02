## 目标
- 将当前“Redis 管理”升级为接近 Another Redis Desktop Manager（RDM）的体验：多连接管理、键树浏览、类型专用编辑器、批处理、Pub/Sub、控制台、集群支持等。
- 保持现有技术栈与代码习惯：后端 Go（Gin+GORM+go-redis v9），前端 Vue3+TS+Pinia+Arco，静态路由+权限树显示。

## 核心功能（RDM 对齐）
1. 多连接与会话
- 连接列表/搜索/收藏、标签页会话并行；支持切换 DB index
- 连接配置支持 TLS、Sentinel、Cluster、超时与重试
- SSH 隧道（计划项，可选）
2. 键树浏览（命名空间）
- 左侧“键树”以冒号命名空间分层（如 `user:1001:profile`）
- SCAN 增量加载、虚拟滚动、类型与大小展示、快速筛选（pattern、type）
3. 类型专用编辑器
- string/hash/list/set/zset/stream/bitmap
- RedisJSON（检测模块可用时启用 JSON 编辑器）
- TTL 查看/编辑、重命名/复制/移动/删除
4. 批处理与导入/导出
- 批量删除、批量 TTL、批量重命名（前缀替换）
- 导出为 JSON/NDJSON；导入 JSON 并可选择覆盖/跳过
- 大批量采用 pipeline、进度与失败重试
5. Pub/Sub 与控制台
- 频道订阅实时查看；发布消息
- 命令控制台（受限白名单），命令历史与收藏
6. 服务信息与监控
- INFO 展示（server/memory/stats），慢查询 slowlog 查看
- 内存占用、连接数、命令统计，定时轮询与节流
7. 集群支持
- Cluster 客户端，节点与槽位视图；按 key 路由到节点
- Sentinel 发现与自动切换（可选）
8. 权限与审计
- 菜单与按钮权限细粒度控制（query/set/del/ttl/conn.manage/console/pubsub/import/export）
- 操作审计（模块 `tools.redis`）记录关键写入动作

## 前端实现
- 目录结构：`web/src/views/tools/redis/`
- 布局
  - 左：连接列表与键树（Tree +虚拟滚动）
  - 右：Tab 工作区（键编辑器、Pub/Sub、控制台、信息）
- 组件
  - `ConnectionPanel`：列表、编辑、测试（抽屉）
  - `KeyTree`：命名空间树、增量加载、筛选
  - `Editor*`：按类型细分（StringEditor/HashEditor/.../JSONEditor）
  - `BulkActions`：批量操作与进度反馈
  - `PubSubPanel`、`ConsolePanel`、`InfoPanel`
- 状态
  - Pinia：连接、当前会话、键缓存、扫描游标、编辑历史、批任务进度
- 交互
  - Tab 会话管理；键打开即生成会话；未保存提示；快捷键（保存/刷新）

## 后端实现
- 客户端与连接
  - 基于 go-redis v9；新增 Cluster 与 Sentinel options
  - 连接缓存按会话隔离；操作超时与重试策略
- API（REST）
  - 连接：CRUD/测试，支持 TLS/Sentinel/Cluster 字段
  - 键：扫描/类型/读/写/删/TTL/重命名/复制/移动
  - 批处理：批量删除/TTL/重命名（前缀替换）
  - 导入/导出：文件上传解析与写入；导出流式返回
  - 信息与监控：`INFO`、`SLOWLOG`、`CONFIG GET`（只读，安全白名单）
- 实时（WebSocket）
  - Pub/Sub 消费与推送；控制台回显；长任务进度
- 安全与资源
  - 命令白名单与速率限制；写操作审计；大批量分批+pipeline；禁止 KEYS；SCAN 限流

## 集群与高级特性
- Cluster：节点列表、槽位分布与健康；按 key 直达节点
- Sentinel：主从发现与切换（可选）
- SSH 隧道：后续版本引入（基于 `golang.org/x/crypto/ssh`）

## 权限与菜单
- 菜单：`tools.redis`
- 按钮权限：`tools.redis.query/set/del/ttl/rename/copy/move/bulk/import/export/console/pubsub/conn.manage`
- 后端校验：按按钮权限限制对应接口；关键写动作写入操作日志

## 交互与性能优化
- 键树增量加载+虚拟滚动；客户端类型缓存（避免重复 TYPE）
- Pipeline 批量；并发限制与进度反馈；错误重试与局部失败汇总
- 大值编辑分块（例如 list/zset 分页操作）

## 交付与里程碑
1. M1：连接面板、键树（string/hash/list/set/zset）、TTL/删除/重命名/复制/移动、导出（JSON）
2. M2：批处理、导入（JSON）、Pub/Sub、控制台（白名单）
3. M3：信息与监控（INFO/slowlog）、Cluster 支持、Sentinel（可选）
4. M4：JSON 模块编辑器、Stream 支持、收藏与历史、SSH 隧道（可选）

## 验收与兼容
- 保持现有路由 `tools/redis` 与权限；渐进式升级 UI
- 提供演示数据与 E2E 测试脚本；前后端单元测试覆盖关键路径
- 文档：使用指南与风险提示（批量操作/导入）

如确认本方案，我将按里程碑逐步落地实现，并保持与现有代码习惯一致。