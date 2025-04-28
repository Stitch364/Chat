package setting

import (
	"chat/global"
	"chat/pkg/tool"
	"flag"
	"github.com/XYYSWK/Lutils/pkg/setting"
	"strings"
)

//配置文件绑定到全局结构体上（默认加载）

var (
	configPaths       string //配置文件路径
	privateConfigName string //private配置文件
	publicConfigName  string //public配置文件
	configType        string //配置文件名
)

// 通过 flag 获取命令行参数
func setupFlag() {
	//命令行参数绑定
	//使用 -name参数 绑定，没有指定默认是第三个参数；
	//第四个参数是帮助信息
	//flag.StringVar相当于是声明
	flag.StringVar(&configPaths, "config_path", global.RootDir+"/config/app", "指定要使用的配置文件的路径，多个路径用逗号 \",\" 隔开 ")
	flag.StringVar(&privateConfigName, "private_config_name", "private", "private 配置文件名")
	flag.StringVar(&publicConfigName, "public_config_name", "public", "public 配置文件名")
	flag.StringVar(&configType, "config_type", "yaml", "配置文件类型")
	flag.Parse() //解析命令行参数，并将它们对应的值赋给相应的变量
}

// 空结构体
// 用来集成，区分方法
type config struct {
}

// Init 读取配置文件，将配置信息映射到结构体中
func (config) Init() {
	//解析命令行参数
	setupFlag()
	var (
		err            error
		publicSetting  *setting.Setting //public 配置信息
		privateSetting *setting.Setting //private 配置信息
	)
	//type Setting struct {
	//	vp  *viper.Viper
	//	all interface{} //用于存储配置文件中的所有配置信息,热重载时存储配置信息和调用BindAll
	//}

	//第一个Init函数，会在其他组将的Init函数执行前，将配置文件绑定到全局变量中
	//DoThat（）方法，err不等于nil就执行后面的函数并返回err，否则返回err
	//基本逻辑就是，如果第一个参数err为nil那就执行后面的函数，并返回新的err
	//一旦err 不为nil，就返回这个err
	//类似于递归。链式调用，确保在某个环节出错时能够及时中断后续操作。
	err = tool.DoThat(err, func() error {
		//初始化 publicSetting 的基础属性
		//strings.Split(configPaths, ",") 用 , 分割
		publicSetting, err = setting.NewSetting(publicConfigName, configType, strings.Split(configPaths, ",")...) //引入配置文件路径
		return tool.DoThat(err, func() error { return publicSetting.BindAll(&global.PublicSetting) })             //将配置文件中的信息解析到全局变量中
	})
	//上面执行后的err也会影响下面
	err = tool.DoThat(err, func() error {
		//初始化 privateSetting 的基础属性
		privateSetting, err = setting.NewSetting(privateConfigName, configType, strings.Split(configPaths, ",")...) //引入配置文件路径
		return tool.DoThat(err, func() error { return privateSetting.BindAll(&global.PrivateSetting) })             //将配置文件中的信息解析到全局变量中
	})
	//处理err
	if err != nil {
		panic("读取配置文件有误：" + err.Error())
	}
}
