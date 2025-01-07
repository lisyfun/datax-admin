-- 创建菜单表
CREATE TABLE IF NOT EXISTS menus (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    parent_id   BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父菜单ID',
    name        VARCHAR(50) NOT NULL COMMENT '菜单名称',
    path        VARCHAR(200) COMMENT '路由路径',
    component   VARCHAR(200) COMMENT '组件路径',
    icon        VARCHAR(50) COMMENT '图标',
    sort        INT NOT NULL DEFAULT 0 COMMENT '排序',
    status      TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    hidden      TINYINT NOT NULL DEFAULT 0 COMMENT '是否隐藏',
    cache       TINYINT NOT NULL DEFAULT 1 COMMENT '是否缓存',
    type        TINYINT NOT NULL DEFAULT 1 COMMENT '类型：1-菜单，2-按钮',
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,
    INDEX idx_parent_id (parent_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单表';
