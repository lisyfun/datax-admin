package controllers

import (
	"datax-admin/services"
	"datax-admin/types"
	"datax-admin/utils/logger"
	"io"
	"net/http"
	"strconv"
	"time"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// TerminalController 终端控制器
type TerminalController struct {
	terminalService *services.TerminalService
}

// NewTerminalController 创建终端控制器
func NewTerminalController() *TerminalController {
	return &TerminalController{
		terminalService: &services.TerminalService{},
	}
}

// CreateTerminal 创建终端
func (c *TerminalController) CreateTerminal(ctx *gin.Context) {
	var req types.CreateTerminalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.terminalService.CreateTerminal(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

// UpdateTerminal 更新终端
func (c *TerminalController) UpdateTerminal(ctx *gin.Context) {
	var req types.UpdateTerminalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的终端ID"})
		return
	}

	if err := c.terminalService.UpdateTerminal(uint(id), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteTerminal 删除终端
func (c *TerminalController) DeleteTerminal(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的终端ID"})
		return
	}

	if err := c.terminalService.DeleteTerminal(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetTerminalList 获取终端列表
func (c *TerminalController) GetTerminalList(ctx *gin.Context) {
	var req types.TerminalListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.terminalService.GetTerminalList(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateTerminalStatus 更新终端状态
func (c *TerminalController) UpdateTerminalStatus(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的终端ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=online offline"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.terminalService.UpdateTerminalStatus(uint(id), req.Status); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}

// GetTerminalByID 获取终端详情
func (c *TerminalController) GetTerminalByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的终端ID"})
		return
	}

	terminal, err := c.terminalService.GetTerminalByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, terminal)
}

// ConnectTerminal WebSocket终端连接处理
func (c *TerminalController) ConnectTerminal(ctx *gin.Context) {
	logger.Info("收到终端连接请求")

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		logger.Error("无效的终端ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的终端ID"})
		return
	}

	logger.Info("正在获取终端信息, ID: %d", id)
	terminal, err := c.terminalService.GetTerminalByID(uint(id))
	if err != nil {
		logger.Error("获取终端信息失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if terminal.Password == "" {
		logger.Error("终端密码未设置, ID: %d", id)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "终端密码未设置"})
		return
	}

	logger.Info("正在升级WebSocket连接, ID: %d", id)
	// 升级HTTP连接为WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 在生产环境中，应该根据实际情况配置允许的域名
		},
		HandshakeTimeout: 10 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
	}

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Error("WebSocket升级失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket升级失败"})
		return
	}
	defer ws.Close()

	logger.Info("正在创建SSH连接, Host: %s, Port: %d", terminal.Host, terminal.Port)
	// 创建SSH连接
	sshClient, err := services.NewSSHClient(terminal.Host, terminal.Port, terminal.Username, terminal.Password)
	if err != nil {
		logger.Error("SSH连接失败: %v", err)
		ws.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": err.Error(),
		})
		return
	}
	defer sshClient.Close()

	logger.Info("正在更新终端状态为在线, ID: %d", id)
	// 更新终端状态为在线
	if err := c.terminalService.UpdateTerminalStatus(uint(id), "online"); err != nil {
		logger.Error("更新终端状态失败: %v", err)
		ws.WriteJSON(map[string]interface{}{
			"type": "error",
			"data": "更新终端状态失败",
		})
		return
	}

	logger.Info("终端连接成功, 开始数据传输, ID: %d", id)

	// 创建用于同步的通道
	done := make(chan bool)

	// 从SSH读取数据并发送到WebSocket
	go func() {
		defer close(done)
		buf := make([]byte, 1024)
		for {
			n, err := sshClient.Read(buf)
			if err != nil {
				if err != io.EOF {
					logger.Error("读取SSH数据失败: %v", err)
					ws.WriteJSON(map[string]interface{}{
						"type": "error",
						"data": err.Error(),
					})
				}
				return
			}
			if n > 0 {
				err = ws.WriteJSON(map[string]interface{}{
					"type": "output",
					"data": string(buf[:n]),
				})
				if err != nil {
					logger.Error("发送WebSocket数据失败: %v", err)
					return
				}
			}
		}
	}()

	// 从WebSocket读取数据并写入SSH
	go func() {
		for {
			var msg struct {
				Type string `json:"type"`
				Data string `json:"data"`
			}
			err := ws.ReadJSON(&msg)
			if err != nil {
				logger.Error("读取WebSocket数据失败: %v", err)
				return
			}

			switch msg.Type {
			case "input":
				_, err = sshClient.Write([]byte(msg.Data))
				if err != nil {
					logger.Error("写入SSH数据失败: %v", err)
					ws.WriteJSON(map[string]interface{}{
						"type": "error",
						"data": err.Error(),
					})
					return
				}
			case "resize":
				var resize struct {
					Cols int `json:"cols"`
					Rows int `json:"rows"`
				}
				if err := json.Unmarshal([]byte(msg.Data), &resize); err == nil {
					if err := sshClient.ResizeTerminal(resize.Cols, resize.Rows); err != nil {
						logger.Error("调整终端大小失败: %v", err)
					}
				}
			}
		}
	}()

	<-done

	// 更新终端状态为离线
	if err := c.terminalService.UpdateTerminalStatus(uint(id), "offline"); err != nil {
		logger.Error("更新终端状态失败: %v", err)
	}
}
