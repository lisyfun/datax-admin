/*
 Navicat Premium Dump SQL

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80033 (8.0.33)
 Source Host           : localhost:13306
 Source Schema         : datax_admin

 Target Server Type    : MySQL
 Target Server Version : 80033 (8.0.33)
 File Encoding         : 65001

 Date: 17/01/2025 12:52:23
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for job_histories
-- ----------------------------
DROP TABLE IF EXISTS `job_histories`;
CREATE TABLE `job_histories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `job_id` bigint unsigned NOT NULL,
  `status` bigint DEFAULT NULL,
  `start_time` datetime(3) DEFAULT NULL,
  `end_time` datetime(3) DEFAULT NULL,
  `duration` bigint NOT NULL,
  `output` longtext COLLATE utf8mb4_unicode_ci,
  `error` longtext COLLATE utf8mb4_unicode_ci,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_job_id` (`job_id`)
) ENGINE=InnoDB AUTO_INCREMENT=330 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of job_histories
-- ----------------------------
BEGIN;
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (279, 2, 1, '2025-01-14 09:10:00.010', '2025-01-14 09:10:00.012', 2, '{\"message\":\"pong\"}', '', '2025-01-14 09:10:00.013', '2025-01-14 09:10:00.013', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (280, 2, 1, '2025-01-14 09:15:00.002', '2025-01-14 09:15:00.009', 6, '{\"message\":\"pong\"}', '', '2025-01-14 09:15:00.014', '2025-01-14 09:15:00.014', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (281, 2, 1, '2025-01-14 09:36:30.820', '2025-01-14 09:36:30.823', 3, '{\"message\":\"pong\"}', '', '2025-01-14 09:36:30.871', '2025-01-14 09:36:30.871', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (282, 2, 1, '2025-01-14 09:40:00.012', '2025-01-14 09:40:00.015', 2, '{\"message\":\"pong\"}', '', '2025-01-14 09:40:00.066', '2025-01-14 09:40:00.066', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (283, 2, 1, '2025-01-14 09:45:00.012', '2025-01-14 09:45:00.016', 3, '{\"message\":\"pong\"}', '', '2025-01-14 09:45:00.031', '2025-01-14 09:45:00.031', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (284, 2, 1, '2025-01-14 09:50:00.010', '2025-01-14 09:50:00.013', 3, '{\"message\":\"pong\"}', '', '2025-01-14 09:50:00.066', '2025-01-14 09:50:00.066', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (285, 2, 1, '2025-01-14 10:10:41.804', '2025-01-14 10:10:41.808', 4, '{\"message\":\"pong\"}', '', '2025-01-14 10:10:41.854', '2025-01-14 10:10:41.854', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (286, 2, 1, '2025-01-14 10:15:00.009', '2025-01-14 10:15:00.012', 2, '{\"message\":\"pong\"}', '', '2025-01-14 10:15:00.067', '2025-01-14 10:15:00.067', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (287, 2, 1, '2025-01-14 10:20:00.007', '2025-01-14 10:20:00.008', 1, '{\"message\":\"pong\"}', '', '2025-01-14 10:20:00.041', '2025-01-14 10:20:00.041', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (288, 2, 1, '2025-01-14 10:25:00.008', '2025-01-14 10:25:00.011', 3, '{\"message\":\"pong\"}', '', '2025-01-14 10:25:00.070', '2025-01-14 10:25:00.070', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (289, 2, 1, '2025-01-14 10:30:00.007', '2025-01-14 10:30:00.010', 3, '{\"message\":\"pong\"}', '', '2025-01-14 10:30:00.072', '2025-01-14 10:30:00.072', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (290, 2, 1, '2025-01-14 10:35:00.007', '2025-01-14 10:35:00.009', 2, '{\"message\":\"pong\"}', '', '2025-01-14 10:35:00.055', '2025-01-14 10:35:00.055', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (291, 2, 1, '2025-01-14 10:40:00.005', '2025-01-14 10:40:00.007', 1, '{\"message\":\"pong\"}', '', '2025-01-14 10:40:00.019', '2025-01-14 10:40:00.019', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (292, 2, 1, '2025-01-14 10:45:00.009', '2025-01-14 10:45:00.013', 3, '{\"message\":\"pong\"}', '', '2025-01-14 10:45:00.068', '2025-01-14 10:45:00.068', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (293, 2, 1, '2025-01-14 10:50:00.009', '2025-01-14 10:50:00.013', 3, '{\"message\":\"pong\"}', '', '2025-01-14 10:50:00.060', '2025-01-14 10:50:00.060', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (294, 2, 1, '2025-01-14 10:55:00.004', '2025-01-14 10:55:00.006', 2, '{\"message\":\"pong\"}', '', '2025-01-14 10:55:00.034', '2025-01-14 10:55:00.034', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (295, 2, 1, '2025-01-14 11:00:00.008', '2025-01-14 11:00:00.011', 2, '{\"message\":\"pong\"}', '', '2025-01-14 11:00:00.053', '2025-01-14 11:00:00.053', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (296, 2, 1, '2025-01-14 11:05:00.009', '2025-01-14 11:05:00.012', 2, '{\"message\":\"pong\"}', '', '2025-01-14 11:05:00.055', '2025-01-14 11:05:00.055', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (297, 2, 1, '2025-01-14 11:26:01.508', '2025-01-14 11:26:01.511', 3, '{\"message\":\"pong\"}', '', '2025-01-14 11:26:01.561', '2025-01-14 11:26:01.561', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (298, 2, 1, '2025-01-14 11:30:00.000', '2025-01-14 11:30:00.011', 10, '{\"message\":\"pong\"}', '', '2025-01-14 11:30:00.056', '2025-01-14 11:30:00.056', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (299, 2, 1, '2025-01-14 11:35:00.005', '2025-01-14 11:35:00.010', 4, '{\"message\":\"pong\"}', '', '2025-01-14 11:35:00.071', '2025-01-14 11:35:00.071', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (300, 2, 1, '2025-01-14 11:40:00.011', '2025-01-14 11:40:00.013', 2, '{\"message\":\"pong\"}', '', '2025-01-14 11:40:00.058', '2025-01-14 11:40:00.058', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (301, 2, 1, '2025-01-14 11:45:00.009', '2025-01-14 11:45:00.012', 2, '{\"message\":\"pong\"}', '', '2025-01-14 11:45:00.068', '2025-01-14 11:45:00.068', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (302, 2, 1, '2025-01-14 11:50:00.007', '2025-01-14 11:50:00.010', 3, '{\"message\":\"pong\"}', '', '2025-01-14 11:50:00.071', '2025-01-14 11:50:00.071', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (303, 2, 1, '2025-01-14 11:55:00.009', '2025-01-14 11:55:00.012', 2, '{\"message\":\"pong\"}', '', '2025-01-14 11:55:00.062', '2025-01-14 11:55:00.062', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (304, 2, 1, '2025-01-14 12:00:00.008', '2025-01-14 12:00:00.012', 3, '{\"message\":\"pong\"}', '', '2025-01-14 12:00:00.047', '2025-01-14 12:00:00.047', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (305, 2, 1, '2025-01-14 12:05:00.003', '2025-01-14 12:05:00.005', 1, '{\"message\":\"pong\"}', '', '2025-01-14 12:05:00.033', '2025-01-14 12:05:00.033', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (306, 2, 1, '2025-01-14 12:10:00.000', '2025-01-14 12:10:00.012', 11, '{\"message\":\"pong\"}', '', '2025-01-14 12:10:00.067', '2025-01-14 12:10:00.067', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (307, 2, 1, '2025-01-14 12:15:00.008', '2025-01-14 12:15:00.012', 3, '{\"message\":\"pong\"}', '', '2025-01-14 12:15:00.059', '2025-01-14 12:15:00.059', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (308, 2, 1, '2025-01-14 12:20:00.009', '2025-01-14 12:20:00.012', 3, '{\"message\":\"pong\"}', '', '2025-01-14 12:20:00.062', '2025-01-14 12:20:00.062', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (309, 2, 1, '2025-01-14 12:25:00.007', '2025-01-14 12:25:00.011', 4, '{\"message\":\"pong\"}', '', '2025-01-14 12:25:00.070', '2025-01-14 12:25:00.070', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (310, 2, 1, '2025-01-14 12:30:00.005', '2025-01-14 12:30:00.008', 2, '{\"message\":\"pong\"}', '', '2025-01-14 12:30:00.036', '2025-01-14 12:30:00.036', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (311, 2, 1, '2025-01-14 12:35:00.002', '2025-01-14 12:35:00.007', 5, '{\"message\":\"pong\"}', '', '2025-01-14 12:35:00.040', '2025-01-14 12:35:00.040', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (312, 2, 1, '2025-01-14 12:40:00.007', '2025-01-14 12:40:00.010', 3, '{\"message\":\"pong\"}', '', '2025-01-14 12:40:00.060', '2025-01-14 12:40:00.060', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (313, 2, 1, '2025-01-14 12:45:00.011', '2025-01-14 12:45:00.013', 2, '{\"message\":\"pong\"}', '', '2025-01-14 12:45:00.048', '2025-01-14 12:45:00.048', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (314, 2, 1, '2025-01-14 12:50:00.006', '2025-01-14 12:50:00.009', 2, '{\"message\":\"pong\"}', '', '2025-01-14 12:50:00.068', '2025-01-14 12:50:00.068', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (315, 2, 1, '2025-01-14 12:55:00.006', '2025-01-14 12:55:00.010', 3, '{\"message\":\"pong\"}', '', '2025-01-14 12:55:00.063', '2025-01-14 12:55:00.063', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (316, 2, 1, '2025-01-14 13:00:00.008', '2025-01-14 13:00:00.013', 4, '{\"message\":\"pong\"}', '', '2025-01-14 13:00:00.069', '2025-01-14 13:00:00.069', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (317, 2, 1, '2025-01-14 13:05:00.002', '2025-01-14 13:05:00.007', 5, '{\"message\":\"pong\"}', '', '2025-01-14 13:05:00.029', '2025-01-14 13:05:00.029', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (318, 2, 1, '2025-01-14 13:10:00.004', '2025-01-14 13:10:00.006', 2, '{\"message\":\"pong\"}', '', '2025-01-14 13:10:00.052', '2025-01-14 13:10:00.052', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (319, 2, 1, '2025-01-14 13:15:00.007', '2025-01-14 13:15:00.012', 4, '{\"message\":\"pong\"}', '', '2025-01-14 13:15:00.062', '2025-01-14 13:15:00.062', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (320, 2, 1, '2025-01-14 13:20:00.008', '2025-01-14 13:20:00.012', 4, '{\"message\":\"pong\"}', '', '2025-01-14 13:20:00.084', '2025-01-14 13:20:00.084', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (321, 2, 1, '2025-01-14 13:25:00.008', '2025-01-14 13:25:00.012', 3, '{\"message\":\"pong\"}', '', '2025-01-14 13:25:00.061', '2025-01-14 13:25:00.061', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (322, 2, 1, '2025-01-14 13:30:00.010', '2025-01-14 13:30:00.013', 3, '{\"message\":\"pong\"}', '', '2025-01-14 13:30:00.059', '2025-01-14 13:30:00.059', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (323, 2, 1, '2025-01-14 13:35:00.008', '2025-01-14 13:35:00.011', 3, '{\"message\":\"pong\"}', '', '2025-01-14 13:35:00.012', '2025-01-14 13:35:00.012', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (324, 2, 1, '2025-01-14 13:40:00.007', '2025-01-14 13:40:00.011', 4, '{\"message\":\"pong\"}', '', '2025-01-14 13:40:00.063', '2025-01-14 13:40:00.063', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (325, 9, 1, '2025-01-14 13:45:46.651', '2025-01-14 13:45:46.663', 11, '命令信息:\n执行命令: /app/bin/datax -job /Users/lisongyu/workspace/go/src/datax-admin/backend/datax-1757582301.json\n工作目录: /app\n配置内容:\n{\n    \"job\": {\n        \"content\": [\n            {\n                \"reader\": {\n                    \"name\": \"oraclereader\",\n                    \"parameter\": {\n                        \"username\": \"sys\",\n                        \"password\": \"123456\",\n                        \"host\": \"127.0.0.1\",\n                        \"port\": 1521,\n                        \"service\": \"ORCL\",\n                        \"table\": \"USERS\",\n                        \"columns\": [],\n                        \"selectSql\": \"select ID, NAME, EMAIL, AGE from USERS\",\n                        \"where\": \"id > 0\",\n                        \"batchSize\": 20000\n                    }\n                },\n                \"writer\": {\n                    \"name\": \"mysqlwriter\",\n                    \"parameter\": {\n                        \"username\": \"root\",\n                        \"password\": \"123456\",\n                        \"host\": \"127.0.0.1\",\n                        \"port\": 13306,\n                        \"database\": \"datax_target\",\n                        \"table\": \"users\",\n                        \"columns\": [\n                            \"id\",\n                            \"username\",\n                            \"email\",\n                            \"age\"\n                        ],\n                        \"preSql\": [\n                            \"select count(*) from users\",\n                            \"truncate users\"\n                        ],\n                        \"postSql\": [\n                            \"select count(*) from users\"\n                        ],\n                        \"batchSize\": 1000,\n                        \"writeMode\": \"insert\"\n                    }\n                }\n            }\n        ],\n        \"setting\": {\n            \"speed\": {\n                \"channel\": 5\n            }\n        }\n    }\n}\n\n无输出', '', '2025-01-14 13:45:46.683', '2025-01-14 13:45:46.683', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (326, 8, 1, '2025-01-14 13:45:46.664', '2025-01-14 13:45:46.675', 10, '命令信息:\n执行命令: /app/bin/datax -job /Users/lisongyu/workspace/go/src/datax-admin/backend/datax-196421810.json\n工作目录: /app\n配置内容:\n{\n    \"job\": {\n        \"content\": [\n            {\n                \"reader\": {\n                    \"name\": \"mysqlreader\",\n                    \"parameter\": {\n                        \"username\": \"root\",\n                        \"password\": \"123456\",\n                        \"host\": \"localhost\",\n                        \"port\": 13306,\n                        \"database\": \"datax_source1\",\n                        \"table\": \"users\",\n                        \"columns\": [\n                            \"username\",\n                            \"email\",\n                            \"age\",\n                            \"created_at\",\n                            \"updated_at\"\n                        ],\n                        \"selectSql\": \"select username, email, age, created_at, updated_at from users\",\n                        \"where\": \"1=1\",\n                        \"batchSize\": 20000\n                    }\n                },\n                \"writer\": {\n                    \"name\": \"postgresqlwriter\",\n                    \"parameter\": {\n                        \"username\": \"postgres\",\n                        \"password\": \"123456\",\n                        \"host\": \"localhost\",\n                        \"port\": 15432,\n                        \"database\": \"target_db\",\n                        \"schema\": \"public\",\n                        \"table\": \"users\",\n                        \"columns\": [\n                            \"username\",\n                            \"email\",\n                            \"age\",\n                            \"created_at\",\n                            \"updated_at\"\n                        ],\n                        \"batchSize\": 20000,\n                        \"preSql\": [\n                            \"select count(*) as total_count from users\",\n                            \"truncate table users\"\n                        ],\n                        \"postSql\": [\n                            \"select count(*) as total_count from users\"\n                        ],\n                        \"writeMode\": \"insert\"\n                    }\n                }\n            }\n        ],\n        \"setting\": {\n            \"speed\": {\n                \"channel\": 24,\n                \"bytes\": 52428800\n            },\n            \"errorLimit\": {\n                \"record\": 0,\n                \"percentage\": 0.02\n            }\n        }\n    }\n}\n\n无输出', '', '2025-01-14 13:45:46.686', '2025-01-14 13:45:46.686', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (327, 4, 1, '2025-01-14 13:45:46.668', '2025-01-14 13:45:46.679', 11, '1\n', '', '2025-01-14 13:45:46.697', '2025-01-14 13:45:46.697', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (328, 7, 1, '2025-01-14 13:45:46.667', '2025-01-14 13:45:46.675', 8, '命令信息:\n执行命令: /app/bin/datax -job /Users/lisongyu/workspace/go/src/datax-admin/backend/datax-3856738629.json\n工作目录: /app\n配置内容:\n{\n    \"job\": {\n        \"content\": [\n            {\n                \"reader\": {\n                    \"name\": \"mysqlreader\",\n                    \"parameter\": {\n                        \"username\": \"root\",\n                        \"password\": \"123456\",\n                        \"host\": \"192.168.228.1\",\n                        \"port\": 13306,\n                        \"database\": \"datax_source\",\n                        \"table\": \"users\",\n                        \"columns\": [],\n                        \"selectSql\": \"select username, email, age, created_at, updated_at from users\",\n                        \"where\": \"1=1\",\n                        \"batchSize\": 20000\n                    }\n                },\n                \"writer\": {\n                    \"name\": \"mysqlwriter\",\n                    \"parameter\": {\n                        \"username\": \"root\",\n                        \"password\": \"123456\",\n                        \"host\": \"192.168.228.1\",\n                        \"port\": 13306,\n                        \"database\": \"datax_target\",\n                        \"table\": \"users\",\n                        \"columns\": [\n                            \"username\",\n                            \"email\",\n                            \"age\",\n                            \"created_at\",\n                            \"updated_at\"\n                        ],\n                        \"batchSize\": 10000,\n                        \"preSql\": [\n                            \"select count(*) from users\",\n                            \"truncate users\"\n                        ],\n                        \"postSql\": [\n                            \"select count(*) from users\"\n                        ],\n                        \"writeMode\": \"replace\"\n                    }\n                }\n            }\n        ],\n        \"setting\": {\n            \"speed\": {\n                \"channel\": 24,\n                \"bytes\": 52428800\n            },\n            \"errorLimit\": {\n                \"record\": 0,\n                \"percentage\": 0.02\n            }\n        }\n    }\n}\n\n无输出', '', '2025-01-14 13:45:46.697', '2025-01-14 13:45:46.697', NULL);
INSERT INTO `job_histories` (`id`, `job_id`, `status`, `start_time`, `end_time`, `duration`, `output`, `error`, `created_at`, `updated_at`, `deleted_at`) VALUES (329, 2, 1, '2025-01-14 13:45:46.669', '2025-01-14 13:45:46.681', 11, '{\"message\":\"pong\"}', '', '2025-01-14 13:45:46.697', '2025-01-14 13:45:46.697', NULL);
COMMIT;

-- ----------------------------
-- Table structure for jobs
-- ----------------------------
DROP TABLE IF EXISTS `jobs`;
CREATE TABLE `jobs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cron_expr` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint DEFAULT '0',
  `timeout` int DEFAULT '0',
  `retry_count` int DEFAULT '0',
  `retry_delay` int DEFAULT '0',
  `params` text COLLATE utf8mb4_unicode_ci,
  `creator` bigint unsigned NOT NULL,
  `updater` bigint unsigned DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_jobs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of jobs
-- ----------------------------
BEGIN;
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, '数据库备份', 'shell', '每天凌晨3点备份MySQL数据库', '0 0 3 * * *', 0, 3600, 3, 300, '{\"command\":\"mysqldump -u root -p123456 mydb \\u003e /backup/mydb_$(date +%Y%m%d).sql\",\"work_dir\":\"/backup\",\"environment\":{\"MYSQL_PWD\":\"123456\"}}', 1, 0, '2025-01-06 13:19:27.000', '2025-01-10 12:30:43.000', '2025-01-10 12:30:44.000');
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, '服务健康检查', 'http', '每5分钟检查一次服务健康状态', '0 0/5 * * * ?', 0, 30, 3, 60, '{\"url\":\"http://localhost:8080/api/v1/ping\",\"method\":\"GET\",\"headers\":{},\"body\":\"\",\"success_code\":[200]}', 1, 1, '2025-01-06 13:19:33.000', '2025-01-14 13:47:22.117', NULL);
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, '数据同步', 'datax', '每小时同步一次数据', '0 0 * * * *', 0, 7200, 2, 600, '{\"job_path\":\"/datax/job/mysql_to_mysql.json\",\"parameters\":{\"batch_size\":\"1000\",\"from_date\":\"${yyyy-MM-dd}\"},\"jvm_options\":[\"-Xms1g\",\"-Xmx2g\"],\"speed\":10000}', 1, 0, '2025-01-06 13:19:42.000', '2025-01-10 12:24:03.000', '2025-01-10 12:24:04.000');
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, '测试任务', 'shell', '这是一个测试任务', '* 0/1 * * * ?', 0, 3600, 3, 60, '{\"command\":\"echo 1\",\"work_dir\":\"\",\"environment\":{}}', 1, 1, '2025-01-06 16:57:31.000', '2025-01-12 12:01:29.713', NULL);
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 'test', 'shell', 'test', '* * * * * ?', 0, 2, 2, 34, '{\"command\":\"echo 1\",\"work_dir\":\"\",\"environment\":{}}', 1, 1, '2025-01-06 17:27:00.000', '2025-01-06 17:54:40.000', '2025-01-06 17:54:41.000');
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 'HTTP任务', 'http', '测试HTTP任务', '0 2 1 3 * ?', 0, 30, 3, 5, '{\"url\":\"http://localhost:8080/api/v1/ping\",\"method\":\"GET\",\"headers\":{\"Content-Type\":\"application/json\"},\"body\":\"\",\"success_code\":[200,201]}', 1, 1, '2025-01-06 17:34:12.000', '2025-01-10 12:30:41.000', '2025-01-10 12:30:41.000');
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, 'mysql到 mysql', 'datax', 'mysql到 mysql', '0 0 0 * * ?', 0, 0, 0, 0, '{\"job_config\":\"{\\\"job\\\":{\\\"content\\\":[{\\\"reader\\\":{\\\"name\\\":\\\"mysqlreader\\\",\\\"parameter\\\":{\\\"username\\\":\\\"root\\\",\\\"password\\\":\\\"123456\\\",\\\"host\\\":\\\"192.168.228.1\\\",\\\"port\\\":13306,\\\"database\\\":\\\"datax_source\\\",\\\"table\\\":\\\"users\\\",\\\"columns\\\":[],\\\"selectSql\\\":\\\"select username, email, age, created_at, updated_at from users\\\",\\\"where\\\":\\\"1=1\\\",\\\"batchSize\\\":20000}},\\\"writer\\\":{\\\"name\\\":\\\"mysqlwriter\\\",\\\"parameter\\\":{\\\"username\\\":\\\"root\\\",\\\"password\\\":\\\"123456\\\",\\\"host\\\":\\\"192.168.228.1\\\",\\\"port\\\":13306,\\\"database\\\":\\\"datax_target\\\",\\\"table\\\":\\\"users\\\",\\\"columns\\\":[\\\"username\\\",\\\"email\\\",\\\"age\\\",\\\"created_at\\\",\\\"updated_at\\\"],\\\"batchSize\\\":10000,\\\"preSql\\\":[\\\"select count(*) from users\\\",\\\"truncate users\\\"],\\\"postSql\\\":[\\\"select count(*) from users\\\"],\\\"writeMode\\\":\\\"replace\\\"}}}],\\\"setting\\\":{\\\"speed\\\":{\\\"channel\\\":24,\\\"bytes\\\":52428800},\\\"errorLimit\\\":{\\\"record\\\":0,\\\"percentage\\\":0.02}}}}\",\"parameters\":{\"datax.job.setting.speed.channel\":\"5\"}}', 1, 1, NULL, '2025-01-13 07:22:28.722', NULL);
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, 'MySQL到PostgreSQL同步', 'datax', '从源MySQL数据库同步users表到目标PostgreSQL数据库', '0 0 * * * *', 0, 3600, 3, 300, '{\"job_config\":\"{\\\"job\\\":{\\\"content\\\":[{\\\"reader\\\":{\\\"name\\\":\\\"mysqlreader\\\",\\\"parameter\\\":{\\\"username\\\":\\\"root\\\",\\\"password\\\":\\\"123456\\\",\\\"host\\\":\\\"localhost\\\",\\\"port\\\":13306,\\\"database\\\":\\\"datax_source1\\\",\\\"table\\\":\\\"users\\\",\\\"columns\\\":[\\\"username\\\",\\\"email\\\",\\\"age\\\",\\\"created_at\\\",\\\"updated_at\\\"],\\\"selectSql\\\":\\\"select username, email, age, created_at, updated_at from users\\\",\\\"where\\\":\\\"1=1\\\",\\\"batchSize\\\":20000}},\\\"writer\\\":{\\\"name\\\":\\\"postgresqlwriter\\\",\\\"parameter\\\":{\\\"username\\\":\\\"postgres\\\",\\\"password\\\":\\\"123456\\\",\\\"host\\\":\\\"localhost\\\",\\\"port\\\":15432,\\\"database\\\":\\\"target_db\\\",\\\"schema\\\":\\\"public\\\",\\\"table\\\":\\\"users\\\",\\\"columns\\\":[\\\"username\\\",\\\"email\\\",\\\"age\\\",\\\"created_at\\\",\\\"updated_at\\\"],\\\"batchSize\\\":20000,\\\"preSql\\\":[\\\"select count(*) as total_count from users\\\",\\\"truncate table users\\\"],\\\"postSql\\\":[\\\"select count(*) as total_count from users\\\"],\\\"writeMode\\\":\\\"insert\\\"}}}],\\\"setting\\\":{\\\"speed\\\":{\\\"channel\\\":24,\\\"bytes\\\":52428800},\\\"errorLimit\\\":{\\\"record\\\":0,\\\"percentage\\\":0.02}}}}\",\"parameters\":{}}', 1, 1, '2025-01-10 11:35:20.000', '2025-01-11 14:19:53.000', NULL);
INSERT INTO `jobs` (`id`, `name`, `type`, `description`, `cron_expr`, `status`, `timeout`, `retry_count`, `retry_delay`, `params`, `creator`, `updater`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, 'Oracle到MySQL同步', 'datax', '从源Oracle数据库同步users表到目标MySQL数据库', '0 0 * * * *', 0, 3600, 3, 300, '{\"job_config\":\"{\\\"job\\\":{\\\"content\\\":[{\\\"reader\\\":{\\\"name\\\":\\\"oraclereader\\\",\\\"parameter\\\":{\\\"username\\\":\\\"sys\\\",\\\"password\\\":\\\"123456\\\",\\\"host\\\":\\\"127.0.0.1\\\",\\\"port\\\":1521,\\\"service\\\":\\\"ORCL\\\",\\\"table\\\":\\\"USERS\\\",\\\"columns\\\":[],\\\"selectSql\\\":\\\"select ID, NAME, EMAIL, AGE from USERS\\\",\\\"where\\\":\\\"id \\u003e 0\\\",\\\"batchSize\\\":20000}},\\\"writer\\\":{\\\"name\\\":\\\"mysqlwriter\\\",\\\"parameter\\\":{\\\"username\\\":\\\"root\\\",\\\"password\\\":\\\"123456\\\",\\\"host\\\":\\\"127.0.0.1\\\",\\\"port\\\":13306,\\\"database\\\":\\\"datax_target\\\",\\\"table\\\":\\\"users\\\",\\\"columns\\\":[\\\"id\\\",\\\"username\\\",\\\"email\\\",\\\"age\\\"],\\\"preSql\\\":[\\\"select count(*) from users\\\",\\\"truncate users\\\"],\\\"postSql\\\":[\\\"select count(*) from users\\\"],\\\"batchSize\\\":1000,\\\"writeMode\\\":\\\"insert\\\"}}}],\\\"setting\\\":{\\\"speed\\\":{\\\"channel\\\":5}}}}\",\"parameters\":{\"datax.job.setting.speed.channel\":\"5\"},\"speed\":10000}', 1, 0, '2025-01-10 11:35:36.000', '2025-01-10 11:35:36.000', NULL);
COMMIT;

-- ----------------------------
-- Table structure for login_logs
-- ----------------------------
DROP TABLE IF EXISTS `login_logs`;
CREATE TABLE `login_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `ip` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录IP',
  `login_time` datetime(3) NOT NULL COMMENT '登录时间',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of login_logs
-- ----------------------------
BEGIN;
INSERT INTO `login_logs` (`id`, `username`, `ip`, `login_time`, `created_at`, `updated_at`, `user_id`) VALUES (1, 'admin', '::1', '2025-01-11 19:21:12.471', '2025-01-11 19:21:12.471', '2025-01-11 19:21:12.471', 1);
INSERT INTO `login_logs` (`id`, `username`, `ip`, `login_time`, `created_at`, `updated_at`, `user_id`) VALUES (2, 'admin', '::1', '2025-01-11 19:41:16.074', '2025-01-11 19:41:16.075', '2025-01-11 19:41:16.075', 1);
INSERT INTO `login_logs` (`id`, `username`, `ip`, `login_time`, `created_at`, `updated_at`, `user_id`) VALUES (3, 'admin', '::1', '2025-01-11 23:22:36.613', '2025-01-11 23:22:36.614', '2025-01-11 23:22:36.614', 1);
INSERT INTO `login_logs` (`id`, `username`, `ip`, `login_time`, `created_at`, `updated_at`, `user_id`) VALUES (4, 'admin', '::1', '2025-01-11 23:24:57.203', '2025-01-11 23:24:57.203', '2025-01-11 23:24:57.203', 1);
INSERT INTO `login_logs` (`id`, `username`, `ip`, `login_time`, `created_at`, `updated_at`, `user_id`) VALUES (5, 'admin', '::1', '2025-01-12 15:53:07.836', '2025-01-12 15:53:07.838', '2025-01-12 15:53:07.838', 1);
INSERT INTO `login_logs` (`id`, `username`, `ip`, `login_time`, `created_at`, `updated_at`, `user_id`) VALUES (6, 'admin', '192.168.215.1', '2025-01-13 07:09:24.335', '2025-01-13 07:09:24.336', '2025-01-13 07:09:24.336', 1);
INSERT INTO `login_logs` (`id`, `username`, `ip`, `login_time`, `created_at`, `updated_at`, `user_id`) VALUES (7, 'admin', 'fd07:b51a:cc66:1::1', '2025-01-13 16:16:29.165', '2025-01-13 16:16:29.165', '2025-01-13 16:16:29.165', 1);
INSERT INTO `login_logs` (`id`, `username`, `ip`, `login_time`, `created_at`, `updated_at`, `user_id`) VALUES (8, 'admin', '::1', '2025-01-14 08:58:49.735', '2025-01-14 08:58:49.739', '2025-01-14 08:58:49.739', 1);
INSERT INTO `login_logs` (`id`, `username`, `ip`, `login_time`, `created_at`, `updated_at`, `user_id`) VALUES (9, 'admin', '192.168.215.1', '2025-01-16 15:03:11.884', '2025-01-16 15:03:11.884', '2025-01-16 15:03:11.884', 1);
INSERT INTO `login_logs` (`id`, `username`, `ip`, `login_time`, `created_at`, `updated_at`, `user_id`) VALUES (10, 'admin1', '192.168.215.1', '2025-01-16 15:05:19.161', '2025-01-16 15:05:19.162', '2025-01-16 15:05:19.162', 5);
INSERT INTO `login_logs` (`id`, `username`, `ip`, `login_time`, `created_at`, `updated_at`, `user_id`) VALUES (11, 'admin', '::1', '2025-01-16 15:26:46.345', '2025-01-16 15:26:46.348', '2025-01-16 15:26:46.348', 1);
INSERT INTO `login_logs` (`id`, `username`, `ip`, `login_time`, `created_at`, `updated_at`, `user_id`) VALUES (12, 'admin', '::1', '2025-01-16 18:36:07.777', '2025-01-16 18:36:07.782', '2025-01-16 18:36:07.782', 1);
INSERT INTO `login_logs` (`id`, `username`, `ip`, `login_time`, `created_at`, `updated_at`, `user_id`) VALUES (13, 'admin', 'fd07:b51a:cc66:1::1', '2025-01-16 18:55:12.679', '2025-01-16 18:55:12.679', '2025-01-16 18:55:12.679', 1);
COMMIT;

-- ----------------------------
-- Table structure for menus
-- ----------------------------
DROP TABLE IF EXISTS `menus`;
CREATE TABLE `menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` bigint unsigned DEFAULT '0',
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `path` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `component` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `icon` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sort` bigint DEFAULT '0',
  `status` bigint DEFAULT '1',
  `hidden` tinyint(1) DEFAULT '0',
  `cache` tinyint(1) DEFAULT '1',
  `type` bigint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_menus_parent_id` (`parent_id`),
  KEY `idx_menus_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of menus
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for permissions
-- ----------------------------
DROP TABLE IF EXISTS `permissions`;
CREATE TABLE `permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `parent_id` bigint unsigned DEFAULT NULL,
  `path` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `component` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `icon` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sort` bigint DEFAULT '0',
  `status` bigint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_permissions_code` (`code`),
  KEY `idx_permissions_parent_id` (`parent_id`),
  KEY `idx_permissions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of permissions
-- ----------------------------
BEGIN;
INSERT INTO `permissions` (`id`, `name`, `code`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, '仪表盘', 'dashboard', 'menu', NULL, '/dashboard', 'views/dashboard/index', 'dashboard', 1, 1, '2025-01-05 23:44:30.267', '2025-01-06 17:42:08.575', NULL);
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
  KEY `idx_role_permissions_permission_id` (`permission_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of role_permissions
-- ----------------------------
BEGIN;
INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 2, 1, '2025-01-05 23:44:45.922', '2025-01-05 23:44:45.922', NULL);
COMMIT;

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` bigint DEFAULT '1',
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `is_admin` tinyint(1) DEFAULT '0',
  `remark` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_roles_code` (`code`),
  UNIQUE KEY `uni_roles_name` (`name`),
  KEY `idx_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of roles
-- ----------------------------
BEGIN;
INSERT INTO `roles` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`, `is_admin`, `remark`) VALUES (1, '超级管理员', 'admin', '系统超级管理员', 1, '2025-01-05 23:43:34.342', '2025-01-06 11:00:04.362', NULL, 0, NULL);
INSERT INTO `roles` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`, `is_admin`, `remark`) VALUES (2, '普通用户', 'user', '普通用户角色', 1, '2025-01-05 23:44:22.952', '2025-01-07 15:21:51.069', NULL, 0, NULL);
INSERT INTO `roles` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`, `is_admin`, `remark`) VALUES (3, '测试', 'test', '测试', 1, '2025-01-06 14:19:45.851', '2025-01-06 14:19:45.851', '2025-01-06 14:19:49.046', 0, NULL);
COMMIT;

-- ----------------------------
-- Table structure for terminals
-- ----------------------------
DROP TABLE IF EXISTS `terminals`;
CREATE TABLE `terminals` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  `name` varchar(100) COLLATE utf8mb4_general_ci NOT NULL COMMENT '终端名称',
  `host` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '主机地址',
  `port` int NOT NULL DEFAULT '22' COMMENT 'SSH端口',
  `username` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名',
  `password` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '密码',
  `status` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'offline' COMMENT '状态(online/offline)',
  `last_seen` datetime(3) DEFAULT NULL COMMENT '最后在线时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_terminals_name` (`name`),
  UNIQUE KEY `uk_terminals_host_username` (`host`,`username`),
  KEY `idx_terminals_deleted_at` (`deleted_at`),
  KEY `idx_terminals_name` (`name`),
  KEY `idx_terminals_host` (`host`),
  KEY `idx_terminals_status` (`status`),
  CONSTRAINT `terminals_chk_1` CHECK (((`port` >= 1) and (`port` <= 65535))),
  CONSTRAINT `terminals_chk_2` CHECK ((`status` in (_utf8mb4'online',_utf8mb4'offline')))
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='终端管理';

-- ----------------------------
-- Records of terminals
-- ----------------------------
BEGIN;
INSERT INTO `terminals` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `host`, `port`, `username`, `password`, `status`, `last_seen`) VALUES (1, '2025-01-16 16:31:08.588', '2025-01-17 09:50:06.090', NULL, '192.168.10.22', '192.168.10.22', 22, 'parallels', 'lsy@6894', 'online', '2025-01-17 09:50:06.090');
COMMIT;

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
  CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`),
  CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (13, 2, 2, '2025-01-11 18:15:54.800', '2025-01-11 18:15:54.800', NULL);
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (14, 1, 1, '2025-01-11 23:30:38.668', '2025-01-11 23:30:38.668', NULL);
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (15, 4, 2, '2025-01-11 23:34:11.146', '2025-01-11 23:34:11.146', NULL);
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (16, 3, 2, '2025-01-11 23:34:20.982', '2025-01-11 23:34:20.982', '2025-01-11 23:56:23.264');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (17, 3, 1, '2025-01-11 23:34:20.987', '2025-01-11 23:34:20.987', '2025-01-11 23:56:23.264');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (18, 3, 2, '2025-01-11 23:56:23.275', '2025-01-11 23:56:23.275', '2025-01-11 23:56:29.876');
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (19, 3, 2, '2025-01-11 23:56:29.877', '2025-01-11 23:56:29.877', NULL);
INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `created_at`, `updated_at`, `deleted_at`) VALUES (20, 3, 1, '2025-01-11 23:56:29.878', '2025-01-11 23:56:29.878', NULL);
COMMIT;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `password` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码',
  `nickname` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '昵称',
  `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '邮箱',
  `avatar` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '头像',
  `status` bigint DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of users
-- ----------------------------
BEGIN;
INSERT INTO `users` (`id`, `username`, `password`, `nickname`, `email`, `avatar`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 'admin', '$2a$10$tt7pZdek/cV68idHLK3XBORJYVFw.P.ezqvmwRW1KvXxQyUhhoffi', '超级管理员', 'admin@example.com', '', 1, '2025-01-05 23:35:27.267', '2025-01-11 18:20:35.964', NULL);
INSERT INTO `users` (`id`, `username`, `password`, `nickname`, `email`, `avatar`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 'admin2', '$2a$10$PJre1.eIQE5XtP0sdon6o.45pSIFlq0jHCMayEdMo52PBpaHdYSOW', '管理员', 'admin2@example.com', '', 1, '2025-01-05 23:44:08.607', '2025-01-11 18:19:36.264', NULL);
INSERT INTO `users` (`id`, `username`, `password`, `nickname`, `email`, `avatar`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 'admin3', '$2a$10$WSoZIJKE8nVn1ECuP6Ng.eENKooqo.z8VRC/e3.iAjeHAA4rtjDYC', 'admin2', 'admin2@qq.com', '', 0, '2025-01-06 00:37:33.577', '2025-01-12 00:00:56.335', NULL);
INSERT INTO `users` (`id`, `username`, `password`, `nickname`, `email`, `avatar`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 'testuser', '$2a$10$CyrfKTj1HHgdabla9Fl/pO4/OHmJxM0wmwdgs77JK.erEBx3wzwr.', '测试用户', 'test@example.com', '', 0, '2025-01-06 11:54:53.364', '2025-01-12 11:56:41.561', '2025-01-06 14:19:01.876');
INSERT INTO `users` (`id`, `username`, `password`, `nickname`, `email`, `avatar`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 'admin1', '$2a$10$./IcLyN1H.GLOlA3Ml9a6uenq/xY4TtORjjEtKcF9JPpmz4Lbb2/q', 'admin1', 'admin1@qq.com', '', 1, '2025-01-16 15:05:07.799', '2025-01-16 15:05:07.799', NULL);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
