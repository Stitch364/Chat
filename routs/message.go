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

		update := r.Group("update")
		{
			update.PUT("pin", api.Apis.Message.UpdateMsgPin)
			update.PUT("top", api.Apis.Message.UpdateMsgTop)
			update.PUT("revoke", api.Apis.Message.RevokeMsg)
			update.PUT("delete", api.Apis.Message.DeleteMsg)
		}
		info := r.Group("info")
		{
			info.GET("top", api.Apis.Message.GetTopMsgByRelationID)
		}
		list := r.Group("list")
		{
			list.POST("time", api.Apis.Message.GetMsgsByRelationIDAndTime)
			//list.GET("offer", api.Apis.Message.OfferMsgsByAccountIDAndTime)
			list.GET("pin", api.Apis.Message.GetPinMsgsByRelationID)
			list.GET("reply", api.Apis.Message.GetRlyMsgsInfoByMsgID)
			list.GET("content", api.Apis.Message.GetMsgsByContent)
		}
	}

}
