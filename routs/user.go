package routs

import (
	"chat/controller/api"
	"chat/middlewares"
	"github.com/gin-gonic/gin"
)

type user struct{}

func (user) Init(router *gin.RouterGroup) {
	r := router.Group("user")
	{
		//用户注册
		r.POST("/register", api.Apis.User.Register)
		//用户登录
		r.POST("/login", api.Apis.User.Login)
		updateGroup := r.Group("update").Use(middlewares.MustUser()) //添加鉴权中间件
		{
			updateGroup.PUT("pwd", api.Apis.User.UpdateUserPassword)
			updateGroup.PUT("email", api.Apis.User.UpdateUserEmail)
			updateGroup.GET("/logout", api.Apis.User.Logout)
		}

	}
}
