package setting

import (
	"chat/global"
	"github.com/XYYSWK/Lutils/pkg/app"
)

type page struct{}

// Init 分页器初始化
func (page) Init() {
	//调用page工具
	global.Page = app.InitPage(global.PublicSetting.Page.DefaultPageSize, global.PublicSetting.Page.MaxPageSize, global.PublicSetting.Page.PageKey, global.PublicSetting.Page.PageSizeKey)
}
