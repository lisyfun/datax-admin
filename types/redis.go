package types

// CreateRedisConnectionRequest 创建 Redis 连接请求
type CreateRedisConnectionRequest struct {
    Name     string `json:"name" binding:"required,max=100"`
    Host     string `json:"host" binding:"required,max=200"`
    Port     int    `json:"port" binding:"min=1,max=65535"`
    Username string `json:"username"`
    Password string `json:"password"`
    DB       int    `json:"db"`
    UseTLS   bool   `json:"use_tls"`
}

// UpdateRedisConnectionRequest 更新 Redis 连接请求
type UpdateRedisConnectionRequest struct {
    Name     string `json:"name" binding:"required,max=100"`
    Host     string `json:"host" binding:"required,max=200"`
    Port     int    `json:"port" binding:"min=1,max=65535"`
    Username string `json:"username"`
    Password string `json:"password"`
    DB       int    `json:"db"`
    UseTLS   bool   `json:"use_tls"`
}

// RedisConnectionInfo 连接信息响应（不含敏感字段）
type RedisConnectionInfo struct {
    ID      uint   `json:"id"`
    Name    string `json:"name"`
    Host    string `json:"host"`
    Port    int    `json:"port"`
    Username string `json:"username"`
    DB      int    `json:"db"`
    UseTLS  bool   `json:"use_tls"`
}

// RedisConnectionListRequest 列表请求
type RedisConnectionListRequest struct {
    Page     int    `form:"page,default=1" binding:"min=1"`
    PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
    Keyword  string `form:"keyword"`
}

// RedisConnectionListResponse 列表响应
type RedisConnectionListResponse struct {
    Total int64                 `json:"total"`
    Items []RedisConnectionInfo `json:"items"`
}

// RedisKeysRequest 键列表请求
type RedisKeysRequest struct {
    ConnID  uint   `form:"conn_id" binding:"required"`
    DB      int    `form:"db,default=0"`
    Pattern string `form:"pattern"`
    Cursor  uint64 `form:"cursor,default=0"`
    Count   int64  `form:"count,default=100"`
    Type    string `form:"type"` // 可选过滤类型
}

// RedisKeyValueResponse 键值响应
type RedisKeyValueResponse struct {
    Type  string      `json:"type"`
    Key   string      `json:"key"`
    Value interface{} `json:"value"`
}

// RedisExpireRequest 设置过期请求
type RedisExpireRequest struct {
    ConnID  uint  `json:"conn_id" binding:"required"`
    DB      int   `json:"db,default=0"`
    Seconds int64 `json:"seconds"`
}
