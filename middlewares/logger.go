package middlewares

import (
	"bytes"
	"chat/global"
	"fmt"
	"github.com/XYYSWK/Lutils/pkg/email"
	"github.com/XYYSWK/Lutils/pkg/times"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
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
		global.Logger.Info(path,
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

func GinLogger2() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// ===== 取当前调用栈（文件名、行号、函数名）=====
		pc, file, line, ok := runtime.Caller(1) // 1 = 调用此函数的栈帧
		fnName := "unknown"
		if ok {
			fn := runtime.FuncForPC(pc)
			fnName = fn.Name()
		}

		c.Next()

		cost := time.Since(start)
		global.Logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
			// 新增字段
			zap.String("caller_file", file),
			zap.Int("caller_line", line),
			zap.String("caller_func", fnName),
		)
	}
}

//// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
//func GinRecovery(stack bool) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		defer func() {
//			if err := recover(); err != nil {
//				// Check for a broken connection, as it is not really a
//				// condition that warrants a panic stack trace.
//				var brokenPipe bool
//				if ne, ok := err.(*net.OpError); ok {
//					if se, ok := ne.Err.(*os.SyscallError); ok {
//						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
//							brokenPipe = true
//						}
//					}
//				}
//
//				httpRequest, _ := httputil.DumpRequest(c.Request, false)
//
//				if brokenPipe {
//					lg.Error(c.Request.URL.Path,
//						zap.Any("error", err),
//						zap.String("request", string(httpRequest)),
//					)
//					// If the connection is dead, we can't write a status to it.
//					c.Error(err.(error)) // nolint: errcheck
//					c.Abort()
//					return
//				}
//
//				if stack {
//					lg.Error("[Recovery from panic]",
//						zap.Any("error", err),
//						zap.String("request", string(httpRequest)),
//						zap.String("stack", string(debug.Stack())),
//					)
//				} else {
//					lg.Error("[Recovery from panic]",
//						zap.Any("error", err),
//						zap.String("request", string(httpRequest)),
//					)
//				}
//				c.AbortWithStatus(http.StatusInternalServerError)
//			}
//		}()
//		c.Next()
//	}
//}

func Recovery(stack bool) gin.HandlerFunc {
	defaultMailer := email.NewEmail(&email.SMTPInfo{
		Port:     global.PrivateSetting.Email.Port,
		IsSSL:    global.PrivateSetting.Email.IsSSL,
		Host:     global.PrivateSetting.Email.Host,
		UserName: global.PrivateSetting.Email.Username,
		Password: global.PrivateSetting.Email.Password,
		From:     global.PrivateSetting.Email.From,
	})
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查连接是否断开，因为这并不是真正需要进行恐慌堆栈跟踪的情况
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connections reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 将请求对象转换为字节切片
				httpRequest, _ := httputil.DumpRequest(ctx.Request, false)
				var body string
				data, ok := ctx.Get(Body)
				if ok {
					body = string(data.([]byte))
				}
				sendErr := defaultMailer.SendMail( // 短信通知
					global.PrivateSetting.Email.To,
					fmt.Sprintf("异常抛出，发生时间：%v\n", time.Now().Format(times.LayoutDate)),
					fmt.Sprintf("错误信息：%s\n请求信息：%s\n请求body:%s\n调用堆栈信息：%s\n", err, string(httpRequest), body, string(debug.Stack())),
				)
				if sendErr != nil {
					global.Logger.Error(fmt.Sprintf("email.SendMail Error: %v", sendErr.Error()))
				}

				if brokenPipe {
					global.Logger.Error(ctx.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)))
					// 如果连接已断开，我们就无法写入状态
					ctx.Error(err.(error)) // 将错误信息与上下文关联
					ctx.Abort()            // 阻止调用后续的处理函数
					return
				}
				if stack { // 如果需要记录堆栈信息
					global.Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack()))) // 记录当前 goroutine 的堆栈跟踪信息到日志中
				} else { // 不需要记录到堆栈信息
					global.Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("body", body))
				}
				ctx.AbortWithStatus(http.StatusInternalServerError) //阻止调用后续的处理函数，并返回“服务器内部错误”的状态码
			}
		}()
		ctx.Next()
	}
}
