package main

import (
	"chat/global"
	"chat/routs/router"
	"chat/setting"
	"github.com/gin-gonic/gin"
	//"web_app/dao/mysql"
	//"web_app/dao/redis"
	//"web_app/logger"
	//"web_app/routs"
	//"web_app/setting"

	"go.uber.org/zap"
)

//Go Web开发较通用的脚手架模板

func main() {
	//if len(os.Args) < 2 {
	//	fmt.Println("need config file.eg")
	//	return
	//}
	//// 1. 加载配置
	//if err := setting.Init(os.Args[1]); err != nil {
	//	fmt.Println("init settings failed, err:", err)
	//	return
	//}
	//// 2. 初始化日志
	//if err := logger.Init(setting.Conf.LogConfig); err != nil {
	//	fmt.Println("init settings failed, err:", err)
	//	return
	//}
	////把缓冲区的文件增加到日志里
	//defer func(l *zap.Logger) {
	//	err := l.Sync()
	//	if err != nil {
	//		zap.L().Error("sync logger failed", zap.Error(err))
	//	}
	//}(zap.L())
	//
	////zap.L().
	//// 3. 初始化MySQL连接
	//if err := mysql.Init(setting.Conf.MySQLConfig); err != nil {
	//	fmt.Println("init settings failed, err:", err)
	//	return
	//}
	//defer mysql.Close()
	//// 4. 初始化Redis连接
	//if err := redis.Init(setting.Conf.RedisConfig); err != nil {
	//	fmt.Println("init settings failed, err:", err)
	//	return
	//}
	//defer redis.Close()
	//// 5. 注册路由
	//r := routs.Setup()

	//1. 初始化项目（配置加载，日志、数据库，雪花算法...初始化等等）
	setting.Inits()
	//设置 Gin 框架为 Release（生产）模式
	if global.PublicSetting.Server.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	//注册成路由
	r := router.NewRouter()
	_ = r.Run(":8080")

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
