package setting

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

// Conf 创建全局结构体变量，存储全部配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	//无论配置文件是啥类型，Tag都用mapstructure
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	Port         int    `mapstructure:"port"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"db"`
	PollSize int    `mapstructure:"poll_size"`
}

func Init(filename string) (err error) {
	//通过外面传入的配置文件地址读取配置
	viper.SetConfigFile(filename)
	//viper.SetConfigFile("config.yaml")
	//viper.SetConfigName("config") // 配置文件名称(无扩展名,会查找到所有叫config的文件)
	//viper.SetConfigType("yaml")   // （专用于从远程获取配置信息时指定配置）
	viper.AddConfigPath(".")   // 还可以在工作目录中查找配置
	err = viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {            // 处理读取配置文件的错误
		fmt.Println("Fatal error config file:  \n", err)
		return
	}
	//把读取到的配置信息反序列到Conf变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Println("viper.Unmarshal error config file:  \n", err)
	}
	//开启监视，并重启读取
	//当配置文件修改后，应该再次去反序列化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		if err := viper.Unmarshal(&Conf); err != nil {
			fmt.Println("viper.Unmarshal error config file:  \n", err)
		}
		return
	})
	return nil
}
