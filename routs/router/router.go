package router

import (
	"chat/global"
	"chat/middlewares"
	"chat/routs"
	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(middlewares.Cors, middlewares.GinLogger(), middlewares.GinRecovery(true))

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
	}

	return r
}
