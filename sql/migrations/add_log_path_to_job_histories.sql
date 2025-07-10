-- 为 job_histories 表添加 log_path 字段
-- 用于存储日志文件路径，替代直接在数据库中存储大量日志内容

-- 添加 log_path 字段
ALTER TABLE `job_histories` ADD COLUMN `log_path` VARCHAR(500) DEFAULT '' COMMENT '日志文件路径';

-- 创建索引以提高查询性能
CREATE INDEX `idx_job_histories_log_path` ON `job_histories` (`log_path`);

-- 注意：现有的 output 和 error 字段将被保留，但在新的实现中不再使用
-- 这些字段将通过 gorm:"-" 标签变成虚拟字段，从日志文件中读取内容
