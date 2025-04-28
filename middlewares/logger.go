package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var lg *zap.Logger

const Body = "body"

// ErrLogMsg 日志数据
func ErrLogMsg(ctx *gin.Context) []zap.Field {
	var body string
	data, ok := ctx.Get(Body)
	if ok {
		body = string(data.([]byte))
	}
	path := ctx.Request.URL.Path
	query := ctx.Request.URL.RawQuery
	fields := []zap.Field{
		zap.Int("status", ctx.Writer.Status()),            //记录响应的状态码
		zap.String("method", ctx.Request.Method),          //记录请求方法
		zap.String("path", path),                          //记录请求的路径
		zap.String("query", query),                        //记录请求的原始查询参数
		zap.String("ip", ctx.ClientIP()),                  //记录客户端的 IP 地址
		zap.String("user-agent", ctx.Request.UserAgent()), //记录客户端的 user-agent
		zap.String("body", body),                          //记录请求的主体数据
	}
	return fields
}

// LogBody 读取 body 内容缓存下来，为之后打印日志做准备（读取请求体的内容并将其存储在 Gin 上下文中）
func LogBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bodyBytes, _ := io.ReadAll(ctx.Request.Body)
		_ = ctx.Request.Body.Close()                                //关闭原始请求主体的读取，以确保资源的正确释放
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 创建一个新地可读取的请求主体，并将之前读取的 bodyBytes 作为内容，最后将其设置回 ctx.Request.Body。
		// 使用一个新的缓冲区来存储请求主体的内容，而不是直接读取原始的请求主体。这样我们就可以在不影响原始请求主体的情况下，对请求主体的内容进行处理和修改
		ctx.Set("body", bodyBytes)
		ctx.Next()
	}
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				if brokenPipe {
					lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
