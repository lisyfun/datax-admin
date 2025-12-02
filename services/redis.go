package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"datax-admin/models"

	"github.com/redis/go-redis/v9"
)

// RedisService 封装 Redis 连接与基础操作
type RedisService struct{}

var (
	redisClientsMu sync.RWMutex
	redisClients   = make(map[string]*redis.Client)
)

func clientKey(connID uint, db int) string { return fmt.Sprintf("%d:%d", connID, db) }

// getClientByConnectionIDAndDB 根据连接ID与DB获取/创建 Redis 客户端
func (s *RedisService) getClientByConnectionIDAndDB(ctx context.Context, connID uint, db int) (*redis.Client, error) {
	redisClientsMu.RLock()
	if c, ok := redisClients[clientKey(connID, db)]; ok {
		redisClientsMu.RUnlock()
		return c, nil
	}
	redisClientsMu.RUnlock()

	// 加载连接配置
	var conn models.RedisConnection
	if err := models.DB.WithContext(ctx).First(&conn, connID).Error; err != nil {
		return nil, err
	}

	// 构建 options
	addr := fmt.Sprintf("%s:%d", conn.Host, conn.Port)
	opts := &redis.Options{
		Addr:      addr,
		Username:  conn.Username,
		Password:  conn.Password,
		DB:        db,
		TLSConfig: nil,
	}

	// TODO: 根据 UseTLS 扩展 TLSConfig（证书路径/跳过验证等）

	client := redis.NewClient(opts)
	// 连接测试
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	// 缓存客户端
	redisClientsMu.Lock()
	redisClients[clientKey(connID, db)] = client
	redisClientsMu.Unlock()

	return client, nil
}

// fingerprint 生成简单指纹（用于日志或调试，不含敏感信息）
func fingerprint(host string, port int, db int) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s:%d/%d", host, port, db)))
	return hex.EncodeToString(h[:8])
}

// ListKeys 使用 SCAN 游标分页
func (s *RedisService) ListKeys(ctx context.Context, connID uint, db int, pattern string, cursor uint64, count int64) (nextCursor uint64, keys []string, err error) {
	client, err := s.getClientByConnectionIDAndDB(ctx, connID, db)
	if err != nil {
		return 0, nil, err
	}
	if count <= 0 {
		count = 100
	}
	iter := client.Scan(ctx, cursor, pattern, count)
	keys, nextCursor, err = iter.Result()
	return
}

// KeyType 获取 key 类型
func (s *RedisService) KeyType(ctx context.Context, connID uint, db int, key string) (string, error) {
	client, err := s.getClientByConnectionIDAndDB(ctx, connID, db)
	if err != nil {
		return "", err
	}
	return client.Type(ctx, key).Result()
}

// GetValue 根据类型读取值
func (s *RedisService) GetValue(ctx context.Context, connID uint, db int, key string) (typ string, value any, err error) {
	client, err := s.getClientByConnectionIDAndDB(ctx, connID, db)
	if err != nil {
		return "", nil, err
	}
	typ, err = client.Type(ctx, key).Result()
	if err != nil {
		return "", nil, err
	}
	switch typ {
	case "string":
		val, e := client.Get(ctx, key).Result()
		return typ, val, e
	case "hash":
		m, e := client.HGetAll(ctx, key).Result()
		return typ, m, e
	case "list":
		arr, e := client.LRange(ctx, key, 0, -1).Result()
		return typ, arr, e
	case "set":
		arr, e := client.SMembers(ctx, key).Result()
		return typ, arr, e
	case "zset":
		zs, e := client.ZRangeWithScores(ctx, key, 0, -1).Result()
		return typ, zs, e
	default:
		return typ, nil, nil
	}
}

