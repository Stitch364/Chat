package task

import (
	"chat/dao"
	"chat/global"
	"chat/model/chat"
	"chat/model/reply"
)

//有关消息的推送任务

func PublishMsg(msg reply.ParamMsgInfoWithRly) func() {
	return func() {
		ctx, cancel := global.DefaultContextWithTimeout()
		defer cancel()
		accountIDs, err := dao.Database.Redis.GetAllAccountsByRelationID(ctx, msg.RelationID)
		if err != nil {
			global.Logger.Error(err.Error())
			return
		}
		for _, accountID := range accountIDs {
			if global.ChatMap.CheckIsOnConnection(accountID) {
				global.ChatMap.Send(accountID, chat.ClientSendMsg, msg)
			} else {
				//用户离线，将消息发送至MQ中

			}

		}
	}
}
