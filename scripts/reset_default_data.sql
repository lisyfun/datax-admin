-- 重置默认数据脚本
-- 注意：此脚本会删除所有用户、角色和权限数据，请谨慎使用！

-- 禁用外键约束检查（MySQL）
-- SET FOREIGN_KEY_CHECKS = 0;

-- 删除关联表数据
DELETE FROM user_roles;
DELETE FROM role_permissions;

-- 删除主表数据
DELETE FROM users;
DELETE FROM roles;
DELETE FROM permissions;

-- 重置自增ID（MySQL）
-- ALTER TABLE users AUTO_INCREMENT = 1;
-- ALTER TABLE roles AUTO_INCREMENT = 1;
-- ALTER TABLE permissions AUTO_INCREMENT = 1;

-- 重置序列（PostgreSQL）
-- SELECT setval('users_id_seq', 1, false);
-- SELECT setval('roles_id_seq', 1, false);
-- SELECT setval('permissions_id_seq', 1, false);

-- 启用外键约束检查（MySQL）
-- SET FOREIGN_KEY_CHECKS = 1;

-- 执行此脚本后，重启应用程序将自动重新初始化默认数据
