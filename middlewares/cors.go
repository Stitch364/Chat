package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
func Cors(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin") // 获取请求的 Origin 头部的值，Origin 头部在跨域请求中很重要，会告诉服务器请求的来源域

	// 接收客户端发送的 Origin（重要）
	c.Header("Access-Control-Allow-Origin", origin)
	// 服务器支持的所有跨域请求的方法
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, HEAD, PUT")
	// 允许跨域设置可以返回其他子段，可以自定义字段
	c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,content-type,Authorization,Content-Length,X-CSRF-AccessToken,AccessToken,session, token")
	// 允许浏览器（客户端）可以解析的头部（重要）
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, token")
	// 允许客户端传递校验信息比如 cookie（重要）
	c.Header("Access-Control-Allow-Credentials", "true")

	// 显式允许 WebSocket 升级头
	if c.GetHeader("Upgrade") == "websocket" {
		c.Writer.Header().Set("Connection", "Upgrade")
		c.Writer.Header().Set("Upgrade", "websocket")
	}

	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

func Cors2(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin")

	// 允许的跨域配置
	if origin != "" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin) // 或 "*"
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "AccountToken, Content-Type, Authorization, AccessToken, X-CSRF-Token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 显式允许 WebSocket 升级头
		if c.GetHeader("Upgrade") == "websocket" {
			c.Writer.Header().Set("Connection", "Upgrade")
			c.Writer.Header().Set("Upgrade", "websocket")
		}
	}

	// 处理预检请求
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()

	c.Next()
}
