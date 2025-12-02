package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"datax-admin/models"
	"datax-admin/services"
	"datax-admin/types"

	"github.com/gin-gonic/gin"
)

// RedisController 控制器
type RedisController struct {
	svc *services.RedisService
}

func NewRedisController() *RedisController {
	return &RedisController{svc: &services.RedisService{}}
}

// 连接管理
func (c *RedisController) ListConnections(ctx *gin.Context) {
	var req types.RedisConnectionListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var total int64
	var conns []models.RedisConnection
	q := models.DB.Model(&models.RedisConnection{})
	if req.Keyword != "" {
		q = q.Where("name LIKE ?", "%"+req.Keyword+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	offset := (req.Page - 1) * req.PageSize
	if err := q.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&conns).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	items := make([]types.RedisConnectionInfo, len(conns))
	for i, v := range conns {
		items[i] = types.RedisConnectionInfo{ID: v.ID, Name: v.Name, Host: v.Host, Port: v.Port, Username: v.Username, DB: v.DB, UseTLS: v.UseTLS}
	}
	ctx.JSON(http.StatusOK, types.RedisConnectionListResponse{Total: total, Items: items})
}

func (c *RedisController) CreateConnection(ctx *gin.Context) {
	var req types.CreateRedisConnectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m := &models.RedisConnection{Name: req.Name, Host: req.Host, Port: req.Port, Username: req.Username, Password: req.Password, DB: req.DB, UseTLS: req.UseTLS}
	if err := models.DB.Create(m).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": m.ID})
}

func (c *RedisController) UpdateConnection(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	var req types.UpdateRedisConnectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]any{"name": req.Name, "host": req.Host, "port": req.Port, "username": req.Username, "password": req.Password, "db": req.DB, "use_tls": req.UseTLS}
	if err := models.DB.Model(&models.RedisConnection{}).Where("id = ?", uint(id)).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (c *RedisController) DeleteConnection(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	if err := models.DB.Delete(&models.RedisConnection{}, uint(id)).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (c *RedisController) TestConnection(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	if _, err := c.svc.KeyType(ctx, uint(id), 0, "__ping__"); err != nil {
		// 使用 svc.getClientByConnectionID + Ping 更直接，此处复用 svc
		if _, err2 := c.svc.GetTTL(ctx, uint(id), 0, "__ping__"); err2 != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err2.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

func (c *RedisController) SelectDB(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	var body struct {
		DB *int `json:"db" form:"db" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		if err2 := ctx.ShouldBindQuery(&body); err2 != nil {
			dbStr := ctx.Query("db")
			if dbStr == "" {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if v, e := strconv.ParseInt(dbStr, 10, 64); e == nil {
				vv := int(v)
				body.DB = &vv
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 db"})
				return
			}
		}
	}
	if body.DB == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "缺少 db"})
		return
	}
	if err := c.svc.SelectDB(ctx, uint(id), *body.DB); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "切换成功"})
}

// 键操作
func (c *RedisController) ListKeys(ctx *gin.Context) {
	var req types.RedisKeysRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	next, keys, err := c.svc.ListKeys(ctx, req.ConnID, req.DB, req.Pattern, req.Cursor, req.Count)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 可选类型过滤
	if req.Type != "" {
		filtered := make([]string, 0, len(keys))
		for _, k := range keys {
			typ, e := c.svc.KeyType(ctx, req.ConnID, req.DB, k)
			if e != nil {
				continue
			}
			if typ == req.Type {
				filtered = append(filtered, k)
			}
		}
		keys = filtered
	}
	ctx.JSON(http.StatusOK, gin.H{"cursor": next, "keys": keys})
}

func (c *RedisController) GetKey(ctx *gin.Context) {
	key := ctx.Param("key")
	if key == "" {
		key = ctx.Query("key")
	}
	connIDStr := ctx.Query("conn_id")
	if connIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "缺少 conn_id"})
		return
	}
	cid, err := strconv.ParseUint(connIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 conn_id"})
		return
	}
	dbStr := ctx.Query("db")
	db := 0
	if dbStr != "" {
		if v, e := strconv.ParseInt(dbStr, 10, 64); e == nil {
			db = int(v)
		}
	}
	typ, val, err := c.svc.GetValue(ctx, uint(cid), db, key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, types.RedisKeyValueResponse{Type: typ, Key: key, Value: val})
}

func (c *RedisController) GetKeyPost(ctx *gin.Context) {
	var body struct {
		ConnID uint   `json:"conn_id" binding:"required"`
		DB     int    `json:"db"`
		Key    string `json:"key" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	typ, val, err := c.svc.GetValue(ctx, body.ConnID, body.DB, body.Key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, types.RedisKeyValueResponse{Type: typ, Key: body.Key, Value: val})
}

func (c *RedisController) SetKey(ctx *gin.Context) {
	var body struct {
		ConnID uint        `json:"conn_id" binding:"required"`
		DB     int         `json:"db"`
		Type   string      `json:"type" binding:"required"`
		Key    string      `json:"key" binding:"required"`
		Value  interface{} `json:"value" binding:"required"`
		TTL    *int64      `json:"ttl_seconds"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.SetValue(ctx, body.ConnID, body.DB, body.Key, body.Type, body.Value, body.TTL); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "设置成功"})
}

func (c *RedisController) DeleteKey(ctx *gin.Context) {
	key := ctx.Param("key")
	var body struct {
		ConnID uint `json:"conn_id" binding:"required"`
		DB     int  `json:"db"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.DeleteKey(ctx, body.ConnID, body.DB, key); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (c *RedisController) DeleteKeyPost(ctx *gin.Context) {
	var body struct {
		ConnID uint   `json:"conn_id" binding:"required"`
		DB     int    `json:"db"`
		Key    string `json:"key" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.DeleteKey(ctx, body.ConnID, body.DB, body.Key); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (c *RedisController) GetTTL(ctx *gin.Context) {
	key := ctx.Param("key")
	connIDStr := ctx.Query("conn_id")
	if connIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "缺少 conn_id"})
		return
	}
	cid, err := strconv.ParseUint(connIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 conn_id"})
		return
	}
	dbStr := ctx.Query("db")
	db := 0
	if dbStr != "" {
		if v, e := strconv.ParseInt(dbStr, 10, 64); e == nil {
			db = int(v)
		}
	}
	ttl, err := c.svc.GetTTL(ctx, uint(cid), db, key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ttl": ttl})
}

func (c *RedisController) SetTTL(ctx *gin.Context) {
	key := ctx.Param("key")
	var req types.RedisExpireRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, err := c.svc.Expire(ctx, req.ConnID, req.DB, key, req.Seconds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": ok})
}

func (c *RedisController) GetTTLPost(ctx *gin.Context) {
	var body struct {
		ConnID uint   `json:"conn_id" binding:"required"`
		DB     int    `json:"db"`
		Key    string `json:"key" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ttl, err := c.svc.GetTTL(ctx, body.ConnID, body.DB, body.Key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ttl": ttl})
}

func (c *RedisController) SetTTLPost(ctx *gin.Context) {
	var body struct {
		ConnID  uint   `json:"conn_id" binding:"required"`
		DB      int    `json:"db"`
		Key     string `json:"key" binding:"required"`
		Seconds int64  `json:"seconds" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, err := c.svc.Expire(ctx, body.ConnID, body.DB, body.Key, body.Seconds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": ok})
}

func (c *RedisController) RenameKey(ctx *gin.Context) {
	var body struct {
		ConnID uint   `json:"conn_id" binding:"required"`
		DB     int    `json:"db"`
		Key    string `json:"key" binding:"required"`
		NewKey string `json:"new_key" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.RenameKey(ctx, body.ConnID, body.DB, body.Key, body.NewKey); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "重命名成功"})
}

func (c *RedisController) CopyKey(ctx *gin.Context) {
	var body struct {
		ConnID  uint   `json:"conn_id" binding:"required"`
		DB      int    `json:"db"`
		Key     string `json:"key" binding:"required"`
		DestKey string `json:"dest_key" binding:"required"`
		Replace bool   `json:"replace"`
		DestDB  *int   `json:"dest_db"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, err := c.svc.CopyKey(ctx, body.ConnID, body.DB, body.Key, body.DestKey, body.Replace, body.DestDB)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": ok})
}

func (c *RedisController) MoveKey(ctx *gin.Context) {
	var body struct {
		ConnID uint   `json:"conn_id" binding:"required"`
		DB     int    `json:"db"`
		Key    string `json:"key" binding:"required"`
		DestDB int    `json:"dest_db" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, err := c.svc.MoveKey(ctx, body.ConnID, body.DB, body.Key, body.DestDB)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": ok})
}

func (c *RedisController) ExportKeys(ctx *gin.Context) {
	connIDStr := ctx.Query("conn_id")
	if connIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "缺少 conn_id"})
		return
	}
	cid, err := strconv.ParseUint(connIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 conn_id"})
		return
	}
	pattern := ctx.Query("pattern")
	countStr := ctx.Query("count")
	var count int64 = 100
	if countStr != "" {
		if v, e := strconv.ParseInt(countStr, 10, 64); e == nil {
			count = v
		}
	}
	dbStr := ctx.Query("db")
	db := 0
	if dbStr != "" {
		if v, e := strconv.ParseInt(dbStr, 10, 64); e == nil {
			db = int(v)
		}
	}
	items, err := c.svc.ExportKeys(ctx, uint(cid), db, pattern, count)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"items": items})
}

func (c *RedisController) CountKeys(ctx *gin.Context) {
	connIDStr := ctx.Query("conn_id")
	if connIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "缺少 conn_id"})
		return
	}
	cid, err := strconv.ParseUint(connIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 conn_id"})
		return
	}
	pattern := ctx.Query("pattern")
	batchStr := ctx.Query("batch")
	var batch int64 = 100
	if batchStr != "" {
		if v, e := strconv.ParseInt(batchStr, 10, 64); e == nil {
			batch = v
		}
	}
	dbStr := ctx.Query("db")
	db := 0
	if dbStr != "" {
		if v, e2 := strconv.ParseInt(dbStr, 10, 64); e2 == nil {
			db = int(v)
		}
	}
	total, e := c.svc.CountKeys(ctx, uint(cid), db, pattern, batch)
	if e != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"total": total})
}

