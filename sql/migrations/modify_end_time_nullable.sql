-- 修改 job_histories 表的 end_time 字段，允许 NULL 值
-- 这样可以支持正在执行中的任务（end_time 为 NULL）

ALTER TABLE `job_histories` MODIFY COLUMN `end_time` DATETIME NULL DEFAULT NULL COMMENT '结束时间';
