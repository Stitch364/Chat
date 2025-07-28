package router

import (
	"chat/global"
	"chat/middlewares"
	"chat/routs"
	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func NewRouter() (*gin.Engine, *socketio.Server) {
	r := gin.New()
	r.Use(middlewares.Cors2, middlewares.GinLogger(), middlewares.Recovery(true))
	//r.Use(middlewares.Cors, middlewares.GinLogger2())

	root := r.Group("api", middlewares.LogBody(), middlewares.PasetoAuth())
	{
		root.GET("/ping", func(c *gin.Context) {
			reply := app.NewResponse(c)
			global.Logger.Info("ping")
			reply.Reply(nil, "pong")
		})

		rg := routs.Routers
		rg.User.Init(root)
		rg.Email.Init(root)
		rg.Account.Init(root)
		rg.Application.Init(root)
		rg.Message.Init(root)
		rg.Setting.Init(root)
		rg.Group.Init(root)
	}

	return r, routs.Routers.Chat.Init(r)
}
