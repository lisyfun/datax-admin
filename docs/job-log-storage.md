# 任务日志存储方案

## 概述

为了解决任务执行日志直接存储在数据库中导致的性能问题，我们重新设计了日志存储方案，将日志内容存储到文件系统中，数据库只保存日志文件的路径和基本信息。

## 问题背景

原有方案的问题：
1. **性能影响**：频繁的日志写入对数据库造成压力
2. **存储成本**：大量日志数据占用数据库存储空间
3. **查询效率**：大量日志数据影响其他查询性能
4. **维护困难**：日志清理和归档变得复杂

## 新方案设计

### 1. 数据库结构变更

在 `job_histories` 表中：
- 添加 `log_path` 字段存储日志文件路径
- `output` 和 `error` 字段变为虚拟字段（gorm:"-"），不再存储到数据库
- 保留原有字段结构，确保向后兼容

### 2. 日志文件存储

**存储位置**：`logs/jobs/` 目录下

**文件命名规则**：
- 格式：`job_{job_id}_{history_id}_{timestamp}.json`
- 示例：`job_1_123_1641234567.json`

**文件内容格式**：
```json
{
  "output": "任务执行的标准输出内容",
  "error": "任务执行的错误信息"
}
```

### 3. 核心功能

#### 日志写入
```go
// 在任务执行完成后写入日志文件
if err := history.WriteLogToFile(output, error); err != nil {
    logger.Info("写入日志文件失败: %v", err)
}
```

#### 日志读取
```go
// 查询历史记录时自动从文件读取日志内容
if err := history.ReadLogFromFile(); err != nil {
    logger.Info("读取日志文件失败: %v", err)
}
```

#### 自动清理
- 支持按天数自动清理过期日志
- 清理孤立的日志文件（数据库中没有对应记录的文件）
- 可配置是否启用自动清理

### 4. 配置选项

在 `config.yaml` 中添加：
```yaml
# 任务日志配置
job_log:
  max_age: 30        # 日志保留天数，默认30天
  auto_cleanup: true # 是否启用自动清理
```

## 实现细节

### 1. 模型变更

```go
type JobHistory struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    JobID     uint      `json:"job_id"`
    // ... 其他字段
    LogPath   string    `json:"log_path" gorm:"size:500"` // 日志文件路径
    Output    string    `json:"output" gorm:"-"`          // 虚拟字段
    Error     string    `json:"error" gorm:"-"`           // 虚拟字段
}
```

### 2. 日志处理方法

```go
// WriteLogToFile 将日志写入文件
func (jh *JobHistory) WriteLogToFile(output, error string) error

// ReadLogFromFile 从文件读取日志
func (jh *JobHistory) ReadLogFromFile() error
```

### 3. 自动清理器

```go
type JobLogCleaner struct {
    maxAge time.Duration // 日志保留时间
}

// CleanOldLogs 清理过期的日志文件
func (c *JobLogCleaner) CleanOldLogs() error

// CleanOrphanedLogs 清理孤立的日志文件
func (c *JobLogCleaner) CleanOrphanedLogs() error
```

## 迁移步骤

### 1. 数据库迁移

执行 SQL 脚本添加 `log_path` 字段：
```sql
ALTER TABLE `job_histories` ADD COLUMN `log_path` VARCHAR(500) DEFAULT '' COMMENT '日志文件路径';
CREATE INDEX `idx_job_histories_log_path` ON `job_histories` (`log_path`);
```

### 2. 代码部署

1. 部署新版本代码
2. 重启应用服务
3. 验证日志文件正常生成

### 3. 历史数据处理

现有的历史记录中的 `output` 和 `error` 字段数据会保留，但新的执行记录将使用文件存储。

## 优势

1. **性能提升**：减少数据库写入压力，提高查询效率
2. **存储优化**：日志文件可以使用文件系统的压缩和归档功能
3. **扩展性好**：支持大量日志数据的存储和管理
4. **维护简单**：自动清理机制，减少人工维护工作
5. **向后兼容**：保持API接口不变，前端无需修改

## 注意事项

1. **文件权限**：确保应用有权限在 `logs/jobs/` 目录下创建和删除文件
2. **磁盘空间**：监控日志目录的磁盘使用情况
3. **备份策略**：如需备份日志，需要同时备份数据库和日志文件
4. **并发安全**：文件操作已考虑并发安全性
5. **错误处理**：日志文件读写失败不会影响任务执行的正常流程

## 监控建议

1. 监控 `logs/jobs/` 目录的磁盘使用情况
2. 定期检查日志清理器的运行状态
3. 监控日志文件读写的错误率
4. 定期验证日志文件的完整性
