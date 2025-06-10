package setting

import (
	"chat/global"
	"chat/manager"
)

type chat struct {
}

func (chat) Init() {
	global.ChatMap = manager.NewChatMap()
}
