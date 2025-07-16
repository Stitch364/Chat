package main

import (
	"chat/global"
	"chat/routs/router"
	"chat/setting"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//Go Web开发较通用的脚手架模板

func main() {

	//1. 初始化项目（配置加载，日志、数据库，雪花算法...初始化等等）
	setting.Inits()
	//设置 Gin 框架为 Release（生产）模式
	if global.PublicSetting.Server.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	//注册成路由
	r, ws := router.NewRouter()
	_ = r.Run(":8080")
	//sever := http.Server{
	//	Addr:           global.PublicSetting.Server.HttpPort, //端口号
	//	Handler:        r,                                    //路由处理器
	//	MaxHeaderBytes: 1 << 20,                              //最大请求头大小（1MB）
	//	//设置合适的 MaxHeaderBytes 值，可以确保服务器能够有效地处理请求头，避免不必要的资源浪费或潜在的安全风险
	//}
	//global.Logger.Info("Server started!") //输出日志，服务器已启动
	//fmt.Println("AppName:", global.PublicSetting.App.Name, "Version:", global.PublicSetting.App.Version, "Address:", global.PublicSetting.Server.HttpPort, "RunMode:", global.PublicSetting.Server.RunMode)

	errChan := make(chan error, 1)
	defer close(errChan) //延迟关闭错误通道

	//go func() {
	//	//启动 HTTP 服务器
	//	err := sever.ListenAndServe()
	//	if err != nil {
	//		errChan <- err //将错误发送到错误通道
	//	}
	//}()

	defer ws.Close()
	// 接收并处理网络连接
	if err := ws.Serve(); err != nil {
		fmt.Println("Socket.IO server error:", err)
		errChan <- err
	}

	//// 启动 Socket.IO 服务器
	//go func() {
	//	defer ws.Close()
	//	// 接收并处理网络连接
	//	fmt.Println(111)
	//	if err := ws.Serve(); err != nil {
	//		errChan <- err
	//	}
	//	fmt.Println(222)
	//}()

	//// 6. 启动服务（优雅关机）
	//srv := &http.Server{
	//	Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
	//	Handler: r,
	//}
	//
	//go func() {
	//	// 开启一个goroutine启动服务
	//	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		log.Fatalf("listen: %s\n", err)
	//	}
	//}()
	//
	//// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	//quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	//// kill 默认会发送 syscall.SIGTERM 信号
	//// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	//// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	//// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	//<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	//zap.L().Info("Shutdown Server ...")
	//// 创建一个5秒超时的context
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	//if err := srv.Shutdown(ctx); err != nil {
	//	zap.L().Fatal("Server Shutdown", zap.Error(err))
	//}

	zap.L().Info("Server exiting")
}