func (c *RedisController) ExecuteCLI(ctx *gin.Context) {
	var body struct {
		ConnID uint   `json:"conn_id" binding:"required"`
		DB     int    `json:"db"`
		Line   string `json:"line" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	argv := splitArgs(body.Line)
	if len(argv) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "命令为空"})
		return
	}
	res, err := c.svc.Execute(ctx, body.ConnID, body.DB, argv)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"argv": argv, "result": res})
}

func splitArgs(s string) []string {
	var args []string
	var b strings.Builder
	inQuote := false
	var qc rune
	for _, r := range s {
		if inQuote {
			if r == qc {
				inQuote = false
			} else {
				b.WriteRune(r)
			}
			continue
		}
		if r == '\'' || r == '"' {
			inQuote = true
			qc = r
			continue
		}
		if r == ' ' || r == '\t' {
			if b.Len() > 0 {
				args = append(args, b.String())
				b.Reset()
			}
			continue
		}
		b.WriteRune(r)
	}
	if b.Len() > 0 {
		args = append(args, b.String())
	}
	return args
}

func (c *RedisController) ExecuteCLIBulk(ctx *gin.Context) {
	var body struct {
		ConnID uint     `json:"conn_id" binding:"required"`
		DB     int      `json:"db"`
		Lines  []string `json:"lines" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(body.Lines) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "命令为空"})
		return
	}
	argvv := make([][]string, len(body.Lines))
	for i, line := range body.Lines {
		argvv[i] = splitArgs(line)
	}
	res, err := c.svc.ExecuteBulk(ctx, body.ConnID, body.DB, argvv)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"results": res})
}
