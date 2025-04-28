package setting

import (
	"chat/global"
	"github.com/XYYSWK/Lutils/pkg/logger"
)

// 空结构体
// 用来集成方法
type log struct {
}

// Init 日志
func (log) Init() {
	// Newlogger就直接初始化好了log对象，之后直接调用log写日志就行
	global.Logger = logger.NewLogger(&logger.InitStruct{
		LogSavePath:   global.PublicSetting.Log.LogSavePath,
		LogFileExt:    global.PublicSetting.Log.LogFileExt,
		MaxSize:       global.PublicSetting.Log.MaxSize,
		MaxBackups:    global.PublicSetting.Log.MaxBackups,
		MaxAge:        global.PublicSetting.Log.MaxAge,
		Compress:      global.PublicSetting.Log.Compress,
		LowLevelFile:  global.PublicSetting.Log.LowLevelFile,
		HighLevelFile: global.PublicSetting.Log.HighLevelFile,
	}, global.PublicSetting.Log.Level)
}
