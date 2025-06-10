package setting

type group struct {
	Config     config
	Database   database
	Log        log
	Page       page
	EmailMark  mark
	Worker     worker
	TokenMaker tokenMaker
	GenerateID generateID
	Chat       chat
	Load       load
}

var Group = new(group)

// Inits 初始化项目
func Inits() {
	Group.Config.Init()
	Group.Database.Init()
	Group.Log.Init()
	Group.Page.Init()
	Group.EmailMark.Init()
	Group.Worker.Init()
	Group.TokenMaker.Init()
	Group.GenerateID.Init()
	Group.Chat.Init()
	Group.Load.Init()
}
