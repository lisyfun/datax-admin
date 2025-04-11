package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityMiddleware 安全中间件
// 用于防止SQL注入和XSLT注入等攻击
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只处理POST、PUT和PATCH请求
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			// 读取请求体
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取请求数据"})
				c.Abort()
				return
			}

			// 重新设置请求体，因为读取后会被消费
			c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

			// 检查JSON格式的请求体是否有SQL注入尝试
			var data map[string]any
			if err := json.Unmarshal(bodyBytes, &data); err == nil {
				// 检查所有字符串值是否有注入尝试
				if containsInjection(data) {
					c.JSON(http.StatusBadRequest, gin.H{"error": "检测到可疑输入"})
					c.Abort()
					return
				}
			}
		}

		// 继续处理请求
		c.Next()
	}
}

// 递归检查所有值是否有注入尝试
func containsInjection(data map[string]any) bool {
	for _, v := range data {
		switch value := v.(type) {
		case string:
			if containsSQLInjection(value) || containsXSLTInjection(value) {
				return true
			}
		case map[string]any:
			if containsInjection(value) {
				return true
			}
		case []any:
			for _, item := range value {
				if mapItem, ok := item.(map[string]any); ok {
					if containsInjection(mapItem) {
						return true
					}
				} else if strItem, ok := item.(string); ok {
					if containsSQLInjection(strItem) || containsXSLTInjection(strItem) {
						return true
					}
				}
			}
		}
	}
	return false
}

// containsSQLInjection 检查是否包含SQL注入
func containsSQLInjection(input string) bool {
	// 检查常见的SQL注入模式
	sqlPatterns := []string{
		"--", ";--", ";", "/*", "*/", "@@", "@",
		"char(", "nchar(", "varchar(", "nvarchar(",
		"alter ", "begin ", "cast(", "create ", "cursor ",
		"declare ", "delete ", "drop ", "end ", "exec ",
		"execute ", "fetch ", "insert ", "kill ", "select ",
		"sys.", "sysobjects", "syscolumns",
		"table ", "update ", "xp_",
		"or 1=1", "or 1 = 1", "or '1'='1'",
		"union ", "UNION ", "HAVING ",
		"'1'='1", "' OR '",
	}

	inputLower := strings.ToLower(input)
	for _, pattern := range sqlPatterns {
		if strings.Contains(inputLower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

// containsXSLTInjection 检查是否包含XSLT注入
func containsXSLTInjection(input string) bool {
	// 检查常见的XSLT注入模式
	xsltPatterns := []string{
		"<xsl:", "</xsl:", "xsl:value-of", "xsl:for-each",
		"xsl:template", "xsl:variable", "xsl:param",
		"system-property", "document(", "unparsed-entity-uri(",
		"<![CDATA[", "<!ENTITY", "<!DOCTYPE",
		"<?xml", "<xml", "xmlns:",
	}

	// 如果包含尖括号和XSLT相关字符串，可能是XSLT注入
	for _, pattern := range xsltPatterns {
		if strings.Contains(strings.ToLower(input), strings.ToLower(pattern)) {
			return true
		}
	}

	// 检查基本的XML标签模式
	if (strings.Contains(input, "<") && strings.Contains(input, ">")) ||
		(strings.Contains(input, "&lt;") && strings.Contains(input, "&gt;")) {
		return true
	}

	return false
}
