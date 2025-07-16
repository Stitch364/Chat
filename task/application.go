package task

import (
	"chat/global"
	"chat/model/chat"
)

func Application(accountID int64) func() {
	return func() {
		global.ChatMap.Send(accountID, chat.ServerApplication)
	}
}
