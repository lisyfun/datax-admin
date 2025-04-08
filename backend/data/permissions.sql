-- 清空权限相关表
TRUNCATE TABLE role_permissions;
TRUNCATE TABLE permissions;

-- 插入权限数据
INSERT INTO permissions (id, name, code, type, parent_id, path, component, icon, sort, status, created_at, updated_at) VALUES
-- 仪表盘
(1, '仪表盘', 'dashboard', 'menu', NULL, '/dashboard', 'dashboard/index', 'icon-dashboard', 1, 1, NOW(), NOW()),
(2, '查看仪表盘', 'dashboard:view', 'button', 1, NULL, NULL, NULL, 1, 1, NOW(), NOW()),

-- 任务管理
(10, '任务管理', 'job', 'menu', NULL, '/job', NULL, 'icon-calendar', 10, 1, NOW(), NOW()),
(11, '任务列表', 'job:list', 'menu', 10, '/job/list', 'job/list/index', 'icon-unordered-list', 1, 1, NOW(), NOW()),
(12, '查看任务', 'job:view', 'button', 11, NULL, NULL, NULL, 1, 1, NOW(), NOW()),
(13, '创建任务', 'job:create', 'button', 11, NULL, NULL, NULL, 2, 1, NOW(), NOW()),
(14, '编辑任务', 'job:update', 'button', 11, NULL, NULL, NULL, 3, 1, NOW(), NOW()),
(15, '删除任务', 'job:delete', 'button', 11, NULL, NULL, NULL, 4, 1, NOW(), NOW()),
(16, '执行任务', 'job:execute', 'button', 11, NULL, NULL, NULL, 5, 1, NOW(), NOW()),

-- 终端管理
(20, '终端管理', 'terminal', 'menu', NULL, '/terminal', NULL, 'icon-desktop', 20, 1, NOW(), NOW()),
(21, '终端列表', 'terminal:list', 'menu', 20, '/terminal/list', 'terminal/list/index', 'icon-unordered-list', 1, 1, NOW(), NOW()),
(22, '查看终端', 'terminal:view', 'button', 21, NULL, NULL, NULL, 1, 1, NOW(), NOW()),
(23, '创建终端', 'terminal:create', 'button', 21, NULL, NULL, NULL, 2, 1, NOW(), NOW()),
(24, '编辑终端', 'terminal:update', 'button', 21, NULL, NULL, NULL, 3, 1, NOW(), NOW()),
(25, '删除终端', 'terminal:delete', 'button', 21, NULL, NULL, NULL, 4, 1, NOW(), NOW()),

-- 工具箱
(30, '工具箱', 'tools', 'menu', NULL, '/tools', NULL, 'icon-bulb', 30, 1, NOW(), NOW()),
(31, 'Kafka工具', 'tools:kafka', 'menu', 30, '/tools/kafka', 'tools/kafka/index', 'icon-cloud', 1, 1, NOW(), NOW()),
(32, '查看Kafka', 'tools:kafka:view', 'button', 31, NULL, NULL, NULL, 1, 1, NOW(), NOW()),
(33, '发送消息', 'tools:kafka:send', 'button', 31, NULL, NULL, NULL, 2, 1, NOW(), NOW()),

-- 系统管理
(40, '系统管理', 'system', 'menu', NULL, '/system', NULL, 'icon-safe', 40, 1, NOW(), NOW()),
(41, '用户管理', 'system:user', 'menu', 40, '/system/user', 'system/user/index', 'icon-user', 1, 1, NOW(), NOW()),
(42, '查看用户', 'system:user:view', 'button', 41, NULL, NULL, NULL, 1, 1, NOW(), NOW()),
(43, '创建用户', 'system:user:create', 'button', 41, NULL, NULL, NULL, 2, 1, NOW(), NOW()),
(44, '编辑用户', 'system:user:update', 'button', 41, NULL, NULL, NULL, 3, 1, NOW(), NOW()),
(45, '删除用户', 'system:user:delete', 'button', 41, NULL, NULL, NULL, 4, 1, NOW(), NOW()),

(46, '角色管理', 'system:role', 'menu', 40, '/system/role', 'system/role/index', 'icon-user-group', 2, 1, NOW(), NOW()),
(47, '查看角色', 'system:role:view', 'button', 46, NULL, NULL, NULL, 1, 1, NOW(), NOW()),
(48, '创建角色', 'system:role:create', 'button', 46, NULL, NULL, NULL, 2, 1, NOW(), NOW()),
(49, '编辑角色', 'system:role:update', 'button', 46, NULL, NULL, NULL, 3, 1, NOW(), NOW()),
(50, '删除角色', 'system:role:delete', 'button', 46, NULL, NULL, NULL, 4, 1, NOW(), NOW()),
(51, '分配权限', 'system:role:permission', 'button', 46, NULL, NULL, NULL, 5, 1, NOW(), NOW()),

(52, '菜单管理', 'system:menu', 'menu', 40, '/system/menu', 'system/menu/index', 'icon-menu', 3, 1, NOW(), NOW()),
(53, '查看菜单', 'system:menu:view', 'button', 52, NULL, NULL, NULL, 1, 1, NOW(), NOW()),
(54, '创建菜单', 'system:menu:create', 'button', 52, NULL, NULL, NULL, 2, 1, NOW(), NOW()),
(55, '编辑菜单', 'system:menu:update', 'button', 52, NULL, NULL, NULL, 3, 1, NOW(), NOW()),
(56, '删除菜单', 'system:menu:delete', 'button', 52, NULL, NULL, NULL, 4, 1, NOW(), NOW()),

(57, '权限管理', 'system:permission', 'menu', 40, '/system/permission', 'system/permission/index', 'icon-lock', 4, 1, NOW(), NOW()),
(58, '查看权限', 'system:permission:view', 'button', 57, NULL, NULL, NULL, 1, 1, NOW(), NOW()),
(59, '创建权限', 'system:permission:create', 'button', 57, NULL, NULL, NULL, 2, 1, NOW(), NOW()),
(60, '编辑权限', 'system:permission:update', 'button', 57, NULL, NULL, NULL, 3, 1, NOW(), NOW()),
(61, '删除权限', 'system:permission:delete', 'button', 57, NULL, NULL, NULL, 4, 1, NOW(), NOW());

-- 为超级管理员角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT 1, id, NOW(), NOW() FROM permissions;
