package routs

import (
	"chat/controller/api"
	"chat/middlewares"
	"github.com/gin-gonic/gin"
)

type message struct {
}

func (message) Init(router *gin.RouterGroup) {
	r := router.Group("/message", middlewares.MustAccount())
	{
		list := r.Group("list")
		{
			list.GET("time", api.Apis.Message.GetMsgsByRelationIDAndTime)
			//list.GET("offer", api.Apis.Message.OfferMsgsByAccountIDAndTime)
			list.GET("pin", api.Apis.Message.GetPinMsgsByRelationID)
			list.GET("reply", api.Apis.Message.GetRlyMsgsInfoByMsgID)
			list.GET("content", api.Apis.Message.GetMsgsByContent)
		}
	}

}