// SetValue 根据类型写入值（string/hash/list/set/zset），可选 ttl 秒
func (s *RedisService) SetValue(ctx context.Context, connID uint, db int, key string, typ string, value any, ttlSeconds *int64) error {
	client, err := s.getClientByConnectionIDAndDB(ctx, connID, db)
	if err != nil {
		return err
	}
	switch typ {
	case "string":
		if err := client.Set(ctx, key, value, 0).Err(); err != nil {
			return err
		}
	case "hash":
		switch v := value.(type) {
		case map[string]string:
			if err := client.HSet(ctx, key, v).Err(); err != nil {
				return err
			}
		case map[string]any:
			if err := client.HSet(ctx, key, v).Err(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported hash value type")
		}
	case "list":
		switch v := value.(type) {
		case []string:
			if err := client.Del(ctx, key).Err(); err != nil {
				return err
			}
			if len(v) > 0 {
				if err := client.RPush(ctx, key, toAnySlice(v)...).Err(); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("unsupported list value type")
		}
	case "set":
		switch v := value.(type) {
		case []string:
			if err := client.Del(ctx, key).Err(); err != nil {
				return err
			}
			if len(v) > 0 {
				if err := client.SAdd(ctx, key, toAnySlice(v)...).Err(); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("unsupported set value type")
		}
	case "zset":
		// value 期望为 []redis.Z
		switch v := value.(type) {
		case []redis.Z:
			if err := client.Del(ctx, key).Err(); err != nil {
				return err
			}
			if len(v) > 0 {
				if err := client.ZAdd(ctx, key, v...).Err(); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("unsupported zset value type")
		}
	default:
		return fmt.Errorf("unsupported type: %s", typ)
	}

	if ttlSeconds != nil {
		if *ttlSeconds <= 0 {
			// persist
			if err := client.Persist(ctx, key).Err(); err != nil {
				return err
			}
		} else {
			if err := client.Expire(ctx, key, toDurationSeconds(*ttlSeconds)).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteKey 删除键
func (s *RedisService) DeleteKey(ctx context.Context, connID uint, db int, key string) error {
	client, err := s.getClientByConnectionIDAndDB(ctx, connID, db)
	if err != nil {
		return err
	}
	return client.Del(ctx, key).Err()
}

// GetTTL 获取 TTL（秒），-1 表示无过期，-2 表示不存在
func (s *RedisService) GetTTL(ctx context.Context, connID uint, db int, key string) (int64, error) {
	client, err := s.getClientByConnectionIDAndDB(ctx, connID, db)
	if err != nil {
		return 0, err
	}
	d, err := client.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if d == -2*time.Second {
		return -2, nil
	}
	if d == -1*time.Second {
		return -1, nil
	}
	return int64(d.Seconds()), nil
}

// Expire 设置过期（秒），<=0 视为持久化
func (s *RedisService) Expire(ctx context.Context, connID uint, db int, key string, seconds int64) (bool, error) {
	client, err := s.getClientByConnectionIDAndDB(ctx, connID, db)
	if err != nil {
		return false, err
	}
	if seconds <= 0 {
		if err := client.Persist(ctx, key).Err(); err != nil {
			return false, err
		}
		return true, nil
	}
	return client.Expire(ctx, key, toDurationSeconds(seconds)).Result()
}

// RenameKey 重命名键
func (s *RedisService) RenameKey(ctx context.Context, connID uint, db int, key string, newKey string) error {
	client, err := s.getClientByConnectionIDAndDB(ctx, connID, db)
	if err != nil {
		return err
	}
	return client.Rename(ctx, key, newKey).Err()
}

// CopyKey 使用 Redis COPY 命令
func (s *RedisService) CopyKey(ctx context.Context, connID uint, db int, key string, destKey string, replace bool, destDB *int) (bool, error) {
	client, err := s.getClientByConnectionIDAndDB(ctx, connID, db)
	if err != nil {
		return false, err
	}
	args := []any{"COPY", key, destKey}
	if destDB != nil {
		args = append(args, "DB", *destDB)
	}
	if replace {
		args = append(args, "REPLACE")
	}
	res, err := client.Do(ctx, args...).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

// MoveKey 移动到目标DB
func (s *RedisService) MoveKey(ctx context.Context, connID uint, db int, key string, destDB int) (bool, error) {
	client, err := s.getClientByConnectionIDAndDB(ctx, connID, db)
	if err != nil {
		return false, err
	}
	res, err := client.Do(ctx, "MOVE", key, destDB).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

func (s *RedisService) SelectDB(ctx context.Context, connID uint, db int) error {
	_, err := s.getClientByConnectionIDAndDB(ctx, connID, db)
	return err
}

// ExportKeys 依据 pattern 导出有限数量键及值
type ExportItem struct {
	Key   string      `json:"key"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
	TTL   int64       `json:"ttl"`
}

func (s *RedisService) ExportKeys(ctx context.Context, connID uint, db int, pattern string, count int64) ([]ExportItem, error) {
	if count <= 0 {
		count = 100
	}
	next, keys, err := s.ListKeys(ctx, connID, db, pattern, 0, count)
	_ = next
	if err != nil {
		return nil, err
	}
	items := make([]ExportItem, 0, len(keys))
	for _, k := range keys {
		typ, val, err := s.GetValue(ctx, connID, db, k)
		if err != nil {
			continue
		}
		ttl, _ := s.GetTTL(ctx, connID, db, k)
		items = append(items, ExportItem{Key: k, Type: typ, Value: val, TTL: ttl})
	}
	return items, nil
}

func (s *RedisService) CountKeys(ctx context.Context, connID uint, db int, pattern string, batch int64) (int64, error) {
	client, err := s.getClientByConnectionIDAndDB(ctx, connID, db)
	if err != nil {
		return 0, err
	}
	if batch <= 0 {
		batch = 100
	}
	var cursor uint64 = 0
	var total int64 = 0
	for {
		keys, next, e := client.Scan(ctx, cursor, pattern, batch).Result()
		if e != nil {
			return 0, e
		}
		total += int64(len(keys))
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return total, nil
}

// 辅助：[]string -> []any
func toAnySlice(ss []string) []any {
	res := make([]any, len(ss))
	for i, s := range ss {
		res[i] = s
	}
	return res
}

// 秒转为 Duration
func toDurationSeconds(sec int64) time.Duration {
	return time.Second * time.Duration(sec)
}

// 结构化便捷方法
func (s *RedisService) HashSet(ctx context.Context, connID uint, db int, key string, fields map[string]string, ttlSeconds *int64) error {
	return s.SetValue(ctx, connID, db, key, "hash", fields, ttlSeconds)
}

func (s *RedisService) ListSet(ctx context.Context, connID uint, db int, key string, items []string, ttlSeconds *int64) error {
	return s.SetValue(ctx, connID, db, key, "list", items, ttlSeconds)
}

func (s *RedisService) SetSet(ctx context.Context, connID uint, db int, key string, members []string, ttlSeconds *int64) error {
	return s.SetValue(ctx, connID, db, key, "set", members, ttlSeconds)
}

func (s *RedisService) ZSetSet(ctx context.Context, connID uint, db int, key string, members []redis.Z, ttlSeconds *int64) error {
	var mv []redis.Z = members
	return s.SetValue(ctx, connID, db, key, "zset", mv, ttlSeconds)
}

func (s *RedisService) Execute(ctx context.Context, connID uint, db int, argv []string) (any, error) {
	if len(argv) == 0 {
		return nil, fmt.Errorf("empty command")
	}
	client, err := s.getClientByConnectionIDAndDB(ctx, connID, db)
	if err != nil {
		return nil, err
	}
	args := make([]any, 0, len(argv))
	for _, a := range argv {
		args = append(args, a)
	}
	res, err := client.Do(ctx, args...).Result()
	if err != nil {
		return nil, err
	}
	return normalizeResult(res), nil
}

func normalizeResult(v any) any {
	switch t := v.(type) {
	case []byte:
		return string(t)
	case []interface{}:
		for i := range t {
			t[i] = normalizeResult(t[i])
		}
		return t
	case map[string]interface{}:
		for k := range t {
			t[k] = normalizeResult(t[k])
		}
		return t
	default:
		return v
	}
}

func (s *RedisService) ExecuteBulk(ctx context.Context, connID uint, db int, lines [][]string) ([]any, error) {
	client, err := s.getClientByConnectionIDAndDB(ctx, connID, db)
	if err != nil {
		return nil, err
	}
	results := make([]any, len(lines))
	for i, argv := range lines {
		if len(argv) == 0 {
			results[i] = nil
			continue
		}
		args := make([]any, 0, len(argv))
		for _, a := range argv {
			args = append(args, a)
		}
		v, e := client.Do(ctx, args...).Result()
		if e != nil {
			results[i] = map[string]any{"error": e.Error()}
		} else {
			results[i] = normalizeResult(v)
		}
	}
	return results, nil
}
