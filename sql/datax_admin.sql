/*
 Navicat Premium Dump SQL

 Source Server         : 127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 80042 (8.0.42)
 Source Host           : 127.0.0.1:13306
 Source Schema         : datax_admin

 Target Server Type    : MySQL
 Target Server Version : 80042 (8.0.42)
 File Encoding         : 65001

 Date: 10/07/2025 21:44:56
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for job_histories
-- ----------------------------
DROP TABLE IF EXISTS `job_histories`;
CREATE TABLE `job_histories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `job_id` bigint unsigned DEFAULT NULL,
  `status` bigint DEFAULT NULL,
  `start_time` datetime(3) DEFAULT NULL,
  `end_time` datetime(3) DEFAULT NULL,
  `duration` bigint DEFAULT NULL,
  `output` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `error` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `log_path` varchar(500) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_job_id` (`job_id`),
  KEY `idx_job_histories_log_path` (`log_path`)
) ENGINE=InnoDB AUTO_INCREMENT=5946 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of job_histories
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for jobs
-- ----------------------------
DROP TABLE IF EXISTS `jobs`;
CREATE TABLE `jobs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `cron_expr` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `status` tinyint DEFAULT '0',
  `timeout` int DEFAULT '0',
  `retry_count` int DEFAULT '0',
  `retry_delay` int DEFAULT '0',
  `params` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `creator` bigint unsigned NOT NULL,
  `updater` bigint unsigned DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_jobs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of jobs
-- ----------------------------
BEGIN;
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, '数据库备份', 'shell', '每天凌晨3点备份MySQL数据库', '0 0 3 * * *', 0, 3600, 3, 300, '{\"command\":\"mysqldump -u root -p123456 mydb \\u003e /backup/mydb_$(date +%Y%m%d).sql\",\"work_dir\":\"/backup\",\"environment\":{\"MYSQL_PWD\":\"123456\"}}', 1, 0, '2025-01-06 13:19:27.000', '2025-01-10 12:30:43.000', '2025-01-10 12:30:44.000');
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, '服务健康检查', 'http', '每1分钟检查一次服务健康状态', '0 0/1 * * * ?', 0, 30, 3, 60, '{\"url\":\"http://localhost:28080/api/v1/ping\",\"method\":\"GET\",\"headers\":{},\"body\":\"\",\"success_code\":[200]}', 1, 1, '2025-01-06 13:19:33.000', '2025-06-18 16:15:18.840', NULL);
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, '数据同步', 'datax', '每小时同步一次数据', '0 0 * * * *', 0, 7200, 2, 600, '{\"job_path\":\"/datax/job/mysql_to_mysql.json\",\"parameters\":{\"batch_size\":\"1000\",\"from_date\":\"${yyyy-MM-dd}\"},\"jvm_options\":[\"-Xms1g\",\"-Xmx2g\"],\"speed\":10000}', 1, 0, '2025-01-06 13:19:42.000', '2025-01-10 12:24:03.000', '2025-01-10 12:24:04.000');
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, '测试任务', 'shell', '这是一个测试任务', '* 0/1 * * * ?', 0, 3600, 3, 60, '{\"command\":\"echo 1\",\"work_dir\":\"\",\"environment\":{}}', 1, 1, '2025-01-06 16:57:31.000', '2025-04-11 14:29:06.342', NULL);
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 'test', 'shell', 'test', '* * * * * ?', 0, 2, 2, 34, '{\"command\":\"echo 1\",\"work_dir\":\"\",\"environment\":{}}', 1, 1, '2025-01-06 17:27:00.000', '2025-01-06 17:54:40.000', '2025-01-06 17:54:41.000');
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 'HTTP任务', 'http', '测试HTTP任务', '0 2 1 3 * ?', 0, 30, 3, 5, '{\"url\":\"http://localhost:8080/api/v1/ping\",\"method\":\"GET\",\"headers\":{\"Content-Type\":\"application/json\"},\"body\":\"\",\"success_code\":[200,201]}', 1, 1, '2025-01-06 17:34:12.000', '2025-01-10 12:30:41.000', '2025-01-10 12:30:41.000');
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, 'mysql到 mysql', 'datax', 'mysql到 mysql', '0 0 0 * * ?', 0, 0, 0, 0, '{\"job_config\":\"{\\\"job\\\":{\\\"content\\\":[{\\\"reader\\\":{\\\"name\\\":\\\"mysqlreader\\\",\\\"parameter\\\":{\\\"username\\\":\\\"root\\\",\\\"password\\\":\\\"123456\\\",\\\"host\\\":\\\"host.docker.internal\\\",\\\"port\\\":13306,\\\"database\\\":\\\"datax_source\\\",\\\"table\\\":\\\"users\\\",\\\"columns\\\":[\\\"username\\\",\\\"email\\\",\\\"password\\\",\\\"full_name\\\",\\\"phone\\\",\\\"status\\\",\\\"created_at\\\",\\\"updated_at\\\"],\\\"selectSql\\\":\\\"select username,email,PASSWORD,full_name,phone,STATUS,created_at,updated_at from users\\\",\\\"where\\\":\\\"id \\u003e1000 and created_at \\u003e= ${time}\\\"}},\\\"writer\\\":{\\\"name\\\":\\\"mysqlwriter\\\",\\\"parameter\\\":{\\\"username\\\":\\\"root\\\",\\\"password\\\":\\\"123456\\\",\\\"host\\\":\\\"host.docker.internal\\\",\\\"port\\\":13306,\\\"database\\\":\\\"datax_target\\\",\\\"table\\\":\\\"users\\\",\\\"columns\\\":[\\\"username\\\",\\\"email\\\",\\\"PASSWORD\\\",\\\"full_name\\\",\\\"phone\\\",\\\"STATUS\\\",\\\"created_at\\\",\\\"updated_at\\\"],\\\"preSql\\\":[\\\"select count(*) from users\\\",\\\"truncate users\\\"],\\\"postSql\\\":[\\\"select count(*) from users\\\"],\\\"writeMode\\\":\\\"replace\\\"}}}],\\\"setting\\\":{\\\"speed\\\":{\\\"channel\\\":24,\\\"bytes\\\":52428800},\\\"errorLimit\\\":{\\\"record\\\":0,\\\"percentage\\\":0.02}}}}\",\"parameters\":{\"datax.job.setting.speed.channel\":\"5\",\"time\":\"2024-01-01\"}}', 1, 1, NULL, '2025-07-10 18:34:04.294', NULL);
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, 'MySQL到PostgreSQL同步', 'datax', '从源MySQL数据库同步users表到目标PostgreSQL数据库', '0 0 * * * *', 0, 3600, 3, 300, '{\"job_config\":\"{\\\"job\\\":{\\\"content\\\":[{\\\"reader\\\":{\\\"name\\\":\\\"mysqlreader\\\",\\\"parameter\\\":{\\\"username\\\":\\\"root\\\",\\\"password\\\":\\\"123456\\\",\\\"host\\\":\\\"localhost\\\",\\\"port\\\":13306,\\\"database\\\":\\\"datax_source1\\\",\\\"table\\\":\\\"users\\\",\\\"columns\\\":[\\\"username\\\",\\\"email\\\",\\\"age\\\",\\\"created_at\\\",\\\"updated_at\\\"],\\\"selectSql\\\":\\\"select username, email, age, created_at, updated_at from users\\\",\\\"where\\\":\\\"1=1\\\",\\\"batchSize\\\":20000}},\\\"writer\\\":{\\\"name\\\":\\\"postgresqlwriter\\\",\\\"parameter\\\":{\\\"username\\\":\\\"postgres\\\",\\\"password\\\":\\\"123456\\\",\\\"host\\\":\\\"localhost\\\",\\\"port\\\":15432,\\\"database\\\":\\\"target_db\\\",\\\"schema\\\":\\\"public\\\",\\\"table\\\":\\\"users\\\",\\\"columns\\\":[\\\"username\\\",\\\"email\\\",\\\"age\\\",\\\"created_at\\\",\\\"updated_at\\\"],\\\"batchSize\\\":20000,\\\"preSql\\\":[\\\"select count(*) as total_count from users\\\",\\\"truncate table users\\\"],\\\"postSql\\\":[\\\"select count(*) as total_count from users\\\"],\\\"writeMode\\\":\\\"insert\\\"}}}],\\\"setting\\\":{\\\"speed\\\":{\\\"channel\\\":24,\\\"bytes\\\":52428800},\\\"errorLimit\\\":{\\\"record\\\":0,\\\"percentage\\\":0.02}}}}\",\"parameters\":{}}', 1, 1, '2025-01-10 11:35:20.000', '2025-02-26 16:38:34.081', NULL);
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, 'Oracle到MySQL同步', 'datax', '从源Oracle数据库同步users表到目标MySQL数据库', '0 0 * * * *', 0, 3600, 3, 300, '{\"job_config\":\"{\\\"job\\\":{\\\"content\\\":[{\\\"reader\\\":{\\\"name\\\":\\\"oraclereader\\\",\\\"parameter\\\":{\\\"username\\\":\\\"sys\\\",\\\"password\\\":\\\"123456\\\",\\\"host\\\":\\\"127.0.0.1\\\",\\\"port\\\":1521,\\\"service\\\":\\\"ORCL\\\",\\\"table\\\":\\\"USERS\\\",\\\"columns\\\":[],\\\"selectSql\\\":\\\"select ID, NAME, EMAIL, AGE from USERS\\\",\\\"where\\\":\\\"id \\u003e 0\\\",\\\"batchSize\\\":20000}},\\\"writer\\\":{\\\"name\\\":\\\"mysqlwriter\\\",\\\"parameter\\\":{\\\"username\\\":\\\"root\\\",\\\"password\\\":\\\"123456\\\",\\\"host\\\":\\\"127.0.0.1\\\",\\\"port\\\":13306,\\\"database\\\":\\\"datax_target\\\",\\\"table\\\":\\\"users\\\",\\\"columns\\\":[\\\"id\\\",\\\"username\\\",\\\"email\\\",\\\"age\\\"],\\\"preSql\\\":[\\\"select count(*) from users\\\",\\\"truncate users\\\"],\\\"postSql\\\":[\\\"select count(*) from users\\\"],\\\"batchSize\\\":1000,\\\"writeMode\\\":\\\"insert\\\"}}}],\\\"setting\\\":{\\\"speed\\\":{\\\"channel\\\":24,\\\"bytes\\\":52428800},\\\"errorLimit\\\":{\\\"record\\\":0,\\\"percentage\\\":0.02}}}}\",\"parameters\":{\"datax.job.setting.speed.channel\":\"5\"}}', 1, 1, '2025-01-10 11:35:36.000', '2025-02-10 15:40:59.143', NULL);
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (10, '1', 'shell', '1', '* * * * * ?', 0, 0, 0, 0, '{\"command\":\"1\",\"work_dir\":\"1\",\"environment\":{}}', 1, 0, '2025-02-10 14:34:10.858', '2025-02-10 14:34:10.858', '2025-02-10 14:37:41.404');
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (11, '2', 'datax', '2', '* * * * * ?', 0, 0, 0, 0, '{\"job_config\":\"{\\\"job\\\":{\\\"content\\\":[{\\\"reader\\\":{\\\"name\\\":\\\"mysqlreader\\\",\\\"parameter\\\":{\\\"username\\\":\\\"lisongyu\\\",\\\"password\\\":\\\"HrQfWc9-WAmNjV4\\\",\\\"host\\\":\\\"localhost\\\",\\\"port\\\":3306,\\\"database\\\":\\\"1\\\",\\\"table\\\":\\\"1\\\",\\\"columns\\\":[\\\"1\\\"],\\\"where\\\":\\\"1=1\\\",\\\"batchSize\\\":20000,\\\"selectSql\\\":\\\"1\\\"}},\\\"writer\\\":{\\\"name\\\":\\\"mysqlwriter\\\",\\\"parameter\\\":{\\\"username\\\":\\\"lisongyu\\\",\\\"password\\\":\\\"HrQfWc9-WAmNjV4\\\",\\\"host\\\":\\\"localhost\\\",\\\"port\\\":3306,\\\"database\\\":\\\"2\\\",\\\"table\\\":\\\"2\\\",\\\"columns\\\":[\\\"2\\\"],\\\"writeMode\\\":\\\"insert\\\",\\\"batchSize\\\":10000,\\\"preSql\\\":[\\\"2\\\"],\\\"postSql\\\":[\\\"2\\\"]}}}],\\\"setting\\\":{\\\"speed\\\":{\\\"channel\\\":24,\\\"bytes\\\":52428800},\\\"errorLimit\\\":{\\\"record\\\":0,\\\"percentage\\\":0.02}}}}\",\"parameters\":{}}', 1, 0, '2025-02-10 14:34:39.992', '2025-02-10 14:34:39.992', '2025-02-10 14:37:39.986');
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (12, 'test', 'datax', '1', '* * * * * ?', 0, 0, 0, 0, '{\"job_config\":\"{\\\"job\\\":{\\\"content\\\":[{\\\"reader\\\":{\\\"name\\\":\\\"mysqlreader\\\",\\\"parameter\\\":{\\\"username\\\":\\\"lisongyu\\\",\\\"password\\\":\\\"HrQfWc9-WAmNjV4\\\",\\\"host\\\":\\\"localhost\\\",\\\"port\\\":3306,\\\"database\\\":\\\"1\\\",\\\"table\\\":\\\"q\\\",\\\"columns\\\":[\\\"username\\\",\\\"email\\\",\\\"age\\\",\\\"created_at\\\",\\\"updated_at\\\"],\\\"where\\\":\\\"1=1\\\",\\\"batchSize\\\":20000,\\\"selectSql\\\":\\\"select username, email, age, created_at, updated_at from users\\\"}},\\\"writer\\\":{\\\"name\\\":\\\"mysqlwriter\\\",\\\"parameter\\\":{\\\"username\\\":\\\"lisongyu\\\",\\\"password\\\":\\\"HrQfWc9-WAmNjV4\\\",\\\"host\\\":\\\"localhost\\\",\\\"port\\\":3306,\\\"database\\\":\\\"1\\\",\\\"table\\\":\\\"2\\\",\\\"columns\\\":[\\\"username\\\",\\\"email\\\",\\\"age\\\",\\\"created_at\\\",\\\"updated_at\\\"],\\\"writeMode\\\":\\\"insert\\\",\\\"batchSize\\\":10000,\\\"preSql\\\":[],\\\"postSql\\\":[]}}}],\\\"setting\\\":{\\\"speed\\\":{\\\"channel\\\":24,\\\"bytes\\\":52428800},\\\"errorLimit\\\":{\\\"record\\\":0,\\\"percentage\\\":0.02}}}}\",\"parameters\":{}}', 1, 0, '2025-02-22 15:15:03.344', '2025-03-13 15:13:33.686', '2025-03-28 21:31:17.442');
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (13, '错误测试任务', 'shell', '测试错误日志', '0 0 1 1 * ?', 0, 0, 0, 0, '{\"command\":\"nonexistentcommand\",\"work_dir\":\"\",\"environment\":{}}', 1, 0, '2025-07-10 13:33:28.533', '2025-07-10 13:33:28.533', NULL);
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (14, '长时间测试任务', 'shell', '测试实时日志功能', '0 0 1 1 * ?', 0, 0, 0, 0, '{\"command\":\"for i in {1..10}; do echo \\\"第 $i 步执行中...\\\"; sleep 2; done; echo \\\"任务完成！\\\"\",\"work_dir\":\"\",\"environment\":{}}', 1, 0, '2025-07-10 13:45:41.019', '2025-07-10 13:45:41.019', NULL);
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (15, '超长时间测试任务', 'shell', '测试实时日志刷新功能', '0 0 1 1 * ?', 0, 0, 0, 0, '{\"command\":\"for i in {1..20}; do echo \\\"第 $i 步执行中，时间: $(date)\\\"; sleep 3; done; echo \\\"任务完成！\\\"\",\"work_dir\":\"\",\"environment\":{}}', 1, 0, '2025-07-10 13:49:43.754', '2025-07-10 13:49:43.754', NULL);
COMMIT;

-- ----------------------------
-- Table structure for kafka_clusters
-- ----------------------------
DROP TABLE IF EXISTS `kafka_clusters`;
CREATE TABLE `kafka_clusters` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(100) NOT NULL,
  `broker_servers` varchar(1000) NOT NULL,
  `security_protocol` varchar(50) DEFAULT NULL,
  `sasl_mechanism` varchar(50) DEFAULT NULL,
  `username` varchar(100) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `description` varchar(500) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delay_message` tinyint(1) DEFAULT '0',
  `status` tinyint(1) DEFAULT '1',
  `last_check_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_kafka_clusters_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Kafka集群配置';

-- ----------------------------
-- Records of kafka_clusters
-- ----------------------------
BEGIN;
INSERT INTO `kafka_clusters` (`id`, `name`, `broker_servers`, `security_protocol`, `sasl_mechanism`, `username`, `password`, `description`, `created_at`, `updated_at`, `delay_message`, `status`, `last_check_time`) VALUES (8, 'host.docker.internal', 'host.docker.internal:9092', 'SASL_PLAINTEXT', 'PLAIN', 'user', '123456', '', '2025-05-13 09:54:35.000', '2025-07-10 21:17:07.530', 0, 1, '2025-07-10 21:17:07.524');
COMMIT;

-- ----------------------------
-- Table structure for kafka_topics
-- ----------------------------
DROP TABLE IF EXISTS `kafka_topics`;
CREATE TABLE `kafka_topics` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `cluster_id` bigint unsigned NOT NULL,
  `name` varchar(191) NOT NULL,
  `partitions` bigint DEFAULT NULL,
  `replicas` bigint DEFAULT NULL,
  `avg_log_size` bigint DEFAULT NULL,
  `total_log_size` bigint DEFAULT NULL,
  `last_refresh_time` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_cluster_topic` (`cluster_id`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of kafka_topics
-- ----------------------------


-- ----------------------------
-- Table structure for login_logs
-- ----------------------------
DROP TABLE IF EXISTS `login_logs`;
CREATE TABLE `login_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名',
  `ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '登录IP',
  `login_time` datetime(3) NOT NULL COMMENT '登录时间',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=273 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of login_logs
-- ----------------------------

-- ----------------------------
-- Table structure for operation_logs
-- ----------------------------
DROP TABLE IF EXISTS `operation_logs`;
CREATE TABLE `operation_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL COMMENT '操作用户ID',
  `username` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '操作用户名',
  `module` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '操作模块',
  `action` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '操作动作',
  `description` varchar(500) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作描述',
  `ip` varchar(50) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作IP',
  `user_agent` varchar(500) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户代理',
  `request_data` text COLLATE utf8mb4_general_ci COMMENT '请求数据',
  `status` bigint DEFAULT '1' COMMENT '操作状态 1:成功 0:失败',
  `error_msg` varchar(500) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '错误信息',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of operation_logs
-- ----------------------------


-- ----------------------------
-- Table structure for permissions
-- ----------------------------
DROP TABLE IF EXISTS `permissions`;
CREATE TABLE `permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `parent_id` bigint unsigned DEFAULT NULL,
  `path` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `component` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `sort` bigint DEFAULT '0',
  `status` bigint DEFAULT '1',
  `hidden` bigint DEFAULT '0',
  `cache` bigint DEFAULT '1',
  `is_external` bigint DEFAULT '0' COMMENT '是否为外部链接 0-否 1-是',
  `external_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '外部链接地址',
  `open_type`  bigint DEFAULT '1' COMMENT '打开方式 0-内嵌 1-新窗口',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_permissions_code` (`code`),
  KEY `idx_permissions_parent_id` (`parent_id`),
  KEY `idx_permissions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='权限表（包含菜单和按钮权限）';

-- ----------------------------
-- Records of permissions
-- ----------------------------
BEGIN;
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, '首页', 'root', 'menu', NULL, '/', 'layouts/default.vue', '', 0, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, '数据面板', 'dashboard', 'menu', 1, 'dashboard', 'views/dashboard/index.vue', 'icon-dashboard', 0, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, '任务管理', 'job', 'menu', 1, 'job', 'views/job/index.vue', 'icon-calendar', 1, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, '任务列表', 'job.list', 'menu', 3, 'list', 'views/job/list/index.vue', 'icon-unordered-list', 0, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, '执行历史', 'job.history', 'menu', 3, 'history', 'views/job/history/index.vue', 'icon-clock-circle', 1, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, '终端管理', 'terminal', 'menu', 1, 'terminal', 'views/terminal/index.vue', 'icon-command', 2, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, '终端列表', 'terminal.list', 'menu', 6, 'list', 'views/terminal/list/index.vue', 'icon-desktop', 0, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, '终端连接', 'terminal.connect', 'menu', 6, 'connect/:id', 'views/terminal/connect/index.vue', '', 1, 1, 1, 1, NULL, '2025-06-21 22:10:48.052', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, '工具管理', 'tools', 'menu', 1, 'tools', 'views/tools/index.vue', 'icon-common', 4, 1, 0, 1, NULL, '2025-06-19 21:59:52.790', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (10, 'Json格式化', 'tools.json-formatter', 'menu', 9, 'json-formatter', 'views/tools/JsonFormatter.vue', 'icon-code', 0, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (11, '加解密工具', 'tools.crypto', 'menu', 9, 'crypto', 'views/tools/Crypto.vue', 'icon-lock', 1, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (12, '消息管理', 'tools.kafka', 'menu', 1, 'kafka', 'views/kafka/index.vue', 'icon-apps', 3, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (13, '集群管理', 'tools.kafka.cluster', 'menu', 12, 'cluster', 'views/kafka/cluster/index.vue', 'icon-apps', 0, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (14, '主题管理', 'tools.kafka.topic', 'menu', 12, 'kafka/cluster/:clusterId/topic', 'views/kafka/topic/index.vue', '', 1, 1, 1, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (15, '消息列表', 'tools.kafka.message', 'menu', 12, 'kafka/clusters/:clusterId/topics/:topicName/messages', 'views/kafka/topic/components/MessageList.vue', '', 2, 1, 1, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (16, '系统管理', 'system', 'menu', 1, 'system', 'views/system/index.vue', 'icon-apps', 999, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (17, '用户管理', 'system.users', 'menu', 16, 'users', 'views/system/users/index.vue', 'icon-user', 0, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (18, '角色管理', 'system.roles', 'menu', 16, 'roles', 'views/system/roles/index.vue', 'icon-user-group', 1, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (19, '权限管理', 'system.permissions', 'menu', 16, 'permissions', 'views/system/permissions/index.vue', 'icon-safe', 2, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (20, '登录', 'login', 'menu', NULL, '/login', 'views/login/index.vue', '', 0, 1, 1, 0, NULL, '2025-06-19 23:24:00.091', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (21, '注册', 'register', 'menu', NULL, '/register', 'views/register/index.vue', '', 0, 0, 1, 0, NULL, '2025-06-19 23:24:01.129', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (22, '用户查询', 'system.users.query', 'button', 17, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (23, '用户创建', 'system.users.create', 'button', 17, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (24, '用户编辑', 'system.users.update', 'button', 17, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (25, '用户删除', 'system.users.delete', 'button', 17, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (26, '角色查询', 'system.roles.query', 'button', 18, '', '', '', 0, 1, 0, 0, NULL, '2025-06-21 22:48:04.149', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (27, '角色创建', 'system.roles.create', 'button', 18, '', '', '', 0, 1, 0, 0, NULL, '2025-06-21 22:48:04.889', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (28, '角色编辑', 'system.roles.update', 'button', 18, '', '', '', 0, 1, 0, 0, NULL, '2025-06-21 22:48:06.248', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (29, '角色删除', 'system.roles.delete', 'button', 18, '', '', '', 0, 1, 0, 0, NULL, '2025-06-21 22:47:51.038', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (30, '角色权限设置', 'system.roles.permission', 'button', 18, '', '', '', 0, 1, 0, 0, NULL, '2025-06-21 22:47:50.209', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (31, '权限查询', 'system.permissions.query', 'button', 19, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (32, '权限创建', 'system.permissions.create', 'button', 19, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (33, '权限编辑', 'system.permissions.update', 'button', 19, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (34, '权限删除', 'system.permissions.delete', 'button', 19, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (35, '任务查询', 'job.list.query', 'button', 4, '', '', '', 0, 1, 0, 0, NULL, '2025-06-21 23:37:59.856', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (36, '任务创建', 'job.list.create', 'button', 4, '', '', '', 0, 1, 0, 0, NULL, '2025-06-21 23:38:00.481', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (37, '任务编辑', 'job.list.update', 'button', 4, '', '', '', 0, 1, 0, 0, NULL, '2025-06-21 23:38:01.120', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (38, '任务删除', 'job.list.delete', 'button', 4, '', '', '', 0, 1, 0, 0, NULL, '2025-06-21 23:38:02.007', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (39, '任务执行', 'job.list.execute', 'button', 4, '', '', '', 0, 1, 0, 0, NULL, '2025-06-21 23:38:03.803', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (40, '历史查询', 'job.history.query', 'button', 5, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (41, '历史详情', 'job.history.detail', 'button', 5, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (42, '终端查询', 'terminal.list.query', 'button', 7, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (43, '终端创建', 'terminal.list.create', 'button', 7, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (44, '终端编辑', 'terminal.list.update', 'button', 7, '', '', '', 0, 1, 0, 0, NULL, '2025-06-21 22:57:59.562', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (45, '终端删除', 'terminal.list.delete', 'button', 7, '', '', '', 0, 1, 0, 0, NULL, '2025-06-21 22:14:20.455', NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (46, '终端连接', 'terminal.list.connect', 'button', 7, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (47, '集群查询', 'tools.kafka.cluster.query', 'button', 13, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (48, '集群创建', 'tools.kafka.cluster.create', 'button', 13, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (49, '集群编辑', 'tools.kafka.cluster.update', 'button', 13, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (50, '集群删除', 'tools.kafka.cluster.delete', 'button', 13, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (51, '操作管理', 'system.logs', 'menu', 16, 'logs', 'views/system/logs/index.vue', 'icon-file', 2, 1, 0, 1, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (56, '终端上传', 'terminal.list.upload', 'button', 7, '', '', '', 0, 1, 0, 0, NULL, NULL, NULL);
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `hidden`, `cache`, `created_at`, `updated_at`, `deleted_at`) VALUES (57, '终端下载', 'terminal.list.download', 'button', 7, '', '', '', 0, 1, 0, 0, NULL, '2025-06-24 11:19:25.955', NULL);
COMMIT;

-- ----------------------------
-- Table structure for role_permissions
-- ----------------------------
DROP TABLE IF EXISTS `role_permissions`;
CREATE TABLE `role_permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `role_id` bigint unsigned NOT NULL,
  `permission_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_role_permissions_deleted_at` (`deleted_at`),
  KEY `idx_role_permissions_role_id` (`role_id`),
  KEY `idx_role_permissions_permission_id` (`permission_id`),
  KEY `idx_role_perm` (`role_id`,`permission_id`)
) ENGINE=InnoDB AUTO_INCREMENT=215 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of role_permissions
-- ----------------------------
BEGIN;
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 1, 1, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 1, 2, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 1, 3, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 1, 4, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 1, 5, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 1, 6, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, 1, 7, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, 1, 8, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, 1, 9, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (10, 1, 10, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (11, 1, 11, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (12, 1, 12, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (13, 1, 13, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (14, 1, 14, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (15, 1, 15, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (16, 1, 16, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (17, 1, 17, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (18, 1, 18, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (19, 1, 19, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (20, 1, 20, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (21, 1, 21, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (22, 1, 22, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (23, 1, 23, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (24, 1, 24, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (25, 1, 25, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (26, 1, 26, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (27, 1, 27, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (28, 1, 28, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (29, 1, 29, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (30, 1, 30, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (31, 1, 31, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (32, 1, 32, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (33, 1, 33, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (34, 1, 34, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (35, 1, 35, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (36, 1, 36, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (37, 1, 37, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (38, 1, 38, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (39, 1, 39, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (40, 1, 40, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (41, 1, 41, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (42, 1, 42, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (43, 1, 43, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (44, 1, 44, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (45, 1, 45, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (46, 1, 46, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (47, 1, 47, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (48, 1, 48, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (49, 1, 49, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (50, 1, 50, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (64, 1, 51, '2025-06-19 22:48:48.138', '2025-06-19 22:48:48.138', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (65, 1, 52, '2025-06-19 22:48:48.139', '2025-06-19 22:48:48.139', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (66, 1, 53, '2025-06-19 22:48:48.141', '2025-06-19 22:48:48.141', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (67, 1, 54, '2025-06-19 22:48:48.141', '2025-06-19 22:48:48.141', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (68, 1, 55, '2025-06-19 22:48:48.142', '2025-06-19 22:48:48.142', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (141, 2, 2, '2025-06-19 23:15:00.343', '2025-06-19 23:15:00.343', '2025-06-19 23:21:56.047');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (142, 2, 3, '2025-06-19 23:15:00.344', '2025-06-19 23:15:00.344', '2025-06-19 23:21:56.047');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (143, 2, 6, '2025-06-19 23:15:00.344', '2025-06-19 23:15:00.344', '2025-06-19 23:21:56.047');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (144, 2, 4, '2025-06-19 23:15:00.344', '2025-06-19 23:15:00.344', '2025-06-19 23:21:56.047');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (145, 2, 7, '2025-06-19 23:15:00.345', '2025-06-19 23:15:00.345', '2025-06-19 23:21:56.047');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (146, 2, 44, '2025-06-19 23:15:00.345', '2025-06-19 23:15:00.345', '2025-06-19 23:21:56.047');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (147, 2, 45, '2025-06-19 23:15:00.347', '2025-06-19 23:15:00.347', '2025-06-19 23:21:56.047');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (148, 2, 8, '2025-06-19 23:15:00.347', '2025-06-19 23:15:00.347', '2025-06-19 23:21:56.047');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (149, 2, 9, '2025-06-19 23:15:00.348', '2025-06-19 23:15:00.348', '2025-06-19 23:21:56.047');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (150, 2, 10, '2025-06-19 23:15:00.349', '2025-06-19 23:15:00.349', '2025-06-19 23:21:56.047');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (151, 2, 11, '2025-06-19 23:15:00.350', '2025-06-19 23:15:00.350', '2025-06-19 23:21:56.047');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (152, 2, 2, '2025-06-19 23:21:56.050', '2025-06-19 23:21:56.050', '2025-06-19 23:44:44.865');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (153, 2, 3, '2025-06-19 23:21:56.051', '2025-06-19 23:21:56.051', '2025-06-19 23:44:44.865');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (154, 2, 9, '2025-06-19 23:21:56.051', '2025-06-19 23:21:56.051', '2025-06-19 23:44:44.865');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (155, 2, 4, '2025-06-19 23:21:56.051', '2025-06-19 23:21:56.051', '2025-06-19 23:44:44.865');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (156, 2, 10, '2025-06-19 23:21:56.052', '2025-06-19 23:21:56.052', '2025-06-19 23:44:44.865');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (157, 2, 11, '2025-06-19 23:21:56.053', '2025-06-19 23:21:56.053', '2025-06-19 23:44:44.865');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (158, 2, 7, '2025-06-19 23:21:56.055', '2025-06-19 23:21:56.055', '2025-06-19 23:44:44.865');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (159, 2, 44, '2025-06-19 23:21:56.059', '2025-06-19 23:21:56.059', '2025-06-19 23:44:44.865');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (160, 2, 45, '2025-06-19 23:21:56.059', '2025-06-19 23:21:56.059', '2025-06-19 23:44:44.865');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (161, 2, 2, '2025-06-19 23:44:44.869', '2025-06-19 23:44:44.869', '2025-06-19 23:45:48.553');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (162, 2, 3, '2025-06-19 23:44:44.873', '2025-06-19 23:44:44.873', '2025-06-19 23:45:48.553');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (163, 2, 9, '2025-06-19 23:44:44.873', '2025-06-19 23:44:44.873', '2025-06-19 23:45:48.553');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (164, 2, 4, '2025-06-19 23:44:44.874', '2025-06-19 23:44:44.874', '2025-06-19 23:45:48.553');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (165, 2, 10, '2025-06-19 23:44:44.874', '2025-06-19 23:44:44.874', '2025-06-19 23:45:48.553');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (166, 2, 11, '2025-06-19 23:44:44.879', '2025-06-19 23:44:44.879', '2025-06-19 23:45:48.553');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (167, 2, 2, '2025-06-19 23:45:48.554', '2025-06-19 23:45:48.554', '2025-06-19 23:52:15.501');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (168, 2, 3, '2025-06-19 23:45:48.555', '2025-06-19 23:45:48.555', '2025-06-19 23:52:15.501');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (169, 2, 9, '2025-06-19 23:45:48.555', '2025-06-19 23:45:48.555', '2025-06-19 23:52:15.501');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (170, 2, 4, '2025-06-19 23:45:48.555', '2025-06-19 23:45:48.555', '2025-06-19 23:52:15.501');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (171, 2, 10, '2025-06-19 23:45:48.556', '2025-06-19 23:45:48.556', '2025-06-19 23:52:15.501');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (172, 2, 11, '2025-06-19 23:45:48.558', '2025-06-19 23:45:48.558', '2025-06-19 23:52:15.501');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (173, 2, 2, '2025-06-19 23:52:15.503', '2025-06-19 23:52:15.503', '2025-06-19 23:53:38.431');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (174, 2, 3, '2025-06-19 23:52:15.503', '2025-06-19 23:52:15.503', '2025-06-19 23:53:38.431');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (175, 2, 9, '2025-06-19 23:52:15.504', '2025-06-19 23:52:15.504', '2025-06-19 23:53:38.431');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (176, 2, 7, '2025-06-19 23:52:15.504', '2025-06-19 23:52:15.504', '2025-06-19 23:53:38.431');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (177, 2, 4, '2025-06-19 23:52:15.504', '2025-06-19 23:52:15.504', '2025-06-19 23:53:38.431');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (178, 2, 10, '2025-06-19 23:52:15.506', '2025-06-19 23:52:15.506', '2025-06-19 23:53:38.431');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (179, 2, 11, '2025-06-19 23:52:15.506', '2025-06-19 23:52:15.506', '2025-06-19 23:53:38.431');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (180, 2, 44, '2025-06-19 23:52:15.507', '2025-06-19 23:52:15.507', '2025-06-19 23:53:38.431');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (181, 2, 45, '2025-06-19 23:52:15.507', '2025-06-19 23:52:15.507', '2025-06-19 23:53:38.431');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (182, 2, 8, '2025-06-19 23:52:15.507', '2025-06-19 23:52:15.507', '2025-06-19 23:53:38.431');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (183, 2, 6, '2025-06-19 23:52:15.507', '2025-06-19 23:52:15.507', '2025-06-19 23:53:38.431');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (184, 2, 2, '2025-06-19 23:53:38.433', '2025-06-19 23:53:38.433', '2025-06-19 23:54:05.225');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (185, 2, 3, '2025-06-19 23:53:38.439', '2025-06-19 23:53:38.439', '2025-06-19 23:54:05.225');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (186, 2, 9, '2025-06-19 23:53:38.439', '2025-06-19 23:53:38.439', '2025-06-19 23:54:05.225');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (187, 2, 7, '2025-06-19 23:53:38.440', '2025-06-19 23:53:38.440', '2025-06-19 23:54:05.225');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (188, 2, 4, '2025-06-19 23:53:38.440', '2025-06-19 23:53:38.440', '2025-06-19 23:54:05.225');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (189, 2, 10, '2025-06-19 23:53:38.441', '2025-06-19 23:53:38.441', '2025-06-19 23:54:05.225');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (190, 2, 11, '2025-06-19 23:53:38.441', '2025-06-19 23:53:38.441', '2025-06-19 23:54:05.225');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (191, 2, 44, '2025-06-19 23:53:38.442', '2025-06-19 23:53:38.442', '2025-06-19 23:54:05.225');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (192, 2, 45, '2025-06-19 23:53:38.442', '2025-06-19 23:53:38.442', '2025-06-19 23:54:05.225');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (193, 2, 8, '2025-06-19 23:53:38.442', '2025-06-19 23:53:38.442', '2025-06-19 23:54:05.225');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (194, 2, 6, '2025-06-19 23:53:38.443', '2025-06-19 23:53:38.443', '2025-06-19 23:54:05.225');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (195, 2, 2, '2025-06-19 23:54:05.227', '2025-06-19 23:54:05.227', '2025-06-19 23:57:51.101');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (196, 2, 3, '2025-06-19 23:54:05.228', '2025-06-19 23:54:05.228', '2025-06-19 23:57:51.101');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (197, 2, 9, '2025-06-19 23:54:05.229', '2025-06-19 23:54:05.229', '2025-06-19 23:57:51.101');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (198, 2, 4, '2025-06-19 23:54:05.229', '2025-06-19 23:54:05.229', '2025-06-19 23:57:51.101');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (199, 2, 10, '2025-06-19 23:54:05.230', '2025-06-19 23:54:05.230', '2025-06-19 23:57:51.101');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (200, 2, 11, '2025-06-19 23:54:05.231', '2025-06-19 23:54:05.231', '2025-06-19 23:57:51.101');
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (201, 2, 2, '2025-06-19 23:57:51.105', '2025-06-19 23:57:51.105', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (202, 2, 3, '2025-06-19 23:57:51.110', '2025-06-19 23:57:51.110', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (203, 2, 9, '2025-06-19 23:57:51.111', '2025-06-19 23:57:51.111', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (204, 2, 4, '2025-06-19 23:57:51.113', '2025-06-19 23:57:51.113', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (205, 2, 10, '2025-06-19 23:57:51.116', '2025-06-19 23:57:51.116', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (206, 2, 11, '2025-06-19 23:57:51.117', '2025-06-19 23:57:51.117', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (207, 2, 12, '2025-06-19 23:57:51.117', '2025-06-19 23:57:51.117', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (208, 2, 13, '2025-06-19 23:57:51.117', '2025-06-19 23:57:51.117', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (209, 2, 14, '2025-06-19 23:57:51.118', '2025-06-19 23:57:51.118', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (210, 2, 15, '2025-06-19 23:57:51.118', '2025-06-19 23:57:51.118', NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (211, 1, 56, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (212, 1, 57, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (213, 2, 56, NULL, NULL, NULL);
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (214, 2, 57, NULL, NULL, NULL);
COMMIT;

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `description` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `status` bigint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `is_admin` tinyint(1) DEFAULT '0',
  `remark` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_roles_code` (`code`),
  UNIQUE KEY `uni_roles_name` (`name`),
  KEY `idx_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of roles
-- ----------------------------
BEGIN;
INSERT INTO `roles` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`, `is_admin`, `remark`) VALUES (1, '超级管理员', 'admin', '系统超级管理员', 1, '2025-01-05 23:43:34.342', '2025-01-06 11:00:04.362', NULL, 0, NULL);
INSERT INTO `roles` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`, `is_admin`, `remark`) VALUES (2, '普通用户', 'user', '普通用户角色', 1, '2025-01-05 23:44:22.952', '2025-06-19 22:00:22.208', NULL, 0, NULL);
INSERT INTO `roles` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`, `is_admin`, `remark`) VALUES (3, '测试', 'test', '测试', 1, '2025-01-06 14:19:45.851', '2025-01-06 14:19:45.851', '2025-01-06 14:19:49.046', 0, NULL);
COMMIT;

-- ----------------------------
-- Table structure for terminals
-- ----------------------------
DROP TABLE IF EXISTS `terminals`;
CREATE TABLE `terminals` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '终端名称',
  `host` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '主机地址',
  `port` int NOT NULL DEFAULT '22' COMMENT 'SSH端口',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名',
  `auth_type` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'password' COMMENT '认证类型(password/key)',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '密码',
  `key_file` text COLLATE utf8mb4_general_ci COMMENT '密钥文件内容',
  `key_passphrase` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '密钥文件密码',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'offline' COMMENT '状态(online/offline)',
  `last_seen` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后在线时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_terminals_host_username` (`host`,`username`),
  KEY `idx_terminals_deleted_at` (`deleted_at`),
  KEY `idx_terminals_name` (`name`),
  KEY `idx_terminals_host` (`host`),
  KEY `idx_terminals_status` (`status`),
  CONSTRAINT `terminals_chk_1` CHECK (((`port` >= 1) and (`port` <= 65535))),
  CONSTRAINT `terminals_chk_2` CHECK ((`status` in (_utf8mb4'online',_utf8mb4'offline')))
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='终端管理';

-- ----------------------------
-- Records of terminals
-- ----------------------------


-- ----------------------------
-- Table structure for user_roles
-- ----------------------------
DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `role_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_roles_user_id` (`user_id`),
  KEY `idx_user_roles_role_id` (`role_id`),
  KEY `idx_user_roles_deleted_at` (`deleted_at`),
  KEY `idx_user_role` (`user_id`,`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of user_roles
-- ----------------------------
BEGIN;
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 1, 1, '2025-01-05 23:43:34.345', '2025-01-05 23:43:34.345', '2025-01-11 23:30:38.663');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 2, 2, '2025-01-05 23:44:55.430', '2025-01-05 23:44:55.430', '2025-01-11 18:15:10.858');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 3, 1, '2025-01-07 14:51:11.761', '2025-01-07 14:51:11.761', '2025-01-07 14:52:22.620');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 3, 2, '2025-01-07 14:51:11.776', '2025-01-07 14:51:11.776', '2025-01-07 14:52:22.620');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 3, 2, '2025-01-07 14:52:22.625', '2025-01-07 14:52:22.625', '2025-01-07 14:53:50.263');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 3, 2, '2025-01-07 14:53:59.105', '2025-01-07 14:53:59.105', '2025-01-07 14:55:57.998');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, 3, 2, '2025-01-07 14:55:58.017', '2025-01-07 14:55:58.017', '2025-01-07 14:56:23.882');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, 3, 2, '2025-01-07 14:56:23.887', '2025-01-07 14:56:23.887', '2025-01-07 15:04:16.042');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, 3, 1, '2025-01-07 14:56:23.892', '2025-01-07 14:56:23.892', '2025-01-07 15:04:16.042');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (10, 3, 2, '2025-01-07 15:04:16.045', '2025-01-07 15:04:16.045', '2025-01-11 23:34:20.976');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (11, 2, 2, '2025-01-11 18:15:10.869', '2025-01-11 18:15:10.869', '2025-01-11 18:15:54.796');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (12, 2, 1, '2025-01-11 18:15:10.873', '2025-01-11 18:15:10.873', '2025-01-11 18:15:54.796');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (13, 2, 2, '2025-01-11 18:15:54.800', '2025-01-11 18:15:54.800', '2025-02-28 15:07:14.948');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (14, 1, 1, '2025-01-11 23:30:38.668', '2025-01-11 23:30:38.668', NULL);
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (15, 4, 2, '2025-01-11 23:34:11.146', '2025-01-11 23:34:11.146', '2025-01-23 15:19:45.807');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (16, 3, 2, '2025-01-11 23:34:20.982', '2025-01-11 23:34:20.982', '2025-01-11 23:56:23.264');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (17, 3, 1, '2025-01-11 23:34:20.987', '2025-01-11 23:34:20.987', '2025-01-11 23:56:23.264');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (18, 3, 2, '2025-01-11 23:56:23.275', '2025-01-11 23:56:23.275', '2025-01-11 23:56:29.876');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (19, 3, 2, '2025-01-11 23:56:29.877', '2025-01-11 23:56:29.877', '2025-02-28 15:07:10.530');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (20, 3, 1, '2025-01-11 23:56:29.878', '2025-01-11 23:56:29.878', '2025-02-28 15:07:10.530');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (21, 7, 2, '2025-05-24 21:52:12.091', '2025-05-24 21:52:12.091', NULL);
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (22, 8, 2, '2025-05-25 10:25:46.870', '2025-05-25 10:25:46.870', NULL);
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (23, 9, 2, '2025-06-19 22:42:00.908', '2025-06-19 22:42:00.908', NULL);
COMMIT;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名',
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '昵称',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '邮箱',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '头像',
  `status` bigint DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of users
-- ----------------------------
BEGIN;
INSERT INTO `users` (`id`, `username`, `password`, `nickname`, `email`, `avatar`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 'admin', '$2a$10$rGo.7OK8UTGrS8sHtv9gw.SDpHDVOQgGfl/qKuQHaigqQMxQmwmyW', '超级管理员', 'admin@example.com', '', 1, '2025-01-05 23:35:27.267', '2025-06-10 14:27:02.105', NULL);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
